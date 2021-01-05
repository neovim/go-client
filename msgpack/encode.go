package msgpack

import (
	"reflect"
	"sync"
)

// Marshaler is the interface implemented by objects that can encode themselves
// to a MessagePack stream.
type Marshaler interface {
	MarshalMsgPack(e *Encoder) error
}

type encodeTypeError struct {
	Type reflect.Type
}

func (e *encodeTypeError) Error() string {
	return "msgpack: unsupported type: " + e.Type.String()
}

func encodeUnsupportedType(e *Encoder, v reflect.Value) {
	abort(&encodeTypeError{v.Type()})
}

// Encode writes the MessagePack encoding of v to the stream.
//
// Encode traverses the value v recursively. If an encountered value implements
// the Marshaler interface Encode calls its MarshalMsgPack method to write the
// value to the stream.
//
// Otherwise, Encode uses the following type-dependent default encodings:
//
//  Go Type             MessagePack Type
//  bool                true or false
//  float32, float64    float64
//  string              string
//  []byte              binary
//  slices, arrays      array
//  struct, map         map
//
// Struct values encode as maps or arrays. If any struct field tag specifies
// the "array" option, then the struct is encoded as an array. Otherwise, the
// struct is encoded as a map.  Each exported struct field becomes a member of
// the map unless
//   - the field's tag is "-", or
//   - the field is empty and its tag specifies the "omitempty" option.
//
// Anonymous struct fields are marshaled as if their inner exported fields
// were fields in the outer struct.
//
// The struct field tag "empty" specifies a default value when decoding and the
// empty value for the "omitempty" option.
//
// Pointer values encode as the value pointed to. A nil pointer encodes as the
// MessagePack nil value.
//
// Interface values encode as the value contained in the interface. A nil
// interface value encodes as the MessagePack nil value.
func (e *Encoder) Encode(v interface{}) (err error) {
	if v == nil {
		return e.PackNil()
	}
	defer handleAbort(&err)

	rv := reflect.ValueOf(v)
	encoderForType(rv.Type(), nil)(e, rv)

	return nil
}

type encodeFunc func(e *Encoder, v reflect.Value)

type encodeBuilder struct {
	m map[reflect.Type]encodeFunc
}

var encodeFuncCache struct {
	sync.RWMutex
	m map[reflect.Type]encodeFunc
}

func encoderForType(t reflect.Type, b *encodeBuilder) encodeFunc {
	encodeFuncCache.RLock()
	f, ok := encodeFuncCache.m[t]
	encodeFuncCache.RUnlock()
	if ok {
		return f
	}

	save := false
	if b == nil {
		b = &encodeBuilder{m: make(map[reflect.Type]encodeFunc)}
		save = true
	} else if f, ok := b.m[t]; ok {
		return f
	}

	// Add temporary entry to break recursion.
	b.m[t] = func(e *Encoder, v reflect.Value) {
		f(e, v)
	}
	f = b.encoder(t)
	b.m[t] = f

	if save {
		encodeFuncCache.Lock()

		if encodeFuncCache.m == nil {
			encodeFuncCache.m = make(map[reflect.Type]encodeFunc)
		}
		for t, f := range b.m {
			encodeFuncCache.m[t] = f
		}

		encodeFuncCache.Unlock()
	}

	return f
}

func (b *encodeBuilder) encoder(t reflect.Type) encodeFunc {
	if t.Implements(marshalerType) {
		return b.marshalEncoder(t)
	}

	var f encodeFunc
	switch t.Kind() {
	case reflect.Bool:
		f = boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f = intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		f = uintEncoder
	case reflect.Float32, reflect.Float64:
		f = floatEncoder
	case reflect.String:
		f = stringEncoder
	case reflect.Array:
		f = b.arrayEncoder(t)
	case reflect.Slice:
		f = b.sliceEncoder(t)
	case reflect.Map:
		f = b.mapEncoder(t)
	case reflect.Interface:
		f = interfaceEncoder
	case reflect.Struct:
		f = b.structEncoder(t)
	case reflect.Ptr:
		f = b.ptrEncoder(t)
	default:
		f = encodeUnsupportedType
	}

	if t.Kind() != reflect.Ptr && reflect.PtrTo(t).Implements(marshalerType) {
		f = marshalAddrEncoder{f}.encode
	}

	return f
}

func nilEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackNil(); err != nil {
		abort(err)
	}
}

func boolEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackBool(v.Bool()); err != nil {
		abort(err)
	}
}

func intEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackInt(v.Int()); err != nil {
		abort(err)
	}
}

func uintEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackUint(v.Uint()); err != nil {
		abort(err)
	}
}

func floatEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackFloat(v.Float()); err != nil {
		abort(err)
	}
}

func stringEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackString(v.String()); err != nil {
		abort(err)
	}
}

func byteSliceEncoder(e *Encoder, v reflect.Value) {
	if err := e.PackBinary(v.Bytes()); err != nil {
		abort(err)
	}
}

func interfaceEncoder(e *Encoder, v reflect.Value) {
	if !v.IsValid() || v.IsNil() {
		nilEncoder(e, v)
		return
	}

	v = v.Elem()
	encoderForType(v.Type(), nil)(e, v)
}

type ptrEncoder struct{ elem encodeFunc }

func (enc ptrEncoder) encode(e *Encoder, v reflect.Value) {
	if v.IsNil() {
		nilEncoder(e, v)
		return
	}

	enc.elem(e, v.Elem())
}

