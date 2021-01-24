package msgpack

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()

	benchs := map[string]struct {
		value interface{}
	}{
		"Bool/False": {
			value: false,
		},
		"Bool/True": {
			value: true,
		},
		"Uint8": {
			value: math.MaxUint8,
		},
		"Uint16": {
			value: math.MaxUint16,
		},
		"Uint32": {
			value: math.MaxUint32,
		},
		"Uint64": {
			value: uint(math.MaxUint64),
		},
		"Int8": {
			value: math.MaxInt8,
		},
		"Int16": {
			value: math.MaxInt16,
		},
		"Int32": {
			value: math.MaxInt32,
		},
		"Int64": {
			value: math.MaxInt64,
		},
		"Float32": {
			value: math.MaxFloat32,
		},
		"Float64": {
			value: math.MaxFloat64,
		},
		"String": {
			value: makeString(math.MaxUint8),
		},
		"Array": {
			value: []int{math.MaxUint8, math.MaxUint16, math.MaxUint32},
		},
		"Map": {
			value: map[string]uint{
				"Uint8":  math.MaxUint8,
				"Uint16": math.MaxUint16,
				"Uint32": math.MaxUint32,
				"Uint64": math.MaxUint64,
			},
		},
	}
	for name, bb := range benchs {
		bb := bb
		b.Run(name, func(b *testing.B) {
			enc := NewEncoder(ioutil.Discard)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := enc.Encode(bb.value); err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	for name, bb := range benchs {
		b.Run("Parallel/"+name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					enc := NewEncoder(ioutil.Discard)
					if err := enc.Encode(bb.value); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkPackBool(b *testing.B) {
	b.ReportAllocs()

	b.Run("False", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			v := false
			if err := enc.PackBool(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("True", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			v := true
			if err := enc.PackBool(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPackUint8(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := uint64(math.MaxUint8)
		if err := enc.PackUint(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackUint16(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := uint64(math.MaxUint16)
		if err := enc.PackUint(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackUint32(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := uint64(math.MaxUint32)
		if err := enc.PackUint(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackUint64(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := uint64(math.MaxUint64)
		if err := enc.PackUint(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackInt8(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := int64(math.MaxInt8)
		if err := enc.PackInt(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackInt16(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := int64(math.MaxInt16)
		if err := enc.PackInt(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackInt32(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := int64(math.MaxInt32)
		if err := enc.PackInt(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackInt64(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := int64(math.MaxInt64)
		if err := enc.PackInt(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackFloat32(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := float64(math.MaxFloat32)
		if err := enc.PackFloat(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackFloat64(b *testing.B) {
	b.ReportAllocs()

	enc := NewEncoder(ioutil.Discard)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := float64(math.MaxFloat64)
		if err := enc.PackFloat(v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPackString(b *testing.B) {
	b.ReportAllocs()

	b.Run("MaxUint8", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		s := makeString(math.MaxUint8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackString(s); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint8+1", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		s := makeString(math.MaxUint8 + 1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackString(s); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint16", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		s := makeString(math.MaxUint16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackString(s); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPackStringBytes(b *testing.B) {
	b.ReportAllocs()

	b.Run("MaxUint8", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		p := []byte(makeString(math.MaxUint8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackStringBytes(p); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint8+1", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		p := []byte(makeString(math.MaxUint8 + 1))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackStringBytes(p); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint16", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		p := []byte(makeString(math.MaxUint16))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackStringBytes(p); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPackBinary(b *testing.B) {
	b.ReportAllocs()

	b.Run("MaxUint8", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		p := []byte(makeString(math.MaxUint8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackBinary(p); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint8+1", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		p := []byte(makeString(math.MaxUint8 + 1))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackBinary(p); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint16", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		p := []byte(makeString(math.MaxUint16))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackBinary(p); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPackArrayLen(b *testing.B) {
	b.ReportAllocs()

	b.Run("MaxUint8", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		v := int64(math.MaxUint8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackArrayLen(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint16", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		v := int64(math.MaxUint16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackArrayLen(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint32", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		v := int64(math.MaxUint32)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackArrayLen(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPackMapLen(b *testing.B) {
	b.ReportAllocs()

	b.Run("MaxUint8", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		v := int64(math.MaxUint8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackMapLen(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint16", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		v := int64(math.MaxUint16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackMapLen(v); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MaxUint32", func(b *testing.B) {
		enc := NewEncoder(ioutil.Discard)
		v := int64(math.MaxUint32)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.PackMapLen(v); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()

	b.Run("Bool/False", func(b *testing.B) {
		p := []byte{falseCode, byte(0)}
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v bool
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != false {
				b.Fatal("not MaxUint8")
			}
		}
	})

	b.Run("Bool/True", func(b *testing.B) {
		p := []byte{trueCode, byte(1)}
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v bool
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != true {
				b.Fatal("not MaxUint8")
			}
		}
	})

	b.Run("Uint8", func(b *testing.B) {
		p := []byte{uint8Code, byte(math.MaxUint8)}
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v uint8
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxUint8 {
				b.Fatal("not MaxUint8")
			}
		}
	})

	b.Run("Uint16", func(b *testing.B) {
		p := make([]byte, 3)
		p[0] = uint16Code
		binary.BigEndian.PutUint16(p[1:], math.MaxUint16)
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v uint16
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxUint16 {
				b.Fatal("not MaxUint16")
			}
		}
	})

	b.Run("Uint32", func(b *testing.B) {
		p := make([]byte, 5)
		p[0] = uint32Code
		binary.BigEndian.PutUint32(p[1:], math.MaxUint32)
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v uint32
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxUint32 {
				b.Fatal("not MaxUint32")
			}
		}
	})

	b.Run("Uint64", func(b *testing.B) {
		p := make([]byte, 9)
		p[0] = uint64Code
		binary.BigEndian.PutUint64(p[1:], math.MaxUint64)
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v uint64
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxUint64 {
				b.Fatal("not MaxUint64")
			}
		}
	})

	b.Run("FixInt8", func(b *testing.B) {
		p := []byte{fixIntCodeMax, byte(math.MaxInt8)}
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v int8
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxInt8 {
				b.Fatal("not MaxInt8")
			}
		}
	})

	b.Run("Int8", func(b *testing.B) {
		p := []byte{int8Code, byte(math.MaxInt8)}
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v int8
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxInt8 {
				b.Fatal("not MaxInt8")
			}
		}
	})

	b.Run("Int16", func(b *testing.B) {
		p := make([]byte, 3)
		p[0] = int16Code
		binary.BigEndian.PutUint16(p[1:], math.MaxInt16)
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v int16
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxInt16 {
				b.Fatal("not MaxInt16")
			}
		}
	})

	b.Run("Int32", func(b *testing.B) {
		p := make([]byte, 5)
		p[0] = int32Code
		binary.BigEndian.PutUint32(p[1:], math.MaxInt32)
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v int32
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxInt32 {
				b.Fatal("not MaxInt32")
			}
		}
	})

	b.Run("Int64", func(b *testing.B) {
		p := make([]byte, 9)
		p[0] = int64Code
		binary.BigEndian.PutUint64(p[1:], math.MaxInt64)
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v int64
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxInt64 {
				b.Fatal("not MaxInt64")
			}
		}
	})

	b.Run("Float32", func(b *testing.B) {
		p := make([]byte, 5)
		p[0] = float32Code
		binary.BigEndian.PutUint32(p[1:], math.Float32bits(math.MaxFloat32))
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v float32
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxFloat32 {
				b.Fatal("not MaxFloat32")
			}
		}
	})

	b.Run("Float64", func(b *testing.B) {
		p := make([]byte, 9)
		p[0] = float64Code
		binary.BigEndian.PutUint64(p[1:], math.Float64bits(math.MaxFloat64))
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v float64
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != math.MaxFloat64 {
				b.Fatal("not MaxFloat64")
			}
		}
	})

	b.Run("String", func(b *testing.B) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		s := makeString(math.MaxUint16)
		if err := enc.PackString(s); err != nil {
			b.Fatal(err)
		}
		dec := NewDecoder(NewTestReader(buf.Bytes()))

		b.SetBytes(int64(len(s)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v string
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if v != s {
				b.Fatalf("not equal: got %s want %s", v, s)
			}
		}
	})

	b.Run("StringBytes", func(b *testing.B) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		var p []byte
		for _, byt := range [][]byte{[]byte(makeString(math.MaxUint8)), []byte(makeString(math.MaxUint8 + 1)), []byte(makeString(math.MaxUint16))} {
			p = append(p, byt...)
		}
		if err := enc.PackBinary(p); err != nil {
			b.Fatal(err)
		}
		dec := NewDecoder(NewTestReader(buf.Bytes()))

		b.SetBytes(int64(int64(len(p))))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v []byte
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(v, p) {
				b.Fatalf("not equal:\n got %s\nwant %s", string(v), string(p))
			}
		}
	})

	b.Run("Binary", func(b *testing.B) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		var p []byte
		for _, by := range [][]byte{[]byte(makeString(math.MaxUint8)), []byte(makeString(math.MaxUint8 + 1)), []byte(makeString(math.MaxUint16))} {
			p = append(p, by...)
		}
		if err := enc.PackBinary(p); err != nil {
			b.Fatal(err)
		}
		dec := NewDecoder(NewTestReader(buf.Bytes()))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v []byte
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(v, p) {
				b.Fatalf("not equal: got %s want %s", string(v), string(p))
			}
		}
	})

	b.Run("Array/String", func(b *testing.B) {
		ss := []string{makeString(math.MaxUint8), makeString(math.MaxUint8 + 1), makeString(math.MaxUint16)}
		builder := NewTestArrayBuilder(b)
		for _, s := range ss {
			builder.Add(s)
		}
		dec := NewDecoder(NewTestReader(builder.Bytes()))

		b.SetBytes(builder.Count())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v []string
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if !reflect.DeepEqual(v, ss) {
				b.Fatalf("not equal: got %#v want %#v", v, ss)
			}
		}
	})

	b.Run("Array/Uint64", func(b *testing.B) {
		u64s := []uint64{math.MaxUint8, math.MaxUint16, math.MaxUint32}
		builder := NewTestArrayBuilder(b)
		for _, u64 := range u64s {
			builder.Add(u64)
		}
		dec := NewDecoder(NewTestReader(builder.Bytes()))

		b.SetBytes(builder.Count())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v []uint64
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if !reflect.DeepEqual(v, u64s) {
				b.Fatalf("not equal: got %#v want %#v", v, u64s)
			}
		}
	})

	b.Run("Map/String", func(b *testing.B) {
		m := map[string]string{
			"foo": makeString(math.MaxUint8),
			"bar": makeString(math.MaxUint8 + 1),
			"baz": makeString(math.MaxUint8 + 2),
		}
		builder := NewTestMapBuilder(b)
		for k, v := range m {
			builder.Add(k, v)
		}
		dec := NewDecoder(NewTestReader(builder.Bytes()))

		b.SetBytes(builder.Count())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			v := make(map[string]string)
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if !reflect.DeepEqual(v, m) {
				b.Fatalf("not equal: got %#v want %#v", v, m)
			}
		}
	})

	b.Run("Map/Interface", func(b *testing.B) {
		m := map[string]interface{}{
			"uint8":  uint64(math.MaxUint8),
			"uint16": uint64(math.MaxUint16),
			"uint32": uint64(math.MaxUint32),
			"uint64": uint64(math.MaxUint64),
		}
		builder := NewTestMapBuilder(b)
		for k, v := range m {
			builder.Add(k, v)
		}
		dec := NewDecoder(NewTestReader(builder.Bytes()))

		b.SetBytes(builder.Count())
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			v := make(map[string]interface{})
			if err := dec.Decode(&v); err != nil {
				b.Fatal(err)
			}
			if !reflect.DeepEqual(v, m) {
				b.Fatalf("not equal: got %#v want %#v", v, m)
			}
		}
	})
}

func benchmarkDecoderBool(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new(bool)).Elem()
		for i := 0; i < b.N; i++ {
			boolDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderBool(b *testing.B) {
	b.ReportAllocs()

	b.Run("Bool/True", benchmarkDecoderBool(&decodeState{Decoder: &Decoder{n: uint64(1), t: Bool}}))
	b.Run("Bool/False", benchmarkDecoderBool(&decodeState{Decoder: &Decoder{n: uint64(0), t: Bool}}))

	b.Run("Int/True", benchmarkDecoderBool(&decodeState{Decoder: &Decoder{n: uint64(1), t: Int}}))
	b.Run("Int/False", benchmarkDecoderBool(&decodeState{Decoder: &Decoder{n: uint64(0), t: Int}}))

	b.Run("Uint/True", benchmarkDecoderBool(&decodeState{Decoder: &Decoder{n: uint64(1), t: Uint}}))
	b.Run("Uint/False", benchmarkDecoderBool(&decodeState{Decoder: &Decoder{n: uint64(0), t: Uint}}))
}

func benchmarkDecoderInt(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new(int64)).Elem()
		for i := 0; i < b.N; i++ {
			intDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderInt(b *testing.B) {
	b.ReportAllocs()

	b.Run("Int", benchmarkDecoderInt(&decodeState{Decoder: &Decoder{n: uint64(math.MaxInt64), t: Int}}))
	b.Run("Uint", benchmarkDecoderInt(&decodeState{Decoder: &Decoder{n: uint64(math.MaxUint32), t: Uint}}))
	b.Run("Float", benchmarkDecoderInt(&decodeState{Decoder: &Decoder{n: math.Float64bits(float64(math.MaxFloat64)), t: Float}}))
}

func benchmarkDecoderUint(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new(uint64)).Elem()
		for i := 0; i < b.N; i++ {
			uintDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderUint(b *testing.B) {
	b.ReportAllocs()

	b.Run("Int", benchmarkDecoderUint(&decodeState{Decoder: &Decoder{n: uint64(math.MaxInt64), t: Int}}))
	b.Run("Uint", benchmarkDecoderUint(&decodeState{Decoder: &Decoder{n: uint64(math.MaxUint64), t: Uint}}))
	b.Run("Float", benchmarkDecoderUint(&decodeState{Decoder: &Decoder{n: math.Float64bits(float64(math.MaxFloat64)), t: Float}}))
}

func benchmarkDecoderFloat(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new(float64)).Elem()
		for i := 0; i < b.N; i++ {
			floatDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderFloat(b *testing.B) {
	b.ReportAllocs()

	b.Run("Int", benchmarkDecoderFloat(&decodeState{Decoder: &Decoder{n: math.Float64bits(float64(math.MaxInt64)), t: Int}}))
	b.Run("Uint", benchmarkDecoderFloat(&decodeState{Decoder: &Decoder{n: math.Float64bits(float64(math.MaxUint64)), t: Uint}}))
	b.Run("Float", benchmarkDecoderFloat(&decodeState{Decoder: &Decoder{n: math.Float64bits(float64(math.MaxFloat64)), t: Float}}))
}

func benchmarkDecoderString(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new(string)).Elem()
		for i := 0; i < b.N; i++ {
			stringDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderString(b *testing.B) {
	b.ReportAllocs()

	b.Run("Binary", benchmarkDecoderString(&decodeState{Decoder: &Decoder{p: []byte(makeString(math.MaxInt16)), t: Binary}}))
	b.Run("String", benchmarkDecoderString(&decodeState{Decoder: &Decoder{p: []byte(makeString(math.MaxInt16)), t: String}}))
}

func benchmarkDecoderByteSlice(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new([]byte)).Elem()
		for i := 0; i < b.N; i++ {
			byteSliceDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderByteSlice(b *testing.B) {
	b.ReportAllocs()

	b.Run("Binary", benchmarkDecoderByteSlice(&decodeState{Decoder: &Decoder{p: []byte(makeString(math.MaxInt16)), t: Binary}}))
	b.Run("String", benchmarkDecoderByteSlice(&decodeState{Decoder: &Decoder{p: []byte(makeString(math.MaxInt16)), t: String}}))
}

func benchmarkDecoderInterface(ds *decodeState) func(*testing.B) {
	return func(b *testing.B) {
		v := reflect.ValueOf(new(interface{})).Elem()
		for i := 0; i < b.N; i++ {
			interfaceDecoder(ds, v)
		}
	}
}

func BenchmarkDecoderInterface(b *testing.B) {
	b.ReportAllocs()

	b.Run("Bool/True", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{n: 1, t: Bool}}))
	b.Run("Bool/False", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{n: 0, t: Bool}}))
	b.Run("Int", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{n: uint64(math.MaxInt64), t: Int}}))
	b.Run("Uint", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{n: uint64(math.MaxUint64), t: Uint}}))
	b.Run("Float", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{n: math.Float64bits(float64(math.MaxFloat64)), t: Float}}))
	b.Run("String", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{p: []byte(makeString(math.MaxInt16)), t: String}}))

	ab := NewTestArrayBuilder(b)
	ab.Add(math.MaxUint8)
	ab.Add(math.MaxUint16)
	ab.Add(math.MaxUint32)
	b.Run("ArrayLen", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{p: ab.Bytes(), t: ArrayLen}}))

	msb := NewTestMapBuilder(b)
	msb.Add("foo", makeString(math.MaxUint8))
	msb.Add("bar", makeString(math.MaxUint8+1))
	msb.Add("baz", makeString(math.MaxUint16))
	b.Run("MapLen/String", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{p: msb.Bytes(), t: MapLen}}))

	mub := NewTestMapBuilder(b)
	mub.Add("uint8", math.MaxUint8)
	mub.Add("uint16", math.MaxUint16)
	mub.Add("uint32", math.MaxUint32)
	b.Run("MapLen/Interface", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{p: mub.Bytes(), t: MapLen}}))

	b.Run("Nil", benchmarkDecoderInterface(&decodeState{Decoder: &Decoder{p: nil, t: Nil}}))
}

func BenchmarkUnpackUint8(b *testing.B) {
	b.ReportAllocs()

	p := []byte{uint8Code, byte(math.MaxUint8)}
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Uint(); v != math.MaxUint8 {
			b.Fatalf("expected math.MaxUint8: %d", v)
		}
	}
}

func BenchmarkUnpackUint16(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 3)
	p[0] = uint16Code
	binary.BigEndian.PutUint16(p[1:], math.MaxUint16)
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Uint(); v != math.MaxUint16 {
			b.Fatalf("expected math.MaxUint16: %d", v)
		}
	}
}

func BenchmarkUnpackUint32(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 5)
	p[0] = uint32Code
	binary.BigEndian.PutUint32(p[1:], math.MaxUint32)
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Uint(); v != math.MaxUint32 {
			b.Fatalf("expected math.MaxUint32: %d", v)
		}
	}
}

func BenchmarkUnpackUint64(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 9)
	p[0] = uint64Code
	binary.BigEndian.PutUint64(p[1:], math.MaxUint64)
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Uint(); v != math.MaxUint64 {
			b.Fatalf("expected math.MaxUint64: %d", v)
		}
	}
}

func BenchmarkUnpackFixInt8(b *testing.B) {
	b.ReportAllocs()

	p := []byte{fixIntCodeMax, byte(math.MaxInt8)}
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Int(); v != math.MaxInt8 {
			b.Fatalf("expected math.MaxInt8: %d", v)
		}
	}
}

func BenchmarkUnpackInt8(b *testing.B) {
	b.ReportAllocs()

	p := []byte{int8Code, byte(math.MaxInt8)}
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Int(); v != math.MaxInt8 {
			b.Fatalf("expected math.MaxInt8: %d", v)
		}
	}
}

func BenchmarkUnpackInt16(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 3)
	p[0] = int16Code
	binary.BigEndian.PutUint16(p[1:], math.MaxInt16)
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Int(); v != math.MaxInt16 {
			b.Fatalf("expected math.MaxInt16: %d", v)
		}
	}
}

func BenchmarkUnpackInt32(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 5)
	p[0] = int32Code
	binary.BigEndian.PutUint32(p[1:], math.MaxInt32)
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Int(); v != math.MaxInt32 {
			b.Fatalf("expected math.MaxInt32: %d", v)
		}
	}
}

func BenchmarkUnpackInt64(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 9)
	p[0] = int64Code
	binary.BigEndian.PutUint64(p[1:], math.MaxInt64)
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Int(); v != math.MaxInt64 {
			b.Fatalf("expected math.MaxInt64: %d", v)
		}
	}
}

func BenchmarkUnpackFloat32(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 5)
	p[0] = float32Code
	binary.BigEndian.PutUint32(p[1:], math.Float32bits(math.MaxFloat32))
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Float(); v != math.MaxFloat32 {
			b.Fatalf("expected math.MaxFloat32: %f", v)
		}
	}
}

func BenchmarkUnpackFloat64(b *testing.B) {
	b.ReportAllocs()

	p := make([]byte, 9)
	p[0] = float64Code
	binary.BigEndian.PutUint64(p[1:], math.Float64bits(math.MaxFloat64))
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.Float(); v != math.MaxFloat64 {
			b.Fatalf("expected math.MaxInt64: %f", v)
		}
	}
}

