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
	"io"
	"reflect"
	"testing"
)

type testDecStruct struct {
	I  interface{}
	S  string
	N  int
	U  uint
	F  float64
	Sl []string
	M  map[string]interface{}
}

type testDecArrayStruct struct {
	I int `msgpack:",array"`
	S string
}

func ptrInt(i int) *int {
	return &i
}

var decodeTests = []struct {
	// arg is argument for Decode().
	arg func() interface{}
	// data is data to decode.
	data []interface{}
	// expected is the expected decoded value.
	expected interface{}
}{
	// int
	{func() interface{} { return new(int) }, []interface{}{int64(1234)}, int(1234)},
	{func() interface{} { return new(int) }, []interface{}{float64(4321)}, int(4321)},
	{func() interface{} { return new(int) }, []interface{}{uint64(5678)}, int(5678)},

	// uint
	{func() interface{} { return new(uint) }, []interface{}{int64(1234)}, uint(1234)},
	{func() interface{} { return new(uint) }, []interface{}{float64(4321)}, uint(4321)},
	{func() interface{} { return new(uint) }, []interface{}{uint64(5678)}, uint(5678)},

	// float
	{func() interface{} { return new(float64) }, []interface{}{int64(1234)}, float64(1234)},
	{func() interface{} { return new(float64) }, []interface{}{float64(4321)}, float64(4321)},
	{func() interface{} { return new(float64) }, []interface{}{uint64(5678)}, float64(5678)},

	// bool
	{func() interface{} { return new(bool) }, []interface{}{true}, true},
	{func() interface{} { return new(bool) }, []interface{}{false}, false},

	// string
	{func() interface{} { return new(string) }, []interface{}{"hello"}, "hello"},
	{func() interface{} { return new(string) }, []interface{}{[]byte("world")}, "world"},

	// []byte
	{func() interface{} { return new([]byte) }, []interface{}{"hello"}, []byte("hello")},
	{func() interface{} { return new([]byte) }, []interface{}{[]byte("world")}, []byte("world")},

	// Pointer
	{func() interface{} { return new(*int) }, []interface{}{int64(-1)}, ptrInt(-1)},

	// Interface
	// *int is set
	{func() interface{} { return &testDecStruct{I: ptrInt(1234)} }, []interface{}{mapLen(1), "I", int64(5678)}, testDecStruct{I: ptrInt(5678)}},
	// []string elements set, but not resized.
	{func() interface{} { return &testDecStruct{I: []string{"hello", "world"}} }, []interface{}{mapLen(1), "I", arrayLen(1), "foo"}, testDecStruct{I: []string{"foo", ""}}},

	// *Slice
	{func() interface{} { return new([]string) }, []interface{}{arrayLen(2), "foo", "bar"}, []string{"foo", "bar"}},
	{func() interface{} { x := make([]string, 1); return &x }, []interface{}{arrayLen(2), "foo", "bar"}, []string{"foo", "bar"}},
	{func() interface{} { x := make([]string, 3); return &x }, []interface{}{arrayLen(2), "foo", "bar"}, []string{"foo", "bar"}},

	// Slice
	{func() interface{} { return []string{"", ""} }, []interface{}{arrayLen(2), "foo", "bar"}, []string{"foo", "bar"}},
	{func() interface{} { return []string{""} }, []interface{}{arrayLen(2), "foo", "bar"}, []string{"foo"}},
	{func() interface{} { return []string{"", "bar"} }, []interface{}{arrayLen(1), "foo"}, []string{"foo", ""}},

	// Array
	{func() interface{} { x := [...]string{"foo", "bar", "quux"}; return &x }, []interface{}{arrayLen(2), "hello", "world"}, [...]string{"hello", "world", ""}},
	{func() interface{} { x := [...]string{"foo"}; return &x }, []interface{}{arrayLen(2), "hello", "world"}, [...]string{"hello"}},

	// Struct array
	{func() interface{} { return new(testDecArrayStruct) }, []interface{}{arrayLen(2), int64(22), "skidoo"}, testDecArrayStruct{22, "skidoo"}},

	// Map
	{func() interface{} { return make(map[string]string) }, []interface{}{mapLen(1), "foo", "bar"}, map[string]string{"foo": "bar"}},

	// *Map
	{func() interface{} { return new(map[string]string) }, []interface{}{mapLen(1), "foo", "bar"}, map[string]string{"foo": "bar"}},

	// TODO: test errors like the following:
	// {func() interface{} { return &testDecStruct{I: 1234} }, []interface{}{mapLen(1), "I", int64(5678)}, testDecStruct{I: 1234}},
}

func TestDecode(t *testing.T) {
	for _, tt := range decodeTests {
		data, err := pack(tt.data...)
		if err != nil {
			t.Errorf("pack(%+v) returned error %v", tt.data, err)
			continue
		}
		dec := NewDecoder(bytes.NewReader(data))
		buf, _ := dec.r.Peek(0)

		arg := tt.arg()
		if err := dec.Decode(arg); err != nil {
			t.Errorf("decode(%+v, %T) returned error %v", tt.data, arg, err)
			continue
		}

		// scribble on bufio.Reader buffer to test that Decoder.Bytes() return value is copied.
		buf = buf[:cap(buf)]
		for i := range buf {
			buf[i] = 0xff
		}

		rv := reflect.ValueOf(arg)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		v := rv.Interface()
		if !reflect.DeepEqual(v, tt.expected) {
			t.Errorf("decode(%+v, %T) returned %#v, want %#v", tt.data, arg, v, tt.expected)
		}

		// Decode should read to EOF.
		if _, err := dec.r.ReadByte(); err != io.EOF {
			t.Errorf("decode(%+v, %T) did not read to EOF", tt.data, arg)
		}
	}
}