func (b *encodeBuilder) ptrEncoder(t reflect.Type) encodeFunc {
	return ptrEncoder{encoderForType(t.Elem(), b)}.encode
}

type mapEncoder struct{ key, elem encodeFunc }

func (enc *mapEncoder) encode(e *Encoder, v reflect.Value) {
	if v.IsNil() {
		nilEncoder(e, v)
		return
	}

	if err := e.PackMapLen(int64(v.Len())); err != nil {
		abort(err)
	}

	for _, k := range v.MapKeys() {
		enc.key(e, k)
		enc.elem(e, v.MapIndex(k))
	}
}

func (b *encodeBuilder) mapEncoder(t reflect.Type) encodeFunc {
	enc := &mapEncoder{key: encoderForType(t.Key(), b), elem: encoderForType(t.Elem(), b)}
	return enc.encode
}

type sliceArrayEncoder struct{ elem encodeFunc }

func (enc sliceArrayEncoder) encodeArray(e *Encoder, v reflect.Value) {
	if err := e.PackArrayLen(int64(v.Len())); err != nil {
		abort(err)
	}

	for i := 0; i < v.Len(); i++ {
		enc.elem(e, v.Index(i))
	}
}

func (b *encodeBuilder) arrayEncoder(t reflect.Type) encodeFunc {
	return sliceArrayEncoder{encoderForType(t.Elem(), b)}.encodeArray
}

func (enc sliceArrayEncoder) encodeSlice(e *Encoder, v reflect.Value) {
	if v.IsNil() {
		nilEncoder(e, v)
		return
	}

	enc.encodeArray(e, v)
}

func (b *encodeBuilder) sliceEncoder(t reflect.Type) encodeFunc {
	if t.Elem().Kind() == reflect.Uint8 {
		return byteSliceEncoder
	}

	return sliceArrayEncoder{encoderForType(t.Elem(), b)}.encodeSlice
}

var marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()

func marshalPtrEncoder(e *Encoder, v reflect.Value) {
	if v.IsNil() {
		nilEncoder(e, v)
		return
	}

	marshalEncoder(e, v)
}

func marshalEncoder(e *Encoder, v reflect.Value) {
	m := v.Interface().(Marshaler)

	if err := m.MarshalMsgPack(e); err != nil {
		abort(err)
	}
}

func (b *encodeBuilder) marshalEncoder(t reflect.Type) encodeFunc {
	if t.Kind() == reflect.Ptr {
		return marshalPtrEncoder
	}

	return marshalEncoder
}

type marshalAddrEncoder struct{ f encodeFunc }

func (enc marshalAddrEncoder) encode(e *Encoder, v reflect.Value) {
	if v.CanAddr() {
		marshalEncoder(e, v.Addr())
		return
	}

	enc.f(e, v)
}

type fieldEnc struct {
	name  string
	empty func(reflect.Value) bool
	f     encodeFunc
	index []int
}

type structEncoder []*fieldEnc

func (enc structEncoder) encode(e *Encoder, v reflect.Value) {
	var n int64
	for _, fe := range enc {
		fv := fieldByIndex(v, fe.index)
		if !fv.IsValid() || (fe.empty != nil && fe.empty(fv)) {
			continue
		}
		n++
	}

	if err := e.PackMapLen(n); err != nil {
		abort(err)
	}

	for _, fe := range enc {
		fv := fieldByIndex(v, fe.index)
		if !fv.IsValid() || (fe.empty != nil && fe.empty(fv)) {
			continue
		}

		if err := e.PackString(fe.name); err != nil {
			abort(err)
		}

		fe.f(e, fv)
	}
}

func (enc structEncoder) encodeArray(e *Encoder, v reflect.Value) {
	if err := e.PackArrayLen(int64(len(enc))); err != nil {
		abort(err)
	}

	for _, fe := range enc {
		fv := fieldByIndex(v, fe.index)
		fe.f(e, fv)
	}
}

func (b *encodeBuilder) structEncoder(t reflect.Type) encodeFunc {
	fields, array := fieldsForType(t)
	enc := make(structEncoder, len(fields))

	for i, f := range fields {
		var empty func(reflect.Value) bool
		if f.omitEmpty {
			empty = emptyFunc(f)
		}
		enc[i] = &fieldEnc{
			name:  f.name,
			empty: empty,
			index: f.index,
			f:     encoderForType(f.typ, b)}
	}

	if array {
		return enc.encodeArray
	}

	return enc.encode
}

func emptyFunc(f *field) func(reflect.Value) bool {
	if f.empty.IsValid() {
		return func(v reflect.Value) bool { return v.Interface() == f.empty.Interface() }
	}

	switch f.typ.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return lenEmpty
	case reflect.Bool:
		return boolEmpty
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEmpty
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEmpty
	case reflect.Float32, reflect.Float64:
		return floatEmpty
	case reflect.Interface, reflect.Ptr:
		return nilEmpty
	default:
		return nil
	}
}

func lenEmpty(v reflect.Value) bool   { return v.Len() == 0 }
func boolEmpty(v reflect.Value) bool  { return !v.Bool() }
func intEmpty(v reflect.Value) bool   { return v.Int() == 0 }
func uintEmpty(v reflect.Value) bool  { return v.Uint() == 0 }
func floatEmpty(v reflect.Value) bool { return v.Float() == 0 }
func nilEmpty(v reflect.Value) bool   { return v.IsNil() }