func BenchmarkUnpackString(b *testing.B) {
	b.ReportAllocs()

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	s := makeString(math.MaxUint16)
	if err := enc.PackString(s); err != nil {
		b.Fatal(err)
	}
	dec := NewDecoder(NewTestReader(buf.Bytes()))

	b.SetBytes(int64(len(s)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		if v := dec.String(); v != s {
			b.Fatalf("got %s but want: %s", v, s)
		}
	}
}

func BenchmarkUnpackArray(b *testing.B) {
	b.ReportAllocs()

	b.Run("Uint", func(b *testing.B) {
		a := NewTestArrayBuilder(b)
		a.Add(math.MaxUint8)
		a.Add(math.MaxUint16)
		a.Add(math.MaxUint32)
		p := a.Bytes()
		dec := NewDecoder(NewTestReader(p))

		b.SetBytes(int64(len(p)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := dec.Unpack(); err != nil {
				b.Fatal(err)
			}
			_ = uint32(dec.Int())
		}
	})

	b.Run("Binary", func(b *testing.B) {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		pp := [][]byte{[]byte(makeString(math.MaxUint8)), []byte(makeString(math.MaxUint8 + 1)), []byte(makeString(math.MaxUint16))}
		for _, p := range pp {
			if err := enc.PackBinary(p); err != nil {
				b.Fatal(err)
			}
		}
		dec := NewDecoder(NewTestReader(buf.Bytes()))

		b.SetBytes(int64(buf.Len()))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := dec.Unpack(); err != nil {
				b.Fatal(err)
			}
			_ = uint32(dec.Int())
		}
	})
}

func BenchmarkUnpackMap(b *testing.B) {
	b.ReportAllocs()

	m := NewTestMapBuilder(b)
	m.Add("uint8", math.MaxUint8)
	m.Add("uint16", math.MaxUint16)
	m.Add("uint32", math.MaxUint32)
	p := m.Bytes()
	dec := NewDecoder(NewTestReader(p))

	b.SetBytes(int64(len(p)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Unpack(); err != nil {
			b.Fatal(err)
		}
		_ = uint32(dec.Int())
	}
}

type namedStruct struct {
	Uint     uint      `msgpack:"uint"`
	Uint8    uint8     `msgpack:"uint8"`
	Uint16   uint16    `msgpack:"uint16"`
	Uint32   uint32    `msgpack:"uint32"`
	Uint64   uint64    `msgpack:"uint64"`
	Int      int       `msgpack:"int"`
	Int8     int8      `msgpack:"int8"`
	Int16    int16     `msgpack:"int16"`
	Int32    int32     `msgpack:"int32"`
	Int64    int64     `msgpack:"int64"`
	String   string    `msgpack:"string"`
	Bool     bool      `msgpack:"bool"`
	SString  []string  `msgpack:"sstring"`
	SInt     []int     `msgpack:"sint"`
	SInt8    []int8    `msgpack:"sint8"`
	SInt16   []int16   `msgpack:"sint16"`
	SInt32   []int32   `msgpack:"sint32"`
	SInt64   []int64   `msgpack:"sint64"`
	SFloat32 []float32 `msgpack:"sfloat32"`
	SFloat64 []float64 `msgpack:"sfloat64"`
	SBool    []bool    `msgpack:"sbool"`
	Struct   struct {
		I64           int64
		U64           uint64
		AnotherStruct struct {
			S string
		} `msgpack:",array"`
	} `msgpack:"struct"`
}

type emptyStruct struct {
	Uint     uint      `msgpack:"uint" empty:"1"`
	Uint8    uint8     `msgpack:"uint8" empty:"1"`
	Uint16   uint16    `msgpack:"uint16" empty:"1"`
	Uint32   uint32    `msgpack:"uint32" empty:"1"`
	Uint64   uint64    `msgpack:"uint64" empty:"1"`
	Int      int       `msgpack:"int" empty:"1"`
	Int8     int8      `msgpack:"int8" empty:"1"`
	Int16    int16     `msgpack:"int16" empty:"1"`
	Int32    int32     `msgpack:"int32" empty:"1"`
	Int64    int64     `msgpack:"int64" empty:"1"`
	String   string    `msgpack:"string" empty:"foo"`
	Bool     bool      `msgpack:"bool" empty:"true"`
	SString  []string  `msgpack:"sstring"`
	SInt     []int     `msgpack:"sint"`
	SInt8    []int8    `msgpack:"sint8"`
	SInt16   []int16   `msgpack:"sint16"`
	SInt32   []int32   `msgpack:"sint32"`
	SInt64   []int64   `msgpack:"sint64"`
	SFloat32 []float32 `msgpack:"sfloat32"`
	SFloat64 []float64 `msgpack:"sfloat64"`
	SBool    []bool    `msgpack:"sbool"`
	Struct   struct {
		I64           int64
		U64           uint64
		AnotherStruct struct {
			S string
		} `msgpack:",array,omitempty"`
	} `msgpack:"struct,omitempty"`
}

type omitemptyStruct struct {
	Uint     uint      `msgpack:"uint,omitempty"`
	Uint8    uint8     `msgpack:"uint8,omitempty"`
	Uint16   uint16    `msgpack:"uint16,omitempty"`
	Uint32   uint32    `msgpack:"uint32,omitempty"`
	Uint64   uint64    `msgpack:"uint64,omitempty"`
	Int      int       `msgpack:"int,omitempty"`
	Int8     int8      `msgpack:"int8,omitempty"`
	Int16    int16     `msgpack:"int16,omitempty"`
	Int32    int32     `msgpack:"int32,omitempty"`
	Int64    int64     `msgpack:"int64,omitempty"`
	String   string    `msgpack:"string,omitempty"`
	Bool     bool      `msgpack:"bool,omitempty"`
	SString  []string  `msgpack:"sstring,omitempty"`
	SInt     []int     `msgpack:"sint,omitempty"`
	SInt8    []int8    `msgpack:"sint8,omitempty"`
	SInt16   []int16   `msgpack:"sint16,omitempty"`
	SInt32   []int32   `msgpack:"sint32,omitempty"`
	SInt64   []int64   `msgpack:"sint64,omitempty"`
	SFloat32 []float32 `msgpack:"sfloat32,omitempty"`
	SFloat64 []float64 `msgpack:"sfloat64,omitempty"`
	SBool    []bool    `msgpack:"sbool,omitempty"`
	Struct   struct {
		I64           int64
		U64           uint64
		AnotherStruct struct {
			S string
		} `msgpack:",array,omitempty"`
	} `msgpack:"struct,omitempty"`
}

func Benchmark_correctFilels(b *testing.B) {
	b.ReportAllocs()

	b.Run("Named", func(b *testing.B) {
		t := reflect.ValueOf(namedStruct{}).Type()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = collectFields(nil, t, make(map[reflect.Type]bool), make(map[string]int), nil)
		}
	})

	b.Run("Empty", func(b *testing.B) {
		t := reflect.ValueOf(namedStruct{}).Type()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = collectFields(nil, t, make(map[reflect.Type]bool), make(map[string]int), nil)
		}
	})

	b.Run("OmitEmpty", func(b *testing.B) {
		t := reflect.ValueOf(omitemptyStruct{}).Type()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = collectFields(nil, t, make(map[reflect.Type]bool), make(map[string]int), nil)
		}
	})
}

