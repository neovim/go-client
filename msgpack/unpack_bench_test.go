package msgpack

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func BenchmarkUnpackInt(b *testing.B) {
	benchs := []struct {
		name string
		v    int64
	}{
		{name: "int64(0x0)", v: int64(0x0)},
		{name: "int64(0x1)", v: int64(0x1)},
		{name: "int64(0x7f)", v: int64(0x7f)},
		{name: "int64(0x80)", v: int64(0x80)},
		{name: "int64(0x7fff)", v: int64(0x7fff)},
		{name: "int64(0x8000)", v: int64(0x8000)},
		{name: "int64(0x7fffffff)", v: int64(0x7fffffff)},
		{name: "int64(0x80000000)", v: int64(0x80000000)},
		{name: "int64(0x7fffffffffffffff)", v: int64(0x7fffffffffffffff)},
		{name: "int64(-0x1)", v: int64(-0x1)},
		{name: "int64(-0x20)", v: int64(-0x20)},
		{name: "int64(-0x21)", v: int64(-0x21)},
		{name: "int64(-0x80)", v: int64(-0x80)},
		{name: "int64(-0x81)", v: int64(-0x81)},
		{name: "int64(-0x8000)", v: int64(-0x8000)},
		{name: "int64(-0x8001)", v: int64(-0x8001)},
		{name: "int64(-0x80000000)", v: int64(-0x80000000)},
		{name: "int64(-0x80000001)", v: int64(-0x80000001)},
		{name: "int64(-0x8000000000000000)", v: int64(-0x8000000000000000)},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.Int()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackUint(b *testing.B) {
	benchs := []struct {
		name string
		v    uint64
	}{
		{name: "uint64(0xff)", v: uint64(0xff)},
		{name: "uint64(0x100)", v: uint64(0x100)},
		{name: "uint64(0xffff)", v: uint64(0xffff)},
		{name: "uint64(0x10000)", v: uint64(0x10000)},
		{name: "uint64(0xffffffff)", v: uint64(0xffffffff)},
		{name: "uint64(0x100000000)", v: uint64(0x100000000)},
		{name: "uint64(0xffffffffffffffff)", v: uint64(0xffffffffffffffff)},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.Uint()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackBool(b *testing.B) {
	benchs := []struct {
		name string
		v    bool
	}{
		{name: "true", v: true},
		{name: "false", v: false},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.Bool()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackFloat(b *testing.B) {
	benchs := []struct {
		name string
		v    float64
	}{
		{name: "float64(123456)", v: float64(123456)},
		{name: "float64(1.23456)", v: float64(1.23456)},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.Float()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackArrayLen(b *testing.B) {
	benchs := []struct {
		name string
		v    arrayLen
	}{
		{name: "arrayLen(0x0)", v: arrayLen(0x0)},
		{name: "arrayLen(0x1)", v: arrayLen(0x1)},
		{name: "arrayLen(0xf)", v: arrayLen(0xf)},
		{name: "arrayLen(0x10)", v: arrayLen(0x10)},
		{name: "arrayLen(0xffff)", v: arrayLen(0xffff)},
		{name: "arrayLen(0x10000)", v: arrayLen(0x10000)},
		{name: "arrayLen(0xffffffff)", v: arrayLen(0xffffffff)},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = arrayLen(dec.Int())
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackMapLen(b *testing.B) {
	benchs := []struct {
		name string
		v    mapLen
	}{
		{name: "mapLen(0x0)", v: mapLen(0x0)},
		{name: "mapLen(0x1)", v: mapLen(0x1)},
		{name: "mapLen(0xf)", v: mapLen(0xf)},
		{name: "mapLen(0x10)", v: mapLen(0x10)},
		{name: "mapLen(0xffff)", v: mapLen(0xffff)},
		{name: "mapLen(0x10000)", v: mapLen(0x10000)},
		{name: "mapLen(0xffffffff)", v: mapLen(0xffffffff)},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = mapLen(dec.Int())
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackString(b *testing.B) {
	benchs := []struct {
		name string
		v    string
	}{
		{name: "string(1234567890123456789012345678901)", v: "1234567890123456789012345678901"},
		{name: "string(12345678901234567890123456789012)", v: "12345678901234567890123456789012"},
		{name: "emptyString", v: ""},
		{name: "string(1)", v: "1"},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.String()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackBinary(b *testing.B) {
	benchs := []struct {
		name string
		v    []byte
	}{
		{name: "[]byte(``)", v: []byte("")},
		{name: "[]byte(`1`)", v: []byte("1")},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.Bytes()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackExtension(b *testing.B) {
	benchs := []struct {
		name string
		v    extension
	}{
		{name: "extension{1, ``}", v: extension{1, ""}},
		{name: "extension{2, `1`}", v: extension{2, "1"}},
		{name: "extension{3, `12`}", v: extension{3, "12"}},
		{name: "extension{4, `1234`}", v: extension{4, "1234"}},
		{name: "extension{5, `12345678`}", v: extension{5, "12345678"}},
		{name: "extension{6, `1234567890123456`}", v: extension{6, "1234567890123456"}},
		{name: "extension{7, `12345678901234567`}", v: extension{7, "12345678901234567"}},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			dec := NewDecoder(&buf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = dec.Extension()
				_ = dec.String()
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnpackNil(b *testing.B) {
	benchs := []struct {
		name string
		v    interface{}
	}{
		{name: "nil", v: nil},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			var buf bytes.Buffer
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// do nothing
			}

			b.SetBytes(int64(buf.Len()))
		})
	}
}

func BenchmarkUnPack(b *testing.B) {
	benchs := []struct {
		name string
		typ  Type
		v    interface{}
		hs   []string
	}{
		{name: "int64(0x0)", typ: Int, v: int64(0x0), hs: []string{"00", "d000", "d10000", "d200000000", "d30000000000000000"}},
		{name: "int64(0x1)", typ: Int, v: int64(0x1), hs: []string{"01", "d001", "d10001", "d200000001", "d30000000000000001"}},
		{name: "int64(0x7f)", typ: Int, v: int64(0x7f), hs: []string{"7f", "d07f", "d1007f", "d20000007f", "d3000000000000007f"}},
		{name: "int64(0x80)", typ: Int, v: int64(0x80), hs: []string{"d10080", "d200000080", "d30000000000000080"}},
		{name: "int64(0x7fff)", typ: Int, v: int64(0x7fff), hs: []string{"d17fff", "d200007fff", "d30000000000007fff"}},
		{name: "int64(0x8000)", typ: Int, v: int64(0x8000), hs: []string{"d200008000", "d30000000000008000"}},
		{name: "int64(0x7fffffff)", typ: Int, v: int64(0x7fffffff), hs: []string{"d27fffffff", "d3000000007fffffff"}},
		{name: "int64(0x80000000)", typ: Int, v: int64(0x80000000), hs: []string{"d30000000080000000"}},
		{name: "int64(0x7fffffffffffffff)", typ: Int, v: int64(0x7fffffffffffffff), hs: []string{"d37fffffffffffffff"}},
		{name: "int64(-0x1)", typ: Int, v: int64(-0x1), hs: []string{"ff", "d0ff", "d1ffff", "d2ffffffff", "d3ffffffffffffffff"}},
		{name: "int64(-0x20)", typ: Int, v: int64(-0x20), hs: []string{"e0", "d0e0", "d1ffe0", "d2ffffffe0", "d3ffffffffffffffe0"}},
		{name: "int64(-0x21)", typ: Int, v: int64(-0x21), hs: []string{"d0df", "d1ffdf", "d2ffffffdf", "d3ffffffffffffffdf"}},
		{name: "int64(-0x80)", typ: Int, v: int64(-0x80), hs: []string{"d080", "d1ff80", "d2ffffff80", "d3ffffffffffffff80"}},
		{name: "int64(-0x81)", typ: Int, v: int64(-0x81), hs: []string{"d1ff7f", "d2ffffff7f", "d3ffffffffffffff7f"}},
		{name: "int64(-0x8000)", typ: Int, v: int64(-0x8000), hs: []string{"d18000", "d2ffff8000", "d3ffffffffffff8000"}},
		{name: "int64(-0x8001)", typ: Int, v: int64(-0x8001), hs: []string{"d2ffff7fff", "d3ffffffffffff7fff"}},
		{name: "int64(-0x80000000)", typ: Int, v: int64(-0x80000000), hs: []string{"d280000000", "d3ffffffff80000000"}},
		{name: "int64(-0x80000001)", typ: Int, v: int64(-0x80000001), hs: []string{"d3ffffffff7fffffff"}},
		{name: "int64(-0x8000000000000000)", typ: Int, v: int64(-0x8000000000000000), hs: []string{"d38000000000000000"}},
		{name: "uint64(0xff)", typ: Uint, v: uint64(0xff), hs: []string{"ccff", "cd00ff", "ce000000ff", "cf00000000000000ff"}},
		{name: "uint64(0x100)", typ: Uint, v: uint64(0x100), hs: []string{"cd0100", "ce00000100", "cf0000000000000100"}},
		{name: "uint64(0xffff)", typ: Uint, v: uint64(0xffff), hs: []string{"cdffff", "ce0000ffff", "cf000000000000ffff"}},
		{name: "uint64(0x10000)", typ: Uint, v: uint64(0x10000), hs: []string{"ce00010000", "cf0000000000010000"}},
		{name: "uint64(0xffffffff)", typ: Uint, v: uint64(0xffffffff), hs: []string{"ceffffffff", "cf00000000ffffffff"}},
		{name: "uint64(0x100000000)", typ: Uint, v: uint64(0x100000000), hs: []string{"cf0000000100000000"}},
		{name: "uint64(0xffffffffffffffff)", typ: Uint, v: uint64(0xffffffffffffffff), hs: []string{"cfffffffffffffffff"}},
		{name: "true", typ: Bool, v: true, hs: []string{"c3"}},
		{name: "false", typ: Bool, v: false, hs: []string{"c2"}},
		{name: "float64(123456)", typ: Float, v: float64(123456), hs: []string{"ca47f12000"}},
		{name: "float64(1.23456)", typ: Float, v: float64(1.23456), hs: []string{"cb3ff3c0c1fc8f3238"}},
		{name: "arrayLen(0x0)", typ: ArrayLen, v: arrayLen(0x0), hs: []string{"90", "dc0000", "dd00000000"}},
		{name: "arrayLen(0x1)", typ: ArrayLen, v: arrayLen(0x1), hs: []string{"91", "dc0001", "dd00000001"}},
		{name: "arrayLen(0xf)", typ: ArrayLen, v: arrayLen(0xf), hs: []string{"9f", "dc000f", "dd0000000f"}},
		{name: "arrayLen(0x10)", typ: ArrayLen, v: arrayLen(0x10), hs: []string{"dc0010", "dd00000010"}},
		{name: "arrayLen(0xffff)", typ: ArrayLen, v: arrayLen(0xffff), hs: []string{"dcffff", "dd0000ffff"}},
		{name: "arrayLen(0x10000)", typ: ArrayLen, v: arrayLen(0x10000), hs: []string{"dd00010000"}},
		{name: "arrayLen(0xffffffff)", typ: ArrayLen, v: arrayLen(0xffffffff), hs: []string{"ddffffffff"}},
		{name: "mapLen(0x0)", typ: MapLen, v: mapLen(0x0), hs: []string{"80", "de0000", "df00000000"}},
		{name: "mapLen(0x1)", typ: MapLen, v: mapLen(0x1), hs: []string{"81", "de0001", "df00000001"}},
		{name: "mapLen(0xf)", typ: MapLen, v: mapLen(0xf), hs: []string{"8f", "de000f", "df0000000f"}},
		{name: "mapLen(0x10)", typ: MapLen, v: mapLen(0x10), hs: []string{"de0010", "df00000010"}},
		{name: "mapLen(0xffff)", typ: MapLen, v: mapLen(0xffff), hs: []string{"deffff", "df0000ffff"}},
		{name: "mapLen(0x10000)", typ: MapLen, v: mapLen(0x10000), hs: []string{"df00010000"}},
		{name: "mapLen(0xffffffff)", typ: MapLen, v: mapLen(0xffffffff), hs: []string{"dfffffffff"}},
		{name: "emptyString", typ: String, v: "", hs: []string{"a0", "d900", "da0000", "db00000000"}},
		{name: "1", typ: String, v: "1", hs: []string{"a131", "d90131", "da000131", "db0000000131"}},
		{name: "string(1234567890123456789012345678901)", typ: String, v: "1234567890123456789012345678901", hs: []string{
			"bf31323334353637383930313233343536373839303132333435363738393031",
			"d91f31323334353637383930313233343536373839303132333435363738393031",
			"da001f31323334353637383930313233343536373839303132333435363738393031",
			"db0000001f31323334353637383930313233343536373839303132333435363738393031"}},
		{name: "string(12345678901234567890123456789012)", typ: String, v: "12345678901234567890123456789012", hs: []string{
			"d9203132333435363738393031323334353637383930313233343536373839303132",
			"da00203132333435363738393031323334353637383930313233343536373839303132",
			"db000000203132333435363738393031323334353637383930313233343536373839303132"}},
		{name: `Binary("")`, typ: Binary, v: "", hs: []string{"c400", "c50000", "c600000000"}},
		{name: `Binary("1")`, typ: Binary, v: "1", hs: []string{"c40131", "c5000131", "c60000000131"}},
		{name: `extension{1, ""}`, typ: Extension, v: extension{1, ""}, hs: []string{"c70001", "c8000001", "c90000000001"}},
		{name: `extension{2, "1"}`, typ: Extension, v: extension{2, "1"}, hs: []string{"d40231", "c7010231", "c800010231", "c9000000010231"}},
		{name: `extension{3, "12"}`, typ: Extension, v: extension{3, "12"}, hs: []string{"d5033132", "c702033132", "c80002033132", "c900000002033132"}},
		{name: `extension{4, "1234"}`, typ: Extension, v: extension{4, "1234"}, hs: []string{
			"d60431323334",
			"c7040431323334",
			"c800040431323334",
			"c9000000040431323334"}},
		{name: `extension{5, "12345678"}`, typ: Extension, v: extension{5, "12345678"}, hs: []string{
			"d7053132333435363738",
			"c708053132333435363738",
			"c80008053132333435363738",
			"c900000008053132333435363738"}},
		{name: `extension{6, "1234567890123456"}`, typ: Extension, v: extension{6, "1234567890123456"}, hs: []string{
			"d80631323334353637383930313233343536",
			"c7100631323334353637383930313233343536",
			"c800100631323334353637383930313233343536",
			"c9000000100631323334353637383930313233343536"}},
		{name: `extension{7, "12345678901234567"}`, typ: Extension, v: extension{7, "12345678901234567"}, hs: []string{
			"c711073132333435363738393031323334353637",
			"c80011073132333435363738393031323334353637",
			"c900000011073132333435363738393031323334353637"}},
		{name: "nil", typ: Nil, v: nil, hs: []string{"c0"}},
	}

	for _, bb := range benchs {
		b.Run(bb.name, func(b *testing.B) {
			for _, h := range bb.hs {
				p, err := hex.DecodeString(h)
				if err != nil {
					b.Errorf("decode(%s) returned error %v", h, err)
				}
				r := bytes.NewReader(p)
				dec := NewDecoder(r)
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					_ = dec.Unpack()

					switch v := bb.v.(type) {
					case int64, arrayLen, mapLen:
						_ = dec.Int()
					case uint64:
						_ = dec.Uint()
					case bool:
						_ = dec.Bool()
					case float64:
						_ = dec.Float()
					case string:
						_ = dec.String()
					case extension:
						_, _ = dec.Extension(), dec.String()
					case nil:
						// do nothing
					default:
						b.Errorf("no check for %T", v)
					}
				}
			}
		})
	}
}
