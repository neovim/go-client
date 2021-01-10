package msgpack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"
)

type testDecStruct struct {
	B   bool
	I   int
	U   uint
	F64 float64
	S   string
	SS  []string
	M   map[string]interface{}
	IF  interface{}
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

func TestDecode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		// arg is argument for Decode()
		arg func() interface{}
		// data is data to decode
		data []interface{}
		// expected is the expected decoded value
		expected interface{}
		// wantErr is the whether the want error
		wantErr bool
	}{
		"Bool/Bool/True": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{true},
			expected: true,
			wantErr:  false,
		},
		"Bool/Bool/False": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{false},
			expected: false,
			wantErr:  false,
		},
		"Bool/Int64/True": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{int64(1234)},
			expected: true,
			wantErr:  false,
		},
		"Bool/Int64/False": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{int64(0)},
			expected: false,
			wantErr:  false,
		},
		"Bool/Uint64/True": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{uint64(1234)},
			expected: true,
			wantErr:  false,
		},
		"Bool/Uint64/False": {
			arg:      func() interface{} { return new(bool) },
			data:     []interface{}{uint64(0)},
			expected: false,
			wantErr:  false,
		},
		"Int/Int64": {
			arg:      func() interface{} { return new(int) },
			data:     []interface{}{int64(1234)},
			expected: int(1234),
			wantErr:  false,
		},
		"Int/Uint64": {
			arg:      func() interface{} { return new(int) },
			data:     []interface{}{uint64(4321)},
			expected: int(4321),
			wantErr:  false,
		},
		"Int/Float64": {
			arg:      func() interface{} { return new(int) },
			data:     []interface{}{float64(5678)},
			expected: int(5678),
			wantErr:  false,
		},
		"Uint/Int64": {
			arg:      func() interface{} { return new(uint) },
			data:     []interface{}{int64(1234)},
			expected: uint(1234),
			wantErr:  false,
		},
		"Uint/Uint64": {
			arg:      func() interface{} { return new(uint) },
			data:     []interface{}{uint64(4321)},
			expected: uint(4321),
			wantErr:  false,
		},
		"Uint/Float64": {
			arg:      func() interface{} { return new(uint) },
			data:     []interface{}{float64(5678)},
			expected: uint(5678),
			wantErr:  false,
		},
		"Float64/Int64": {
			arg:      func() interface{} { return new(float64) },
			data:     []interface{}{int64(1234)},
			expected: float64(1234),
			wantErr:  false,
		},
		"Float64/Uint64": {
			arg:      func() interface{} { return new(float64) },
			data:     []interface{}{uint64(4321)},
			expected: float64(4321),
			wantErr:  false,
		},
		"Float64/Float64": {
			arg:      func() interface{} { return new(float64) },
			data:     []interface{}{float64(5678)},
			expected: float64(5678),
			wantErr:  false,
		},
		"String/Binary": {
			arg:      func() interface{} { return new(string) },
			data:     []interface{}{[]byte("world")},
			expected: "world",
			wantErr:  false,
		},
		"String/String": {
			arg:      func() interface{} { return new(string) },
			data:     []interface{}{"hello"},
			expected: "hello",
			wantErr:  false,
		},
		"Binary/Nil": {
			arg:      func() interface{} { return new([]byte) },
			data:     []interface{}{nil},
			expected: []byte(nil),
			wantErr:  false,
		},
		"Binary/Binary": {
			arg: func() interface{} { return new([]byte) },
			data: []interface{}{
				[]byte("hello"),
			},
			expected: []byte("hello"),
			wantErr:  false,
		},
		"Binary/String": {
			arg: func() interface{} { return new([]byte) },
			data: []interface{}{
				"world",
			},
			expected: []byte("world"),
			wantErr:  false,
		},
		"StringSlice/ArrayLen/2": {
			arg: func() interface{} { return []string{""} },
			data: []interface{}{
				arrayLen(2),
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
			},
			wantErr: false,
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
			wantErr: false,
		},
		"StringSlice/ArrayLen/1/ValueEmpty": {
			arg: func() interface{} { return []string{"", "bar"} },
			data: []interface{}{
				arrayLen(1),
				"foo",
			},
			expected: []string{
				"foo",
				"",
			},
			wantErr: false,
		},
		"StringSlice/ArrayLen/2/Make/1": {
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
			wantErr: false,
		},
		"StringSlice/ArrayLen/2/Make/3": {
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
			wantErr: false,
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
			wantErr: false,
		},
		"StringArray/ArrayLen/2/ValueValueEmpty": {
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
			wantErr: false,
		},
		"StringArray/ArrayLen/2/Value": {
			arg: func() interface{} { x := [...]string{"foo"}; return &x },
			data: []interface{}{
				arrayLen(2),
				"hello",
				"world",
			},
			expected: [...]string{
				"hello",
			},
			wantErr: false,
		},
		"StructArray/ArrayLen/Int": {
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
		},
		"Pointer/Int": {
			arg: func() interface{} { return new(*int) },
			data: []interface{}{
				int64(-1),
			},
			expected: ptrInt(-1),
			wantErr:  false,
		},
		"Pointer/Uint": {
			arg: func() interface{} { return new(*uint) },
			data: []interface{}{
				uint64(1),
			},
			expected: ptrUint(1),
			wantErr:  false,
		},
		"Struct/Bool": {
			arg: func() interface{} { return &testDecStruct{B: true} },
			data: []interface{}{
				mapLen(1),
				"B",
				bool(true),
			},
			expected: testDecStruct{
				B: bool(true),
			},
			wantErr: false,
		},
		"Struct/Int": {
			arg: func() interface{} { return &testDecStruct{I: int(1234)} },
			data: []interface{}{
				mapLen(1),
				"I",
				int64(1234),
			},
			expected: testDecStruct{
				I: int(1234),
			},
			wantErr: false,
		},
		"Struct/Uint": {
			arg: func() interface{} { return &testDecStruct{U: uint(1234)} },
			data: []interface{}{
				mapLen(1),
				"U",
				uint64(1234),
			},
			expected: testDecStruct{
				U: uint(1234),
			},
			wantErr: false,
		},
		"Struct/Float": {
			arg: func() interface{} { return &testDecStruct{F64: float64(1234)} },
			data: []interface{}{
				mapLen(1),
				"F64",
				float64(1234),
			},
			expected: testDecStruct{
				F64: float64(1234),
			},
			wantErr: false,
		},
		"Struct/String": {
			arg: func() interface{} { return &testDecStruct{S: "hello"} },
			data: []interface{}{
				mapLen(1),
				"S",
				"hello",
			},
			expected: testDecStruct{
				S: "hello",
			},
			wantErr: false,
		},
		"Struct/StringSlice": {
			arg: func() interface{} { return &testDecStruct{SS: []string{"hello", "world"}} },
			data: []interface{}{
				mapLen(1),
				"SS",
				arrayLen(2),
				"hello", "world",
			},
			expected: testDecStruct{
				SS: []string{"hello", "world"},
			},
			wantErr: false,
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
			wantErr: false,
		},
		"Interface/UintPointer": {
			arg: func() interface{} { return &testDecStruct{IF: ptrUint(1234)} },
			data: []interface{}{
				mapLen(1),
				"IF",
				uint64(5678),
			},
			expected: testDecStruct{
				IF: ptrUint(5678),
			},
			wantErr: false,
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
			wantErr: false,
		},
		"Interface/NoReflect/Bool": {
			arg:      func() interface{} { return new(interface{}) },
			data:     []interface{}{bool(true)},
			expected: bool(true),
			wantErr:  false,
		},
		"Interface/NoReflect/Float64": {
			arg:      func() interface{} { return new(interface{}) },
			data:     []interface{}{float64(1234)},
			expected: float64(1234),
			wantErr:  false,
		},
		"Interface/NoReflect/Bytes": {
			arg:      func() interface{} { return new(interface{}) },
			data:     []interface{}{[]byte("hello")},
			expected: []byte("hello"),
			wantErr:  false,
		},
		"Interface/NoReflect/String": {
			arg:      func() interface{} { return new(interface{}) },
			data:     []interface{}{"hello"},
			expected: string("hello"),
			wantErr:  false,
		},
		"Interface/Extensions/ExtensionValue": {
			arg: func() interface{} { return new(interface{}) },
			data: []interface{}{
				extension{
					0, "hello",
				},
			},
			expected: extensionValue{
				kind: 0,
				data: []byte("hello"),
			},
			wantErr: false,
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
			wantErr: false,
		},
		"Extensions/TestExtension": {
			arg: func() interface{} { return new(testExtension1) },
			data: []interface{}{
				extension{
					1, "hello",
				},
			},
			expected: testExtension1{
				data: []byte("hello"),
			},
			wantErr: false,
		},
		"Empty/blank": {
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
			wantErr: false,
		},
		"Empty/NotBlank": {
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
			wantErr: false,
		},
		"Error/Bool": {
			arg:      func() interface{} { return &testDecStruct{B: true} },
			data:     []interface{}{mapLen(1), "B", false},
			expected: testDecStruct{B: true},
			wantErr:  true,
		},
		"Error/Int": {
			arg:      func() interface{} { return &testDecStruct{I: int(1234)} },
			data:     []interface{}{mapLen(1), "I", int64(5678)},
			expected: testDecStruct{I: int(1234)},
			wantErr:  true,
		},
		"Error/Uint": {
			arg:      func() interface{} { return &testDecStruct{U: uint(1234)} },
			data:     []interface{}{mapLen(1), "U", uint64(5678)},
			expected: testDecStruct{U: uint(1234)},
			wantErr:  true,
		},
		"Error/Float": {
			arg:      func() interface{} { return &testDecStruct{F64: float64(1234)} },
			data:     []interface{}{mapLen(1), "F64", float64(5678)},
			expected: testDecStruct{F64: float64(1234)},
			wantErr:  true,
		},
		"Error/String": {
			arg:      func() interface{} { return &testDecStruct{S: "hello"} },
			data:     []interface{}{mapLen(1), "S", "world"},
			expected: testDecStruct{S: "hello"},
			wantErr:  true,
		},
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
			if !reflect.DeepEqual(v, tt.expected) != tt.wantErr {
				t.Fatalf("decode(%+v, %T) returned %T:%#v, want %#v", tt.data, arg, arg, v, tt.expected)
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
		ds      *decodeState
		want    bool
		wantErr bool
	}{
		"Bool/True": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Bool,
				},
			},
			want:    true,
			wantErr: false,
		},
		"Bool/False": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(0),
					t: Bool,
				},
			},
			want:    false,
			wantErr: false,
		},
		"Int/True": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Int,
				},
			},
			want:    true,
			wantErr: false,
		},
		"Int/False": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(0),
					t: Int,
				},
			},
			want:    false,
			wantErr: false,
		},
		"Uint/True": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Uint,
				},
			},
			want:    true,
			wantErr: false,
		},
		"Uint/False": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(0),
					t: Uint,
				},
			},
			want:    false,
			wantErr: false,
		},
		"Invalid": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1),
					t: Invalid,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(bool)).Elem()
			boolDecoder(tt.ds, v)

			if got := tt.ds.Bool(); (got != tt.want) != tt.wantErr {
				t.Fatalf("boolDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}

			if (tt.ds.errSaved != nil) != tt.wantErr {
				t.Fatalf("expected tt.ds.errSaved is not nil: %#v", tt.ds.errSaved)
			}

			if tt.wantErr {
				var decConvErr *DecodeConvertError
				if !errors.As(tt.ds.errSaved, &decConvErr) {
					t.Fatalf("tt.ds.errSaved is not DecodeConvertError: %#v", tt.ds.errSaved)
				}

				expected := fmt.Sprintf("msgpack: cannot convert %s to %s", decConvErr.SrcType, decConvErr.DestType)
				if decConvErr.SrcValue != nil {
					expected = fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", decConvErr.SrcType, decConvErr.SrcValue, decConvErr.DestType)
				}
				if got := decConvErr.Error(); got != expected {
					t.Fatalf("decConvErr.Error() = %s: want: %s", got, expected)
				}
			}
		})
	}
}

