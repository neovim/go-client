package msgpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
)

// pack packs the values vs and returns the result.
//
//  Go Type     Encoder method
//  nil         PackNil
//  bool        PackBool
//  int64       PackInt
//  uint64      PackUint
//  float64     PackFloat
//  arrayLen    PackArrayLen
//  mapLen      PackMapLen
//  string      PackString(s, false)
//  []byte      PackBytes(s, true)
//  extension   PackExtension(k, d)
func pack(vs ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	for _, v := range vs {
		var err error
		switch v := v.(type) {
		case int64:
			err = enc.PackInt(v)
		case uint64:
			err = enc.PackUint(v)
		case bool:
			err = enc.PackBool(v)
		case float64:
			err = enc.PackFloat(v)
		case arrayLen:
			err = enc.PackArrayLen(int64(v))
		case mapLen:
			err = enc.PackMapLen(int64(v))
		case string:
			err = enc.PackString(v)
		case []byte:
			err = enc.PackBinary(v)
		case extension:
			err = enc.PackExtension(v.k, []byte(v.d))
		case nil:
			err = enc.PackNil()
		default:
			err = fmt.Errorf("no pack for type %T", v)
		}

		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// unpack unpacks a byte slice to the following types.
//
//   Type      Go
//   Nil       nil
//   Bool      bool
//   Int       int
//   Uint      int
//   Float     float64
//   ArrayLen  arrayLen
//   MapLen    mapLen
//   String    string
//   Binary    []byte
//   Extension extension
//
// This function is not suitable for unpack tests because the integer and float
// types are mapped to int and float64 respectively.
func unpack(p []byte) ([]interface{}, error) {
	var data []interface{}
	u := NewDecoder(bytes.NewReader(p))

	for {
		err := u.Unpack()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var v interface{}
		switch u.Type() {
		case Nil:
			v = nil
		case Bool:
			v = u.Bool()
		case Int:
			v = int(u.Int())
		case Uint:
			v = int(u.Uint())
		case Float:
			v = u.Float()
		case Binary:
			v = u.Bytes()
		case String:
			v = u.String()
		case ArrayLen:
			v = arrayLen(u.Int())
		case MapLen:
			v = mapLen(u.Int())
		case Extension:
			v = extension{u.Extension(), u.String()}
		default:
			return nil, fmt.Errorf("unpack %d not handled", u.Type())
		}

		data = append(data, v)
	}

	return data, nil
}

type (
	arrayLen uint32
)

type (
	mapLen uint32
)

type extension struct {
	k int
	d string
}

type testExtension1 struct {
	data []byte
}

// Make sure a testExtension1 implements the Marshaler and Unmarshaler interfaces.
var _ Marshaler = (*testExtension1)(nil)
var _ Unmarshaler = (*testExtension1)(nil)

func (x testExtension1) MarshalMsgPack(enc *Encoder) error {
	return enc.PackExtension(1, x.data)
}

func (x *testExtension1) UnmarshalMsgPack(dec *Decoder) error {
	if dec.Type() != Extension || dec.Extension() != 1 {
		err := &DecodeConvertError{
			SrcType:  dec.Type(),
			DestType: reflect.TypeOf(x),
		}
		dec.Skip()
		return err
	}
	x.data = dec.Bytes()
	return nil
}

var testExtensionMap = ExtensionMap{
	1: func(data []byte) (interface{}, error) { return testExtension1{data}, nil },
}

func ptrInt(i int) *int {
	return &i
}

func ptrUint(u uint) *uint {
	return &u
}

func makeString(sz int) string {
	var s strings.Builder
	var x int
	for i := 0; i < sz; i++ {
		if x >= 10 {
			x = 0
		}
		s.WriteByte(byte(x + 48))
		x++
	}
	return s.String()
}

type testReader struct {
	p   []byte
	pos int
}

func NewTestReader(b []byte) io.Reader {
	return &testReader{p: b}
}

func (r *testReader) Read(b []byte) (int, error) {
	n := copy(b, r.p[r.pos:])
	if n < len(r.p) {
		r.pos = r.pos + n
	}

	if r.pos >= len(r.p) {
		r.pos = 0
	}
	return n, nil
}

type WriteByteWriter interface {
	io.Writer
	io.ByteWriter
}

type testArrayBuilder struct {
	buffer []interface{}
	tb     testing.TB
}

func NewTestArrayBuilder(tb testing.TB) *testArrayBuilder {
	return &testArrayBuilder{
		tb: tb,
	}
}

func (e *testArrayBuilder) Add(v interface{}) {
	e.buffer = append(e.buffer, v)
}

func (e *testArrayBuilder) Count() int64 {
	return int64(len(e.buffer))
}

func (e testArrayBuilder) encode(w WriteByteWriter) {
	e.tb.Helper()

	c := len(e.buffer)
	switch {
	case c < 16:
		if err := w.WriteByte(fixArrayCodeMin + byte(c)); err != nil {
			e.tb.Fatalf("msgpack: failed to write fixed array header: %v", err)
		}
	case c < math.MaxUint16:
		if err := w.WriteByte(array16Code); err != nil {
			e.tb.Fatalf("msgpack: failed to write 16-bit array header prefix: %v", err)
		}
		b := make([]byte, 5)
		binary.BigEndian.PutUint16(b, uint16(c))
		if _, err := w.Write(b); err != nil {
			e.tb.Fatalf("msgpack: failed to write 16-bit array header: %v", err)
		}
	case c < math.MaxUint32:
		if err := w.WriteByte(array32Code); err != nil {
			e.tb.Fatalf("msgpack: failed to write 32-bit array header prefix: %v", err)
		}
		b := make([]byte, 7)
		binary.BigEndian.PutUint32(b, uint32(c))
		if _, err := w.Write(b); err != nil {
			e.tb.Fatalf("msgpack: failed to write 32-bit array header: %v", err)
		}
	default:
		e.tb.Fatalf("msgpack: array element count out of range (%d)", c)
	}

	enc := NewEncoder(w)
	for _, v := range e.buffer {
		if err := enc.Encode(v); err != nil {
			e.tb.Fatalf("msgpack: failed to encode array element %s: %v", reflect.TypeOf(v), err)
		}
	}
}

func (e testArrayBuilder) Bytes() []byte {
	var buf bytes.Buffer
	e.encode(&buf)
	return buf.Bytes()
}

type testMapBuilder struct {
	valuas []interface{}
	tb     testing.TB
}

func NewTestMapBuilder(tb testing.TB) *testMapBuilder {
	return &testMapBuilder{
		tb: tb,
	}
}

func (m *testMapBuilder) Add(key string, value interface{}) {
	m.valuas = append(m.valuas, key, value)
}

func (e *testMapBuilder) Count() int64 {
	return int64(len(e.valuas)) / 2
}

func (m *testMapBuilder) encode(w WriteByteWriter) {
	m.tb.Helper()

	c := len(m.valuas) / 2
	switch {
	case c < 16:
		if err := w.WriteByte(fixMapCodeMin + byte(c)); err != nil {
			m.tb.Fatalf("msgpack: failed to write element size prefix: %v", err)
		}
	case c < math.MaxUint16:
		if err := w.WriteByte(map16Code); err != nil {
			m.tb.Fatalf("msgpack: failed to write 16-bit element size prefix: %v", err)
		}
		b := make([]byte, 5)
		binary.BigEndian.PutUint16(b, uint16(c))
		if _, err := w.Write(b); err != nil {
			m.tb.Fatalf("msgpack: failed to write 16-bit element size: %v", err)
		}
	case c < math.MaxUint32:
		if err := w.WriteByte(map32Code); err != nil {
			m.tb.Fatalf("msgpack: failed to write 32-bit element size prefix: %v", err)
		}
		b := make([]byte, 7)
		binary.BigEndian.PutUint32(b, uint32(c))
		if _, err := w.Write(b); err != nil {
			m.tb.Fatalf("msgpack: failed to write 32-bit element size: %v", err)
		}
	default:
		m.tb.Fatalf("msgpack: map element count out of range (%d", c)
	}

	e := NewEncoder(w)
	for i := 0; i < c; i++ {
		if err := e.Encode(m.valuas[i*2]); err != nil {
			m.tb.Fatalf("msgpack: map builder: failed to encode map key %s", m.valuas[i])
		}

		if err := e.Encode(m.valuas[i*2+1]); err != nil {
			m.tb.Fatalf("msgpack: map builder: failed to encode map element for %s: %v", m.valuas[i], err)
		}
	}
}

func (m *testMapBuilder) Bytes() []byte {
	var buf bytes.Buffer
	m.encode(&buf)
	return buf.Bytes()
}
