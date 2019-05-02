// Package rpc implements MessagePack RPC.
package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/neovim/go-client/msgpack"
)

const (
	requestMessage      = 0
	replyMessage        = 1
	notificationMessage = 2
)

const (
	stateInit = iota
	stateClosed
)

var (
	errClosed   = errors.New("msgpack/rpc: session closed")
	errInternal = errors.New("msgpack/rpc: internal error")
)

type Error struct {
	Value interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Value)
}

type Call struct {
	Method string
	Args   interface{}
	Reply  interface{}
	Err    error
	Done   chan *Call
}

func (call *Call) done(e *Endpoint, err error) {
	call.Err = err
	select {
	case call.Done <- call:
		// ok
	default:
		e.logf("msgpack/rpc: done channel over capacity for method %s", call.Method)
	}
}

type handler struct {
	fn   reflect.Value
	args []reflect.Value
}

type notification struct {
	call   func([]reflect.Value) []reflect.Value
	args   []reflect.Value
	method string
	next   *notification
}

// Endpoint represents a MessagePack RPC peer.
type Endpoint struct {
	logf   func(fmt string, args ...interface{})
	arg    reflect.Value
	dec    *msgpack.Decoder
	closer io.Closer
	done   chan struct{}

	packMu sync.Mutex
	enc    *msgpack.Encoder
	bw     *bufio.Writer

	mu      sync.Mutex
	id      uint64
	pending map[uint64]*Call
	state   int
	err     error

	handlersMu sync.RWMutex
	handlers   map[string]*handler

	notificationsMu   sync.Mutex
	notificationsCond *sync.Cond
	notifications     []*notification
}

// NewEndpoint returns a new endpoint with the specified options.
func NewEndpoint(r io.Reader, w io.Writer, c io.Closer, options ...Option) (*Endpoint, error) {
	bw := bufio.NewWriter(w)
	e := &Endpoint{
		done:     make(chan struct{}),
		handlers: make(map[string]*handler),
		pending:  make(map[uint64]*Call),
		closer:   c,
		bw:       bw,
		enc:      msgpack.NewEncoder(bw),
		dec:      msgpack.NewDecoder(r),
	}
	for _, option := range options {
		option.f(e)
	}
	return e, nil

}

type Option struct{ f func(*Endpoint) }

func WithExtensions(extensions msgpack.ExtensionMap) Option {
	return Option{func(e *Endpoint) {
		e.dec.SetExtensions(extensions)
	}}
}

func WithLogf(f func(fmt string, args ...interface{})) Option {
	return Option{func(e *Endpoint) {
		e.logf = f
	}}
}

func (e *Endpoint) decodeUint(what string) (uint64, error) {
	if err := e.dec.Unpack(); err != nil {
		return 0, err
	}
	t := e.dec.Type()
	if t != msgpack.Uint && t != msgpack.Int {
		return 0, fmt.Errorf("msgpack/rpc: error decoding %s, found %s", what, e.dec.Type())
	}
	return e.dec.Uint(), nil
}

func (e *Endpoint) decodeString(what string) (string, error) {
	if err := e.dec.Unpack(); err != nil {
		return "", err
	}
	if e.dec.Type() != msgpack.String {
		return "", fmt.Errorf("msgpack/rpc: error decoding %s, found %s", what, e.dec.Type())
	}
	return e.dec.String(), nil
}

func (e *Endpoint) skip(n int) error {
	for i := 0; i < n; i++ {
		if err := e.dec.Unpack(); err != nil {
			return err
		}
		if err := e.dec.Skip(); err != nil {
			return err
		}
	}
	return nil
}

var errorType = reflect.ValueOf(new(error)).Elem().Type()