func Test_intDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds      *decodeState
		want    int64
		wantErr bool
	}{
		"Int": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1234),
					t: Int,
				},
			},
			want:    int64(1234),
			wantErr: false,
		},
		"Uint": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(4321),
					t: Uint,
				},
			},
			want:    int64(4321),
			wantErr: false,
		},
		"Float": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: math.Float64bits(float64(5678)),
					t: Float,
				},
			},
			want:    int64(math.Float64bits(float64(5678))),
			wantErr: false,
		},
		"Invalid": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(8765),
					t: Invalid,
				},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(int64)).Elem()
			intDecoder(tt.ds, v)

			if got := tt.ds.Int(); (got != tt.want) != tt.wantErr {
				t.Fatalf("intDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}

			if (tt.ds.errSaved != nil) != tt.wantErr {
				t.Fatalf("expected tt.ds.errSaved is not nil: %#v", tt.ds.errSaved)
			}

			if tt.wantErr {
				var decConvErr *DecodeConvertError
				if !errors.As(tt.ds.errSaved, &decConvErr) {
					t.Fatalf("tt.ds.errSaved is not DecodeConvertError: %#v", tt.ds.errSaved)
				}

				expected := fmt.Sprintf("msgpack: cannot convert %s to %s", decConvErr.SrcType, decConvErr.DestType)
				if decConvErr.SrcValue != nil {
					expected = fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", decConvErr.SrcType, decConvErr.SrcValue, decConvErr.DestType)
				}
				if got := decConvErr.Error(); got != expected {
					t.Fatalf("decConvErr.Error() = %s: want: %s", got, expected)
				}
			}
		})
	}
}