type api struct {
	Annotations   []string          `msgpack:"annotations,omitempty"`
	Doc           []string          `msgpack:"doc,omitempty"`
	Parameters    [][2]string       `msgpack:"parameters"`
	ParametersDoc map[string]string `msgpack:"parameters_doc,omitempty"`
	Return        []string          `msgpack:"return,omitempty"`
	SeeAlso       []string          `msgpack:"seealso,omitempty"`
	Signature     string            `msgpack:"signature,omitempty"`
}

type apiMetadata struct {
	CanFail           bool        `msgpack:"can_fail,omitempty"`
	DeprecatedSince   int         `msgpack:"deprecated_since,omitempty"`
	Fast              bool        `msgpack:"fast"`
	ImplName          string      `msgpack:"impl_name,omitempty"`
	Method            bool        `msgpack:"method"`
	Name              string      `msgpack:"name"`
	Parameters        [][2]string `msgpack:"parameters"`
	ReceivesChannelID bool        `msgpack:"receives_channel_id,omitempty"`
	RemoteOnly        bool        `msgpack:"remote_only,omitempty"`
	ReturnType        string      `msgpack:"return_type"`
	Since             int         `msgpack:"since,omitempty"`
}

type apiMetadataEmpty struct {
	CanFail           bool        `msgpack:"can_fail,omitempty" empty:"true"`
	DeprecatedSince   int         `msgpack:"deprecated_since,omitempty" empty:"1"`
	Fast              bool        `msgpack:"fast"`
	ImplName          string      `msgpack:"impl_name,omitempty" empty:"unkonwn"`
	Method            bool        `msgpack:"method"`
	Name              string      `msgpack:"name"`
	Parameters        [][2]string `msgpack:"parameters"`
	ReceivesChannelID bool        `msgpack:"receives_channel_id,omitempty" empty:"true"`
	RemoteOnly        bool        `msgpack:"remote_only,omitempty" empty:"true"`
	ReturnType        string      `msgpack:"return_type"`
	Since             int         `msgpack:"since,omitempty" empty:"1"`
}