// Register registers handler fn for the specified method name.
//
// When servicing a call, the arguments to fn are the values in args followed
// by the values passed from the peer.
func (e *Endpoint) Register(method string, fn interface{}, args ...interface{}) error {
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		return errors.New("msgpack/rpc: handler not a function")
	}
	if t.NumIn() < len(args) {
		return fmt.Errorf("msgpack/rpc: handler must have at least %d args", len(args))
	}

	h := &handler{fn: v, args: make([]reflect.Value, len(args))}

	for i, arg := range args {
		if arg == nil {
			t := t.In(i)
			switch t.Kind() {
			case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
				h.args[i] = reflect.New(t).Elem()
			default:
				return fmt.Errorf("msgpack/rpc: handler arg %d must be interface, pointer, map or slice", i)
			}
		} else {
			h.args[i] = reflect.ValueOf(arg)
			if t.In(i) != h.args[i].Type() {
				return fmt.Errorf("msgpack/rpc: handler arg %d must be type %T", i, arg)
			}
		}
	}

	if t.NumOut() > 2 || (t.NumOut() > 0 && t.Out(t.NumOut()-1) != errorType) {
		return errors.New("msgpack/rpc: handler return must be (), (error) or (valueType, error)")
	}

	e.handlersMu.Lock()
	e.handlers[method] = h
	e.handlersMu.Unlock()
	return nil
}

func (e *Endpoint) close(err error) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.state == stateClosed {
		return e.err
	}
	e.state = stateClosed
	e.err = err
	for _, call := range e.pending {
		call.done(e, errClosed)
	}
	e.pending = nil
	err = e.closer.Close()
	if e.err == nil {
		e.err = err
	}
	return e.err
}

// Close releases the resources used by endpoint.
func (e *Endpoint) Close() error {
	return e.close(nil)
}

func (e *Endpoint) Call(method string, reply interface{}, args ...interface{}) error {
	c := <-e.Go(method, make(chan *Call, 1), reply, args...).Done
	return c.Err
}

func (e *Endpoint) Go(method string, done chan *Call, reply interface{}, args ...interface{}) *Call {
	if args == nil {
		args = []interface{}{}
	}

	if done == nil {
		done = make(chan *Call, 1)
	} else if cap(done) == 0 {
		panic("unbuffered done channel")
	}

	call := &Call{
		Method: method,
		Args:   args,
		Reply:  reply,
		Done:   done,
	}

	e.mu.Lock()
	if e.state == stateClosed {
		call.done(e, errClosed)
		e.mu.Unlock()
		return call
	}
	e.id = (e.id + 1) & 0x7fffffff
	id := e.id
	e.pending[id] = call
	e.mu.Unlock()

	message := &struct {
		Kind   int `msgpack:",array"`
		ID     uint64
		Method string
		Args   []interface{}
	}{
		requestMessage,
		id,
		method,
		args,
	}

	e.packMu.Lock()
	err := e.enc.Encode(message)
	if e := e.bw.Flush(); err == nil {
		err = e
	}
	e.packMu.Unlock()

	if err != nil {
		e.mu.Lock()
		if _, pending := e.pending[id]; pending {
			delete(e.pending, id)
			call.done(e, err)
		}
		e.mu.Unlock()
		e.close(fmt.Errorf("msgpack/rpc: error encoding %s: %v", call.Method, err))
	}

	return call
}

func (e *Endpoint) Notify(method string, args ...interface{}) error {
	if args == nil {
		args = []interface{}{}
	}

	message := &struct {
		Kind   int `msgpack:",array"`
		Method string
		Args   []interface{}
	}{
		notificationMessage,
		method,
		args,
	}

	e.packMu.Lock()
	err := e.enc.Encode(message)
	if e := e.bw.Flush(); err == nil {
		err = e
	}
	e.packMu.Unlock()
	if err != nil {
		e.close(fmt.Errorf("msgpack/rpc: error encoding %s: %v", method, err))
	}
	return err
}

func (e *Endpoint) reply(id uint64, replyErr error, reply interface{}) error {
	e.packMu.Lock()
	defer e.packMu.Unlock()

	err := e.enc.PackArrayLen(4)
	if err != nil {
		return err
	}

	err = e.enc.PackUint(replyMessage)
	if err != nil {
		return err
	}

	err = e.enc.PackUint(id)
	if err != nil {
		return err
	}

	if replyErr == nil {
		err = e.enc.PackNil()
	} else if ee, ok := replyErr.(Error); ok {
		err = e.enc.Encode(ee.Value)
	} else if ee, ok := replyErr.(msgpack.Marshaler); ok {
		err = ee.MarshalMsgPack(e.enc)
	} else {
		err = e.enc.PackString(replyErr.Error())
	}
	if err != nil {
		return err
	}

	err = e.enc.Encode(reply)
	if err != nil {
		return err
	}
	return e.bw.Flush()
}