func Test_uintDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds      *decodeState
		want    uint64
		wantErr bool
	}{
		"Uint": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1234),
					t: Uint,
				},
			},
			want:    uint64(1234),
			wantErr: false,
		},
		"Int": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(4321),
					t: Int,
				},
			},
			want:    uint64(4321),
			wantErr: false,
		},
		"Float": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: math.Float64bits(float64(5678)),
					t: Float,
				},
			},
			want:    math.Float64bits(float64(5678)),
			wantErr: false,
		},
		"Invalid": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(8765),
					t: Invalid,
				},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(uint64)).Elem()
			uintDecoder(tt.ds, v)

			if got := tt.ds.Uint(); (got != tt.want) != tt.wantErr {
				t.Fatalf("uintDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}

			if (tt.ds.errSaved != nil) != tt.wantErr {
				t.Fatalf("expected tt.ds.errSaved is not nil: %#v", tt.ds.errSaved)
			}

			if tt.wantErr {
				var decConvErr *DecodeConvertError
				if !errors.As(tt.ds.errSaved, &decConvErr) {
					t.Fatalf("tt.ds.errSaved is not DecodeConvertError: %#v", tt.ds.errSaved)
				}

				expected := fmt.Sprintf("msgpack: cannot convert %s to %s", decConvErr.SrcType, decConvErr.DestType)
				if decConvErr.SrcValue != nil {
					expected = fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", decConvErr.SrcType, decConvErr.SrcValue, decConvErr.DestType)
				}
				if got := decConvErr.Error(); got != expected {
					t.Fatalf("decConvErr.Error() = %s: want: %s", got, expected)
				}
			}
		})
	}
}