type funcsData struct {
	Args []int `msgpack:"args,omitempty"`
}

func extractMpack(tb testing.TB, path string) []byte {
	tb.Helper()

	f, err := os.Open(path)
	if err != nil {
		tb.Fatal(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		tb.Fatal(err)
	}

	data, err := ioutil.ReadAll(gz)
	if err != nil {
		tb.Fatal(err)
	}

	return data
}

func BenchmarkEncodeMpack(b *testing.B) {
	b.ReportAllocs()

	b.Run("api", func(b *testing.B) {
		var structAPI api

		data := extractMpack(b, filepath.Join("testdata", "api.mpack.gz"))
		dec := NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&structAPI); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			enc := NewEncoder(ioutil.Discard)
			for pb.Next() {
				if err := enc.Encode(&structAPI); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})

	b.Run("api_metadata", func(b *testing.B) {
		var structAPIMetadata []apiMetadata

		data := extractMpack(b, filepath.Join("testdata", "api_metadata.mpack.gz"))
		dec := NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&structAPIMetadata); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			enc := NewEncoder(ioutil.Discard)
			for pb.Next() {
				if err := enc.Encode(&structAPIMetadata); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})

	b.Run("api_metadata/Empty", func(b *testing.B) {
		var structAPIMetadata []apiMetadataEmpty

		data := extractMpack(b, filepath.Join("testdata", "api_metadata.mpack.gz"))
		dec := NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&structAPIMetadata); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			enc := NewEncoder(ioutil.Discard)
			for pb.Next() {
				if err := enc.Encode(&structAPIMetadata); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})

	b.Run("funcs_data", func(b *testing.B) {
		var structFuncsData funcsData

		data := extractMpack(b, filepath.Join("testdata", "funcs_data.mpack.gz"))
		dec := NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&structFuncsData); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			enc := NewEncoder(ioutil.Discard)
			for pb.Next() {
				if err := enc.Encode(&structFuncsData); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})
}

func BenchmarkDecodeMpack(b *testing.B) {
	b.ReportAllocs()

	b.Run("api", func(b *testing.B) {
		data := extractMpack(b, filepath.Join("testdata", "api.mpack.gz"))
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			var r api
			for pb.Next() {
				buf.Write(data)
				if err := dec.Decode(&r); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})

	b.Run("api_metadata", func(b *testing.B) {
		data := extractMpack(b, filepath.Join("testdata", "api_metadata.mpack.gz"))
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			var r []apiMetadata
			for pb.Next() {
				buf.Write(data)
				if err := dec.Decode(&r); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})

	b.Run("api_metadata/Empty", func(b *testing.B) {
		data := extractMpack(b, filepath.Join("testdata", "api_metadata.mpack.gz"))
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			var r []apiMetadataEmpty
			for pb.Next() {
				buf.Write(data)
				if err := dec.Decode(&r); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})

	b.Run("funcs_data", func(b *testing.B) {
		data := extractMpack(b, filepath.Join("testdata", "funcs_data.mpack.gz"))
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			var r funcsData
			for pb.Next() {
				buf.Write(data)
				if err := dec.Decode(&r); err != nil {
					b.Fatalf("Decode: %v", err)
				}
			}
		})

		b.SetBytes(int64(len(data)))
	})
}