// Serve serves incoming requests. Serve blocks until the peer disconnects or
// there is an error.
func (e *Endpoint) Serve() error {
	e.notificationsCond = sync.NewCond(&e.notificationsMu)
	defer e.enqueNotification(nil)
	go e.runNotifications()

	for {
		if err := e.dec.Unpack(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return e.close(err)
		}

		messageLen := e.dec.Len()
		if messageLen < 1 {
			return e.close(fmt.Errorf("msgpack/rpc: invalid message length %d", messageLen))
		}

		messageType, err := e.decodeUint("message type")
		if err != nil {
			return e.close(err)
		}

		switch messageType {
		case requestMessage:
			err = e.handleRequest(messageLen)
		case replyMessage:
			err = e.handleReply(messageLen)
		case notificationMessage:
			err = e.handleNotification(messageLen)
		default:
			err = fmt.Errorf("msgpack/rpc: unknown message type %d", messageType)
		}
		if err != nil {
			return e.close(err)
		}
	}
}

func (e *Endpoint) handleReply(messageLen int) error {
	if messageLen != 4 {
		// messageType, id, error, reply
		return fmt.Errorf("msgpack/rpc: invalid reply message length %d", messageLen)
	}

	id, err := e.decodeUint("response id")
	if err != nil {
		return err
	}

	e.mu.Lock()
	call := e.pending[id]
	delete(e.pending, id)
	e.mu.Unlock()

	if call == nil {
		e.logf("msgpack/rpc: no pending call for reply %d", id)
		return e.skip(2)
	}

	var errorValue interface{}
	if err := e.dec.Decode(&errorValue); err != nil {
		call.done(e, errInternal)
		return fmt.Errorf("msgpack/rpc: error decoding error value: %v", err)
	}

	if errorValue != nil {
		err := e.skip(1)
		call.done(e, Error{errorValue})
		return err
	}

	if call.Reply == nil {
		err = e.skip(1)
	} else {
		err = e.dec.Decode(call.Reply)
		if cvterr, ok := err.(*msgpack.DecodeConvertError); ok {
			call.done(e, cvterr)
			return nil
		}
	}

	if err != nil {
		call.done(e, errInternal)
		return fmt.Errorf("msgpack/rpc: error decoding reply: %v", err)
	}

	call.done(e, nil)
	return nil
}

func (e *Endpoint) createCall(h *handler) (func([]reflect.Value) []reflect.Value, []reflect.Value, error) {
	t := h.fn.Type()
	args := make([]reflect.Value, t.NumIn())
	for i := range h.args {
		args[i] = h.args[i]
	}
	if err := e.dec.Unpack(); err != nil {
		return nil, nil, err
	}
	if e.dec.Type() != msgpack.ArrayLen {
		e.dec.Skip()
		return nil, nil, fmt.Errorf("msgpack/rpc: expected args array, found %s", e.dec.Type())
	}

	// Decode plain arguments.

	var savedErr error

	srcIndex := 0
	srcLen := e.dec.Len()

	dstIndex := len(h.args)
	dstLen := t.NumIn()
	if t.IsVariadic() {
		dstLen--
	}

	for dstIndex < dstLen {
		v := reflect.New(t.In(dstIndex))
		args[dstIndex] = v.Elem()
		dstIndex++
		if srcIndex < srcLen {
			srcIndex++
			err := e.dec.Decode(v.Interface())
			if _, ok := err.(*msgpack.DecodeConvertError); ok {
				if savedErr == nil {
					savedErr = err
				}
			} else if err != nil {
				return nil, nil, err
			}
		}
	}

	if !t.IsVariadic() {
		// Skip extra arguments

		n := srcLen - srcIndex
		if n > 0 {
			err := e.skip(n)
			if err != nil {
				return nil, nil, err
			}
		}

		return h.fn.Call, args, savedErr
	}

	if srcIndex >= srcLen {
		args[dstIndex] = reflect.Zero(t.In(dstIndex))
		return h.fn.CallSlice, args, savedErr
	}

	n := srcLen - srcIndex
	v := reflect.MakeSlice(t.In(dstIndex), n, n)
	args[dstIndex] = v

	for i := 0; i < n; i++ {
		err := e.dec.Decode(v.Index(i).Addr().Interface())
		if _, ok := err.(*msgpack.DecodeConvertError); ok {
			if savedErr == nil {
				savedErr = err
			}
		} else if err != nil {
			return nil, nil, err
		}
	}

	return h.fn.CallSlice, args, nil
}

