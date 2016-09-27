// Copyright 2016 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msgpack

import (
	"bytes"
	"reflect"
	"testing"
)

type renamedString string
type renamedByte byte
type renamedByteSlice []byte
type renamedRenamedByteSlice []renamedByte

type NotStruct int
type AnonFieldNotStruct struct {
	NotStruct
}

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

var encodeTests = []struct {
	v    interface{}
	data []interface{}
}{
	{true, []interface{}{true}},

	{1, []interface{}{(1)}},
	{int8(2), []interface{}{2}},
	{int16(3), []interface{}{3}},
	{int32(4), []interface{}{4}},
	{int64(5), []interface{}{5}},

	{uint(1), []interface{}{1}},
	{uint8(2), []interface{}{2}},
	{uint16(3), []interface{}{3}},
	{uint32(4), []interface{}{4}},
	{uint64(5), []interface{}{5}},

	{float32(5.0), []interface{}{5.0}},
	{float64(6.0), []interface{}{6.0}},

	{"hello", []interface{}{"hello"}},
	{[]byte("world"), []interface{}{[]byte("world")}},
	{renamedString("quux"), []interface{}{"quux"}},
	{renamedByteSlice("foo"), []interface{}{[]byte("foo")}},
	{renamedRenamedByteSlice("bar"), []interface{}{[]byte("bar")}},

	{[]string(nil), []interface{}{nil}},
	{[]string{}, []interface{}{arrayLen(0)}},
	{[]string{"hello", "world"}, []interface{}{arrayLen(2), "hello", "world"}},
	{[2]string{"hello", "world"}, []interface{}{arrayLen(2), "hello", "world"}},

	{map[string]string(nil), []interface{}{nil}},
	{map[string]string{"hello": "world"}, []interface{}{mapLen(1), "hello", "world"}},

	{new(int), []interface{}{0}},

	// Tag names
	{struct {
		A int
		B int `msgpack:"BB"`
		C int `msgpack:"omitempty"`
		D int `msgpack:"-"`
	}{1, 2, 3, 4}, []interface{}{mapLen(3), "A", 1, "BB", 2, "omitempty", 3}},

	// Struct as array
	{struct {
		I int `msgpack:",array"`
		S string
	}{22, "skidoo"}, []interface{}{arrayLen(2), 22, "skidoo"}},

	// Omit Empty
	{struct {
		B  bool `msgpack:"b,omitempty"`
		Bo bool `msgpack:"bo,omitempty"`

		S  string `msgpack:"s,omitempty"`
		So string `msgpack:"so,omitempty"`

		I  int `msgpack:"i,omitempty"`
		Io int `msgpack:"io,omitempty"`

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
		S:  "1",
		I:  2,
		U:  3,
		F:  4.0,
		D:  5.0,
		Sl: []string{"hello"},
		M:  map[string]string{"foo": "bar"},
		P:  new(int),
	},
		[]interface{}{
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
		&ra{"foo", &rb{"bar", &ra{"quux", nil}}},
		[]interface{}{
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
		&me{"hello"},
		[]interface{}{"hello"},
	},
	{
		[]interface{}{"foo", "bar"},
		[]interface{}{arrayLen(2), "foo", "bar"},
	},
	{
		nil,
		[]interface{}{nil},
	},
}

func TestEncode(t *testing.T) {
	for _, tt := range encodeTests {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		if err := enc.Encode(tt.v); err != nil {
			t.Errorf("encode %#v returned error %v", tt.v, err)
			continue
		}
		data, err := unpack(buf.Bytes())
		if err != nil {
			t.Errorf("unpack %#v returned error %v", tt.v, err)
			continue
		}
		if !reflect.DeepEqual(data, tt.data) {
			t.Errorf("encode %#v\n\t got: %#v\n\twant: %#v", tt.v, data, tt.data)
		}
	}
}
