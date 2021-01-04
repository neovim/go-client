package msgpack

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"
)

type testDecStruct struct {
	IF  interface{}
	B   bool
	S   string
	I   int
	U   uint
	F64 float64
	SS  []string
	M   map[string]interface{}
}

type testDecEmptyStruct struct {
	B   bool   `empty:"true"`
	S   string `empty:"blank"`
	I   int    `empty:"1234"`
	I8  int8   `empty:"45"`
	I32 int32  `empty:"6789"`
}

type testDecArrayStruct struct {
	I int `msgpack:",array"`
	S string
}

func ptrInt(i int) *int {
	return &i
}

func TestDecode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// arg is argument for Decode()
		arg func() interface{}
		// data is data to decode
		data []interface{}
		// expected is the expected decoded value
		expected interface{}
	}{
		"Bool/Bool/True": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{true},
			expected: true,
		},
		"Bool/Bool/False": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{false},
			expected: false,
		},
		"Bool/Int64/True": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{int64(1234)},
			expected: true,
		},
		"Bool/Int64/False": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{int64(0)},
			expected: false,
		},
		"Bool/Uint64/True": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{uint64(1234)},
			expected: true,
		},
		"Bool/Uint64/False": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{uint64(0)},
			expected: false,
		},
		"Int/Int64": {
			arg:      func() interface{} { return new(int) },
			data:     []interface{}{int64(1234)},
			expected: int(1234),
		},
		"Int/Uint64": {
			arg:      func() interface{} { return new(int) },
			data:     []interface{}{uint64(4321)},
			expected: int(4321),
		},
		"Int/Float64": {
			arg:      func() interface{} { return new(int) },
			data:     []interface{}{float64(5678)},
			expected: int(5678),
		},
		"Uint/Int64": {
			arg:      func() interface{} { return new(uint) },
			data:     []interface{}{int64(1234)},
			expected: uint(1234),
		},
		"Uint/Uint64": {
			arg:      func() interface{} { return new(uint) },
			data:     []interface{}{uint64(4321)},
			expected: uint(4321),
		},
		"Uint/Float64": {
			arg:      func() interface{} { return new(uint) },
			data:     []interface{}{float64(5678)},
			expected: uint(5678),
		},
		"Float64/Int64": {
			arg:      func() interface{} { return new(float64) },
			data:     []interface{}{int64(1234)},
			expected: float64(1234),
		},
		"Float64/Uint64": {
			arg:      func() interface{} { return new(float64) },
			data:     []interface{}{uint64(4321)},
			expected: float64(4321),
		},
		"Float64/Float64": {
			arg:      func() interface{} { return new(float64) },
			data:     []interface{}{float64(5678)},
			expected: float64(5678),
		},
		"String/Binary": {
			arg:      func() interface{} { return new(string) },
			data:     []interface{}{[]byte("world")},
			expected: "world",
		},
		"String/String": {
			arg:      func() interface{} { return new(string) },
			data:     []interface{}{"hello"},
			expected: "hello",
		},
		"Binary/Nil": {
			arg:      func() interface{} { return new([]byte) },
			data:     []interface{}{nil},
			expected: []byte(nil),
		},
		"Binary/Binary": {
			arg: func() interface{} { return new([]byte) },
			data: []interface{}{
				[]byte("hello"),
			},
			expected: []byte("hello"),
		},
		"Binary/String": {
			arg: func() interface{} { return new([]byte) },
			data: []interface{}{
				"world",
			},
			expected: []byte("world"),
		},
		"Interface/IntPointer": {
			arg: func() interface{} { return &testDecStruct{IF: ptrInt(1234)} },
			data: []interface{}{
				mapLen(1),
				"IF",
				int64(5678),
			},
			expected: testDecStruct{
				IF: ptrInt(5678),
			},
		},
		"Interface/StringSlice": {
			arg: func() interface{} { return &testDecStruct{IF: []string{"hello", "world"}} },
			data: []interface{}{
				mapLen(1),
				"IF",
				arrayLen(1),
				"foo",
			},
			expected: testDecStruct{
				IF: []string{"foo", ""},
			},
		},
		"StringSlice/ArrayLen/1": {
			arg: func() interface{} { return []string{""} },
			data: []interface{}{
				arrayLen(2),
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
			},
		},
		"StringSlice/ArrayLen/2/ValueValue": {
			arg: func() interface{} { return []string{"", ""} },
			data: []interface{}{
				arrayLen(2),
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
				"bar",
			},
		},
		"StringSlice/ArrayLen/2/ValueEmpty": {
			arg: func() interface{} { return []string{"", "bar"} },
			data: []interface{}{
				arrayLen(1),
				"foo",
			},
			expected: []string{
				"foo",
				"",
			},
		},
		"StringSlice/ArrayLen/Make/2": {
			arg: func() interface{} { x := make([]string, 1); return &x },
			data: []interface{}{
				arrayLen(2),
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
				"bar",
			},
		},
		"StringSlice/ArrayLen/Make/3": {
			arg: func() interface{} { x := make([]string, 3); return &x },
			data: []interface{}{
				arrayLen(2),
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
				"bar",
			},
		},
		"StringSlicePointer/ArrayLen/2": {
			arg: func() interface{} { return new([]string) },
			data: []interface{}{
				arrayLen(2),
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
				"bar",
			},
		},
		"StringArray/ArrayLen/3/ValueValueEmpty": {
			arg: func() interface{} { x := [...]string{"foo", "bar", "quux"}; return &x },
			data: []interface{}{
				arrayLen(2),
				"hello",
				"world",
			},
			expected: [...]string{
				"hello",
				"world",
				"",
			},
		},
		"StringArray/ArrayLen/1/Value": {
			arg: func() interface{} { x := [...]string{"foo"}; return &x },
			data: []interface{}{
				arrayLen(2),
				"hello",
				"world",
			},
			expected: [...]string{
				"hello",
			},
		},
		"StructArray/Int64": {
			arg: func() interface{} { return new(testDecArrayStruct) },
			data: []interface{}{
				arrayLen(2),
				int64(22),
				"skidoo",
			},
			expected: testDecArrayStruct{
				I: 22,
				S: "skidoo",
			},
		},
		"Map/StringString": {
			arg: func() interface{} { return make(map[string]string) },
			data: []interface{}{
				mapLen(1),
				"foo", "bar",
			},
			expected: map[string]string{
				"foo": "bar",
			},
		},
		"MapPointer/StringString": {
			arg: func() interface{} { return new(map[string]string) },
			data: []interface{}{
				mapLen(1),
				"foo", "bar",
			},
			expected: map[string]string{
				"foo": "bar",
			},
		},
		"Pointer/Int64": {
			arg: func() interface{} { return new(*int) },
			data: []interface{}{
				int64(-1),
			},
			expected: ptrInt(-1),
		},
		"Interface/Extensions/ExtensionValue": {
			arg: func() interface{} { return new(interface{}) },
			data: []interface{}{
				extension{
					0,
					"hello",
				},
			},
			expected: extensionValue{
				kind: 0,
				data: []byte("hello"),
			},
		},
		"Interface/Extensions/TestExtension": {
			arg: func() interface{} { return new(interface{}) },
			data: []interface{}{
				extension{
					1, "hello",
				},
			},
			expected: testExtension1{
				data: []byte("hello"),
			},
		},
		"TestExtension/Extensions": {
			arg: func() interface{} { return new(testExtension1) },
			data: []interface{}{
				extension{
					1, "hello",
				},
			},
			expected: testExtension1{
				data: []byte("hello"),
			},
		},
		"TestDecEmptyStruct/Empty/blank": {
			arg: func() interface{} { return &testDecEmptyStruct{} },
			data: []interface{}{
				mapLen(0),
			},
			expected: testDecEmptyStruct{
				B:   true,
				S:   "blank",
				I:   1234,
				I8:  45,
				I32: 6789,
			},
		},
		"TestDecEmptyStruct/Empty/NotBlank": {
			arg: func() interface{} { return &testDecEmptyStruct{} },
			data: []interface{}{
				mapLen(1),
				"S", "not blank",
			},
			expected: testDecEmptyStruct{
				B:   true,
				S:   "not blank",
				I:   1234,
				I8:  45,
				I32: 6789,
			},
		},
		// TODO(zchee): test errors like the following:
		// "Errors": {
		// 	arg:      func() interface{} { return &testDecStruct{I: 1234} },
		// 	data:     []interface{}{mapLen(1), "I", int64(5678)},
		// 	expected: testDecStruct{I: 1234},
		// },
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := pack(tt.data...)
			if err != nil {
				t.Fatalf("pack(%+v) returned error %v", tt.data, err)
			}
			dec := NewDecoder(bytes.NewReader(data))
			buf, _ := dec.r.Peek(0)

			dec.SetExtensions(testExtensionMap)

			arg := tt.arg()
			if err := dec.Decode(arg); err != nil {
				t.Fatalf("decode(%+v, %T) returned error %v", tt.data, arg, err)
			}

			// scribble on bufio.Reader buffer to test that Decoder.Bytes() return value is copied
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
				t.Fatalf("decode(%+v, %T) returned %#v, want %#v", tt.data, arg, v, tt.expected)
			}

			// Decode should read to EOF.
			if _, err := dec.r.ReadByte(); err != io.EOF {
				t.Fatalf("decode(%+v, %T) did not read to EOF", tt.data, arg)
			}
		})
	}
}