func (e *Endpoint) handleRequest(messageLen int) error {
	if messageLen != 4 {
		// messageType, id, method, args
		return fmt.Errorf("msgpack/rpc: invalid request message length %d", messageLen)
	}

	id, err := e.decodeUint("request id")
	if err != nil {
		return err
	}

	method, err := e.decodeString("service method name")
	if err != nil {
		return err
	}

	e.handlersMu.RLock()
	h, ok := e.handlers[method]
	e.handlersMu.RUnlock()

	if !ok {
		if err := e.skip(1); err != nil {
			return err
		}
		e.logf("msgpack/rpc: request service method %s not found", method)
		return e.reply(id, fmt.Errorf("unknown request method: %s", method), nil)
	}

	call, args, err := e.createCall(h)
	if _, ok := err.(*msgpack.DecodeConvertError); ok {
		e.logf("msgpack/rpc: %s: %v", method, err)
		return e.reply(id, errors.New("invalid argument"), nil)
	} else if err != nil {
		return err
	}

	go func() {
		out := call(args)
		var replyErr error
		var replyVal interface{}
		switch h.fn.Type().NumOut() {
		case 1:
			replyErr, _ = out[0].Interface().(error)
		case 2:
			replyVal = out[0].Interface()
			replyErr, _ = out[1].Interface().(error)
		}
		if err := e.reply(id, replyErr, replyVal); err != nil {
			e.close(err)
		}
	}()

	return nil
}

func (e *Endpoint) handleNotification(messageLen int) error {
	// messageType, method, args
	if messageLen != 3 {
		return fmt.Errorf("msgpack/rpc: invalid notification message length %d", messageLen)
	}

	method, err := e.decodeString("service method name")
	if err != nil {
		return err
	}

	e.handlersMu.RLock()
	h, ok := e.handlers[method]
	e.handlersMu.RUnlock()

	if !ok {
		e.logf("msgpack/rpc: notification service method %s not found", method)
		return e.skip(1)
	}

	call, args, err := e.createCall(h)
	if err != nil {
		return err
	}

	e.enqueNotification(&notification{call: call, args: args, method: method})
	return nil
}

func (e *Endpoint) enqueNotification(n *notification) {
	e.notificationsMu.Lock()
	e.notifications = append(e.notifications, n)
	e.notificationsCond.Signal()
	e.notificationsMu.Unlock()
}

func (e *Endpoint) dequeueNotifications() []*notification {
	e.notificationsMu.Lock()
	for e.notifications == nil {
		e.notificationsCond.Wait()
	}
	notifications := e.notifications
	e.notifications = nil
	e.notificationsMu.Unlock()
	return notifications
}

// runNotifications runs notifications in a single goroutine to ensure that the
// notifications are processed in order by the application.
func (e *Endpoint) runNotifications() {
	for {
		notifications := e.dequeueNotifications()
		for _, n := range notifications {
			if n == nil {
				// Serve() enqueues nil on return
				return
			}
			out := n.call(n.args)
			if len(out) > 0 {
				replyErr, _ := out[len(out)-1].Interface().(error)
				if replyErr != nil {
					e.logf("msgpack/rpc: service method %s returned %v", n.method, replyErr)
				}
			}
		}
	}
}