func Test_floatDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds      *decodeState
		want    float64
		wantErr bool
	}{
		"Int": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(1234),
					t: Int,
				},
			},
			want:    math.Float64frombits(uint64(1234)),
			wantErr: false,
		},
		"Uint": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(4321),
					t: Uint,
				},
			},
			want:    math.Float64frombits(uint64(4321)),
			wantErr: false,
		},
		"Float": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(5678),
					t: Float,
				},
			},
			want:    math.Float64frombits(uint64(5678)),
			wantErr: false,
		},
		"Invalid": {
			ds: &decodeState{
				Decoder: &Decoder{
					n: uint64(8765),
					t: Invalid,
				},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(float64)).Elem()
			floatDecoder(tt.ds, v)

			if got := tt.ds.Float(); (got != tt.want) != tt.wantErr {
				t.Fatalf("floatDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}

			if (tt.ds.errSaved != nil) != tt.wantErr {
				t.Fatalf("expected tt.ds.errSaved is not nil: %#v", tt.ds.errSaved)
			}

			if tt.wantErr {
				var decConvErr *DecodeConvertError
				if !errors.As(tt.ds.errSaved, &decConvErr) {
					t.Fatalf("tt.ds.errSaved is not DecodeConvertError: %#v", tt.ds.errSaved)
				}

				expected := fmt.Sprintf("msgpack: cannot convert %s to %s", decConvErr.SrcType, decConvErr.DestType)
				if decConvErr.SrcValue != nil {
					expected = fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", decConvErr.SrcType, decConvErr.SrcValue, decConvErr.DestType)
				}
				if got := decConvErr.Error(); got != expected {
					t.Fatalf("decConvErr.Error() = %s: want: %s", got, expected)
				}
			}
		})
	}
}