func Test_boolDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds   *decodeState
		want bool
	}{
		"Bool/True": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Bool,
				},
			},
			want: true,
		},
		"Bool/False": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(0),
					t: Bool,
				},
			},
			want: false,
		},
		"Int/True": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Int,
				},
			},
			want: true,
		},
		"Int/False": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(0),
					t: Int,
				},
			},
			want: false,
		},
		"Uint/True": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Uint,
				},
			},
			want: true,
		},
		"Uint/False": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(0),
					t: Uint,
				},
			},
			want: false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(bool)).Elem()
			boolDecoder(tt.ds, v)
			if got := tt.ds.Bool(); got != tt.want {
				t.Fatalf("boolDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}
		})
	}
}

func Test_intDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds   *decodeState
		want int64
	}{
		"Int": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1234),
					t: Int,
				},
			},
			want: int64(1234),
		},
		"Uint": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(4321),
					t: Uint,
				},
			},
			want: int64(4321),
		},
		"Float": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(5678),
					t: Float,
				},
			},
			want: int64(5678),
		},
		// TODO(zchee): default case
		// "": {
		// 	ds: &decodeState{
		// 		Decoder: &Decoder{
		// 			n: uint64(0),
		// 			t: Invalid,
		// 		},
		// 	},
		// 	want: int64(0),
		// },
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(int64)).Elem()
			intDecoder(tt.ds, v)
			if got := tt.ds.Int(); got != tt.want {
				t.Fatalf("intDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}
		})
	}
}

