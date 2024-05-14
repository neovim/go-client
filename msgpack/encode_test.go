package msgpack

import (
	"bytes"
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

	encodeTests := map[string]struct {
		v    any
		data []any
	}{
		"Bool": {
			v:    true,
			data: []any{true},
		},
		"Int": {
			v:    int(1),
			data: []any{(1)},
		},
		"Int8": {
			v:    int8(2),
			data: []any{2},
		},
		"Int16": {
			v:    int16(3),
			data: []any{3},
		},
		"Int32": {
			v:    int32(4),
			data: []any{4},
		},
		"Int64": {
			v:    int64(5),
			data: []any{5},
		},
		"Uint": {
			v:    uint(1),
			data: []any{1},
		},
		"Uint8": {
			v:    uint8(2),
			data: []any{2},
		},
		"Uint16": {
			v:    uint16(3),
			data: []any{3},
		},
		"Uint32": {
			v:    uint32(4),
			data: []any{4},
		},
		"Uint64": {
			v:    uint64(5),
			data: []any{5},
		},
		"Float32": {
			v:    float32(5.0),
			data: []any{5.0},
		},
		"Float64": {
			v:    float64(6.0),
			data: []any{6.0},
		},
		"String": {
			v:    "hello",
			data: []any{"hello"},
		},
		"Bytes": {
			v:    []byte("world"),
			data: []any{[]byte("world")},
		},
		"typedString": {
			v:    typedString("quux"),
			data: []any{"quux"},
		},
		"typedByteSlice": {
			v: typedByteSlice("foo"),
			data: []any{
				[]byte("foo"),
			},
		},
		"typedTypedByteSlice": {
			v:    typedTypedByteSlice("bar"),
			data: []any{[]byte("bar")},
		},
		"StringSlice/Nil": {
			v:    []string(nil),
			data: []any{nil},
		},
		"StringSlice/Empty": {
			v:    []string{},
			data: []any{arrayLen(0)},
		},
		"StringSlice/Value": {
			v: []string{"hello", "world"},
			data: []any{
				arrayLen(2),
				"hello",
				"world",
			},
		},
		"StringArray/Value": {
			v: [2]string{"hello", "world"},
			data: []any{
				arrayLen(2),
				"hello",
				"world",
			},
		},
		"MapStringString/Nil": {
			v:    map[string]string(nil),
			data: []any{nil},
		},
		"MapStringString/Value": {
			v: map[string]string{"hello": "world"},
			data: []any{
				mapLen(1),
				"hello",
				"world",
			},
		},
		"IntPointer": {
			v:    new(int),
			data: []any{0},
		},
		"TagNames": {
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
			data: []any{
				mapLen(3),
				"A", 1,
				"BB", 2,
				"omitempty", 3,
			},
		},
		"StructAsArray": {
			v: struct {
				I int `msgpack:",array"`
				S string
			}{
				I: 22,
				S: "skidoo",
			},
			data: []any{arrayLen(2), 22, "skidoo"},
		},
		"OmitEmpty": {
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
			data: []any{
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
		"Struct": {
			v: &ra{"foo", &rb{"bar", &ra{"quux", nil}}},
			data: []any{
				mapLen(2),
				"Sa", "foo",
				"Ra", mapLen(2),
				"Sb", "bar",
				"Rb", mapLen(2),
				"Sa", "quux",
				"Ra", nil,
			},
		},
		"MarshalMsgPack": {
			v:    &me{"hello"},
			data: []any{"hello"},
		},
		"Interface": {
			v:    []any{"foo", "bar"},
			data: []any{arrayLen(2), "foo", "bar"},
		},
		"Nil": {
			v:    nil,
			data: []any{nil},
		},
	}
	for name, tt := range encodeTests {
		tt := tt
		t.Run(name, func(t *testing.T) {
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
