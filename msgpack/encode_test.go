package msgpack

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type (
	typedString         string
	typedByte           byte
	typedByteSlice      []byte
	typedTypedByteSlice []typedByte
)

type (
	NotStruct          int
	AnonFieldNotStruct struct {
		NotStruct
	}
)

type ra struct {
	Sa string
	Ra *rb
}

type rb struct {
	Sb string
	Rb *ra
}

type me struct {
	s string
}

func (m me) MarshalMsgPack(enc *Encoder) error {
	return enc.PackString(m.s)
}

func TestEncode(t *testing.T) {
	t.Parallel()

	encodeTests := []struct {
		v    interface{}
		data []interface{}
	}{
		{
			v:    true,
			data: []interface{}{true},
		},
		{
			v:    1,
			data: []interface{}{(1)},
		},
		{
			v:    int8(2),
			data: []interface{}{2},
		},
		{
			v:    int16(3),
			data: []interface{}{3},
		},
		{
			v:    int32(4),
			data: []interface{}{4},
		},
		{
			v:    int64(5),
			data: []interface{}{5},
		},
		{
			v:    uint(1),
			data: []interface{}{1},
		},
		{
			v:    uint8(2),
			data: []interface{}{2},
		},
		{
			v:    uint16(3),
			data: []interface{}{3},
		},
		{
			v:    uint32(4),
			data: []interface{}{4},
		},
		{
			v:    uint64(5),
			data: []interface{}{5},
		},
		{
			v:    float32(5.0),
			data: []interface{}{5.0},
		},
		{
			v:    float64(6.0),
			data: []interface{}{6.0},
		},
		{
			v:    "hello",
			data: []interface{}{"hello"},
		},
		{
			v:    []byte("world"),
			data: []interface{}{[]byte("world")},
		},
		{
			v:    typedString("quux"),
			data: []interface{}{"quux"},
		},
		{
			v: typedByteSlice("foo"),
			data: []interface{}{
				[]byte("foo"),
			},
		},
		{
			v:    typedTypedByteSlice("bar"),
			data: []interface{}{[]byte("bar")},
		},
		{
			v:    []string(nil),
			data: []interface{}{nil},
		},
		{
			v:    []string{},
			data: []interface{}{arrayLen(0)},
		},
		{
			v: []string{"hello", "world"},
			data: []interface{}{
				arrayLen(2),
				"hello",
				"world",
			},
		},
		{
			v: [2]string{"hello", "world"},
			data: []interface{}{arrayLen(2),
				"hello",
				"world",
			},
		},
		{
			v:    map[string]string(nil),
			data: []interface{}{nil},
		},
		{
			v: map[string]string{"hello": "world"},
			data: []interface{}{
				mapLen(1),
				"hello",
				"world",
			},
		},
		{
			v:    new(int),
			data: []interface{}{0},
		},
		// Tag names
		{
			v: struct {
				A int
				B int `msgpack:"BB"`
				C int `msgpack:"omitempty"`
				D int `msgpack:"-"`
			}{
				A: 1,
				B: 2,
				C: 3,
				D: 4,
			},
			data: []interface{}{
				mapLen(3),
				"A", 1,
				"BB", 2,
				"omitempty", 3,
			},
		},
		// Struct as array
		{
			v: struct {
				I int `msgpack:",array"`
				S string
			}{
				I: 22,
				S: "skidoo",
			},
			data: []interface{}{arrayLen(2), 22, "skidoo"},
		},
		// Omit Empty
		{
			v: struct {
				B  bool `msgpack:"b,omitempty"`
				Bo bool `msgpack:"bo,omitempty"`
				Be bool `msgpack:"be,omitempty" empty:"true"`

				S  string `msgpack:"s,omitempty"`
				So string `msgpack:"so,omitempty"`
				Se string `msgpack:"se,omitempty" empty:"blank"`

				I  int `msgpack:"i,omitempty"`
				Io int `msgpack:"io,omitempty"`
				Ie int `msgpack:"ie,omitempty" empty:"-1"`

				U  uint `msgpack:"u,omitempty"`
				Uo uint `msgpack:"uo,omitempty"`

				F  float64 `msgpack:"f,omitempty"`
				Fo float64 `msgpack:"fo,omitempty"`

				D  float64 `msgpack:"d,omitempty"`
				Do float64 `msgpack:"do,omitempty"`

				Sl  []string `msgpack:"sl,omitempty"`
				Slo []string `msgpack:"slo,omitempty"`

				M  map[string]string `msgpack:"m,omitempty"`
				Mo map[string]string `msgpack:"mo,omitempty"`

				P  *int `msgpack:"p,omitempty"`
				Po *int `msgpack:"po,omitempty"`
			}{
				B:  false,
				Be: true,
				S:  "1",
				Se: "blank",
				I:  2,
				Ie: -1,
				U:  3,
				F:  4.0,
				D:  5.0,
				Sl: []string{"hello"},
				M:  map[string]string{"foo": "bar"},
				P:  new(int),
			},
			data: []interface{}{
				mapLen(8),
				"s", "1",
				"i", 2,
				"u", 3,
				"f", 4.0,
				"d", 5.0,
				"sl", arrayLen(1), "hello",
				"m", mapLen(1), "foo", "bar",
				"p", 0,
			},
		},
		{
			v: &ra{"foo", &rb{"bar", &ra{"quux", nil}}},
			data: []interface{}{
				mapLen(2),
				"Sa", "foo",
				"Ra", mapLen(2),
				"Sb", "bar",
				"Rb", mapLen(2),
				"Sa", "quux",
				"Ra", nil,
			},
		},
		{
			v:    &me{"hello"},
			data: []interface{}{"hello"},
		},
		{
			v:    []interface{}{"foo", "bar"},
			data: []interface{}{arrayLen(2), "foo", "bar"},
		},
		{
			v:    nil,
			data: []interface{}{nil},
		},
	}
	for _, tt := range encodeTests {
		tt := tt
		t.Run(fmt.Sprintf("%v", tt.v), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			enc := NewEncoder(&buf)
			if err := enc.Encode(tt.v); err != nil {
				t.Fatalf("encode %#v returned error %v", tt.v, err)
			}
			data, err := unpack(buf.Bytes())
			if err != nil {
				t.Fatalf("unpack %#v returned error %v", tt.v, err)
			}
			if !reflect.DeepEqual(data, tt.data) {
				t.Fatalf("encode %#v\n\t got: %#v\n\twant: %#v", tt.v, data, tt.data)
			}
		})
	}
}