func Test_uintDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds   *decodeState
		want uint64
	}{
		"Uint": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1234),
					t: Uint,
				},
			},
			want: uint64(1234),
		},
		"Int": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(4321),
					t: Int,
				},
			},
			want: uint64(4321),
		},
		"Float": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(5678),
					t: Float,
				},
			},
			want: uint64(5678),
		},
		// TODO(zchee): default case
		// "": {
		// 	ds: &decodeState{
		// 		Decoder: &Decoder{
		// 			n: uint64(0),
		// 			t: Invalid,
		// 		},
		// 	},
		// 	want: uint64(0),
		// },
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(uint64)).Elem()
			uintDecoder(tt.ds, v)
			if got := tt.ds.Uint(); got != tt.want {
				t.Fatalf("uintDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}
		})
	}
}

func Test_floatDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds   *decodeState
		want float64
	}{
		"Int": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1234),
					t: Int,
				},
			},
			want: math.Float64frombits(uint64(1234)),
		},
		"Uint": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(4321),
					t: Uint,
				},
			},
			want: math.Float64frombits(uint64(4321)),
		},
		"Float": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(5678),
					t: Float,
				},
			},
			want: math.Float64frombits(uint64(5678)),
		},
		// TODO(zchee): default case
		// "": {
		// 	ds: &decodeState{
		// 		Decoder: &Decoder{
		// 			n: uint64(0),
		// 			t: Invalid,
		// 		},
		// 	},
		// 	want: float64(0),
		// },
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(float64)).Elem()
			floatDecoder(tt.ds, v)
			if got := tt.ds.Float(); got != tt.want {
				t.Fatalf("floatDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}
		})
	}
}

func Test_stringDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds   *decodeState
		want string
	}{
		"Binary": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("hello"),
					t: Binary,
				},
			},
			want: string("hello"),
		},
		"String": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("world"),
					t: String,
				},
			},
			want: string("world"),
		},
		// TODO(zchee): default case
		// "": {
		// 	ds: &decodeState{
		// 		Decoder: &Decoder{
		// 			p: []byte(nil),
		// 			t: Invalid,
		// 		},
		// 	},
		// 	want: string(""),
		// },
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(string)).Elem()
			stringDecoder(tt.ds, v)
			if got := tt.ds.String(); got != tt.want {
				t.Fatalf("stringDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}
		})
	}
}

func Test_byteSliceDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds   *decodeState
		want []byte
	}{
		"Binary": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("hello"),
					t: Binary,
				},
			},
			want: []byte("hello"),
		},
		"String": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("world"),
					t: String,
				},
			},
			want: []byte("world"),
		},
		// TODO(zchee): default case
		// "": {
		// 	ds: &decodeState{
		// 		Decoder: &Decoder{
		// 			p: []byte(nil),
		// 			t: Invalid,
		// 		},
		// 	},
		// 	want: []byte(""),
		// },
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new([]byte)).Elem()
			byteSliceDecoder(tt.ds, v)
			if got := tt.ds.Bytes(); !bytes.Equal(got, tt.want) {
				t.Fatalf("byteSliceDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}
		})
	}
}