func Test_stringDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds      *decodeState
		want    string
		wantErr bool
	}{
		"Binary": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("hello"),
					t: Binary,
				},
			},
			want:    string("hello"),
			wantErr: false,
		},
		"String": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("world"),
					t: String,
				},
			},
			want:    string("world"),
			wantErr: false,
		},
		"Invalid": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("invalid"),
					t: Invalid,
				},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new(string)).Elem()
			stringDecoder(tt.ds, v)

			if got := tt.ds.String(); (got != tt.want) != tt.wantErr {
				t.Fatalf("stringDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}

			if (tt.ds.errSaved != nil) != tt.wantErr {
				t.Fatalf("expected tt.ds.errSaved is not nil: %#v", tt.ds.errSaved)
			}

			if tt.wantErr {
				var decConvErr *DecodeConvertError
				if !errors.As(tt.ds.errSaved, &decConvErr) {
					t.Fatalf("tt.ds.errSaved is not DecodeConvertError: %#v", tt.ds.errSaved)
				}

				expected := fmt.Sprintf("msgpack: cannot convert %s to %s", decConvErr.SrcType, decConvErr.DestType)
				if decConvErr.SrcValue != nil {
					expected = fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", decConvErr.SrcType, decConvErr.SrcValue, decConvErr.DestType)
				}
				if got := decConvErr.Error(); got != expected {
					t.Fatalf("decConvErr.Error() = %s: want: %s", got, expected)
				}
			}
		})
	}
}

func Test_byteSliceDecoder(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ds      *decodeState
		want    []byte
		wantErr bool
	}{
		"Binary": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("hello"),
					t: Binary,
				},
			},
			want:    []byte("hello"),
			wantErr: false,
		},
		"String": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("world"),
					t: String,
				},
			},
			want:    []byte("world"),
			wantErr: false,
		},
		"Invalid": {
			ds: &decodeState{
				Decoder: &Decoder{
					p: []byte("invalid"),
					t: Invalid,
				},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.ValueOf(new([]byte)).Elem()
			byteSliceDecoder(tt.ds, v)

			if got := tt.ds.Bytes(); !bytes.Equal(got, tt.want) != tt.wantErr {
				t.Fatalf("byteSliceDecoder(%v, %v) = %v: want: %v", tt.ds, v, got, tt.want)
			}

			if (tt.ds.errSaved != nil) != tt.wantErr {
				t.Fatalf("expected tt.ds.errSaved is not nil: %#v", tt.ds.errSaved)
			}

			if tt.wantErr {
				var decConvErr *DecodeConvertError
				if !errors.As(tt.ds.errSaved, &decConvErr) {
					t.Fatalf("tt.ds.errSaved is not DecodeConvertError: %#v", tt.ds.errSaved)
				}

				expected := fmt.Sprintf("msgpack: cannot convert %s to %s", decConvErr.SrcType, decConvErr.DestType)
				if decConvErr.SrcValue != nil {
					expected = fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", decConvErr.SrcType, decConvErr.SrcValue, decConvErr.DestType)
				}
				if got := decConvErr.Error(); got != expected {
					t.Fatalf("decConvErr.Error() = %s: want: %s", got, expected)
				}
			}
		})
	}
}
