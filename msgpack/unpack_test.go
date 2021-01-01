package msgpack

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
)

var unpackTests = []struct {
	// Expected type
	typ Type
	// Expected value
	v interface{}
	// Hex encodings of typ, v
	hs []string
}{
	{
		typ: Int,
		v:   int64(0x0),
		hs:  []string{"00", "d000", "d10000", "d200000000", "d30000000000000000"},
	},
	{
		typ: Int,
		v:   int64(0x1),
		hs:  []string{"01", "d001", "d10001", "d200000001", "d30000000000000001"},
	},
	{
		typ: Int,
		v:   int64(0x7f),
		hs:  []string{"7f", "d07f", "d1007f", "d20000007f", "d3000000000000007f"},
	},
	{
		typ: Int,
		v:   int64(0x80),
		hs:  []string{"d10080", "d200000080", "d30000000000000080"},
	},
	{
		typ: Int,
		v:   int64(0x7fff),
		hs:  []string{"d17fff", "d200007fff", "d30000000000007fff"},
	},
	{
		typ: Int,
		v:   int64(0x8000),
		hs:  []string{"d200008000", "d30000000000008000"},
	},
	{
		typ: Int,
		v:   int64(0x7fffffff),
		hs:  []string{"d27fffffff", "d3000000007fffffff"},
	},
	{
		typ: Int,
		v:   int64(0x80000000),
		hs:  []string{"d30000000080000000"},
	},
	{
		typ: Int,
		v:   int64(0x7fffffffffffffff),
		hs:  []string{"d37fffffffffffffff"},
	},
	{
		typ: Int,
		v:   int64(-0x1),
		hs:  []string{"ff", "d0ff", "d1ffff", "d2ffffffff", "d3ffffffffffffffff"},
	},
	{
		typ: Int,
		v:   int64(-0x20),
		hs:  []string{"e0", "d0e0", "d1ffe0", "d2ffffffe0", "d3ffffffffffffffe0"},
	},
	{
		typ: Int,
		v:   int64(-0x21),
		hs:  []string{"d0df", "d1ffdf", "d2ffffffdf", "d3ffffffffffffffdf"},
	},
	{
		typ: Int,
		v:   int64(-0x80),
		hs:  []string{"d080", "d1ff80", "d2ffffff80", "d3ffffffffffffff80"},
	},
	{
		typ: Int,
		v:   int64(-0x81),
		hs:  []string{"d1ff7f", "d2ffffff7f", "d3ffffffffffffff7f"},
	},
	{
		typ: Int,
		v:   int64(-0x8000),
		hs:  []string{"d18000", "d2ffff8000", "d3ffffffffffff8000"},
	},
	{
		typ: Int,
		v:   int64(-0x8001),
		hs:  []string{"d2ffff7fff", "d3ffffffffffff7fff"},
	},
	{
		typ: Int,
		v:   int64(-0x80000000),
		hs:  []string{"d280000000", "d3ffffffff80000000"},
	},
	{
		typ: Int,
		v:   int64(-0x80000001),
		hs:  []string{"d3ffffffff7fffffff"},
	},
	{
		typ: Int,
		v:   int64(-0x8000000000000000),
		hs:  []string{"d38000000000000000"},
	},
	{
		typ: Uint,
		v:   uint64(0xff),
		hs:  []string{"ccff", "cd00ff", "ce000000ff", "cf00000000000000ff"},
	},
	{
		typ: Uint,
		v:   uint64(0x100),
		hs:  []string{"cd0100", "ce00000100", "cf0000000000000100"},
	},
	{
		typ: Uint,
		v:   uint64(0xffff),
		hs:  []string{"cdffff", "ce0000ffff", "cf000000000000ffff"},
	},
	{
		typ: Uint,
		v:   uint64(0x10000),
		hs:  []string{"ce00010000", "cf0000000000010000"},
	},
	{
		typ: Uint,
		v:   uint64(0xffffffff),
		hs:  []string{"ceffffffff", "cf00000000ffffffff"},
	},
	{
		typ: Uint,
		v:   uint64(0x100000000),
		hs:  []string{"cf0000000100000000"},
	},
	{
		typ: Uint,
		v:   uint64(0xffffffffffffffff),
		hs:  []string{"cfffffffffffffffff"},
	},
	{
		typ: Nil,
		v:   nil,
		hs:  []string{"c0"},
	},
	{
		typ: Bool,
		v:   true,
		hs:  []string{"c3"},
	},
	{
		typ: Bool,
		v:   false,
		hs:  []string{"c2"},
	},
	{
		typ: Float,
		v:   float64(123456),
		hs:  []string{"ca47f12000"},
	},
	{
		typ: Float,
		v:   float64(1.23456),
		hs:  []string{"cb3ff3c0c1fc8f3238"},
	},
	{
		typ: MapLen,
		v:   int64(0x0),
		hs:  []string{"80", "de0000", "df00000000"},
	},
	{
		typ: MapLen,
		v:   int64(0x1),
		hs:  []string{"81", "de0001", "df00000001"},
	},
	{
		typ: MapLen,
		v:   int64(0xf),
		hs:  []string{"8f", "de000f", "df0000000f"},
	},
	{
		typ: MapLen,
		v:   int64(0x10),
		hs:  []string{"de0010", "df00000010"},
	},
	{
		typ: MapLen,
		v:   int64(0xffff),
		hs:  []string{"deffff", "df0000ffff"},
	},
	{
		typ: MapLen,
		v:   int64(0x10000),
		hs:  []string{"df00010000"},
	},
	{
		typ: MapLen,
		v:   int64(0xffffffff),
		hs:  []string{"dfffffffff"},
	},
	{
		typ: ArrayLen,
		v:   int64(0x0),
		hs:  []string{"90", "dc0000", "dd00000000"},
	},
	{
		typ: ArrayLen,
		v:   int64(0x1),
		hs:  []string{"91", "dc0001", "dd00000001"},
	},
	{
		typ: ArrayLen,
		v:   int64(0xf),
		hs:  []string{"9f", "dc000f", "dd0000000f"},
	},
	{
		typ: ArrayLen,
		v:   int64(0x10),
		hs:  []string{"dc0010", "dd00000010"},
	},
	{
		typ: ArrayLen,
		v:   int64(0xffff),
		hs:  []string{"dcffff", "dd0000ffff"},
	},
	{
		typ: ArrayLen,
		v:   int64(0x10000),
		hs:  []string{"dd00010000"},
	},
	{
		typ: ArrayLen,
		v:   int64(0xffffffff),
		hs:  []string{"ddffffffff"},
	},
	{
		typ: String,
		v:   "",
		hs:  []string{"a0", "d900", "da0000", "db00000000"},
	},
	{
		typ: String,
		v:   "1",
		hs:  []string{"a131", "d90131", "da000131", "db0000000131"},
	},
	{
		typ: String,
		v:   "1234567890123456789012345678901",
		hs: []string{
			"bf31323334353637383930313233343536373839303132333435363738393031",
			"d91f31323334353637383930313233343536373839303132333435363738393031",
			"da001f31323334353637383930313233343536373839303132333435363738393031",
			"db0000001f31323334353637383930313233343536373839303132333435363738393031",
		},
	},
	{
		typ: String,
		v:   "12345678901234567890123456789012",
		hs: []string{
			"d9203132333435363738393031323334353637383930313233343536373839303132",
			"da00203132333435363738393031323334353637383930313233343536373839303132",
			"db000000203132333435363738393031323334353637383930313233343536373839303132",
		},
	},
	{
		typ: Binary,
		v:   "",
		hs:  []string{"c400", "c50000", "c600000000"},
	},
	{
		typ: Binary,
		v:   "1",
		hs:  []string{"c40131", "c5000131", "c60000000131"},
	},
	{
		typ: Extension,
		v:   extension{1, ""},
		hs:  []string{"c70001", "c8000001", "c90000000001"},
	},
	{
		typ: Extension,
		v:   extension{2, "1"},
		hs:  []string{"d40231", "c7010231", "c800010231", "c9000000010231"},
	},
	{
		typ: Extension,
		v:   extension{3, "12"},
		hs:  []string{"d5033132", "c702033132", "c80002033132", "c900000002033132"},
	},
	{
		typ: Extension,
		v:   extension{4, "1234"},
		hs: []string{
			"d60431323334",
			"c7040431323334",
			"c800040431323334",
			"c9000000040431323334",
		},
	},
	{
		typ: Extension,
		v:   extension{5, "12345678"},
		hs: []string{
			"d7053132333435363738",
			"c708053132333435363738",
			"c80008053132333435363738",
			"c900000008053132333435363738",
		},
	},
	{
		typ: Extension,
		v:   extension{6, "1234567890123456"},
		hs: []string{
			"d80631323334353637383930313233343536",
			"c7100631323334353637383930313233343536",
			"c800100631323334353637383930313233343536",
			"c9000000100631323334353637383930313233343536",
		},
	},
	{
		typ: Extension,
		v:   extension{7, "12345678901234567"},
		hs: []string{
			"c711073132333435363738393031323334353637",
			"c80011073132333435363738393031323334353637",
			"c900000011073132333435363738393031323334353637",
		},
	},
	{
		typ: Invalid,
		v:   nil,
		hs:  []string{"c1"},
	},
}

func TestUnpack(t *testing.T) {
	t.Parallel()

	for _, tt := range unpackTests {
		tt := tt
		t.Run(tt.typ.String(), func(t *testing.T) {
			t.Parallel()

			for _, h := range tt.hs {
				p, err := hex.DecodeString(h)
				if err != nil {
					t.Fatalf("decode(%s) returned error %v", h, err)
				}
				d := NewDecoder(bytes.NewReader(p))
				err = d.Unpack()
				if err != nil && tt.typ != Invalid {
					t.Fatalf("unpack(%s) returned %v", h, err)
				}
				if d.Type() != tt.typ {
					t.Fatalf("unpack(%s) returned type %d, want %d", h, d.Type(), tt.typ)
				}
				switch v := tt.v.(type) {
				case int64:
					if d.Int() != v {
						t.Fatalf("unpack(%s) returned %x, want %x", h, d.Int(), v)
					}
				case uint64:
					if d.Uint() != v {
						t.Fatalf("unpack(%s) returned %x, want %x", h, d.Uint(), v)
					}
				case bool:
					if d.Bool() != v {
						t.Fatalf("unpack(%s) returned %v, want %v", h, d.Bool(), v)
					}
				case float64:
					if d.Float() != v {
						t.Fatalf("unpack(%s) returned %v, want %v", h, d.Float(), v)
					}
				case string:
					if d.String() != v {
						t.Fatalf("unpack(%s) returned %q, want %q", h, d.String(), v)
					}
				case extension:
					k, d := d.Extension(), d.String()
					if k != v.k || d != v.d {
						t.Fatalf("unpack(%s) returned (%d, %q) want (%d, %q)", h, k, d, v.k, v.d)
					}
				case nil:
					// nothing to do
				default:
					t.Fatalf("no check for %T", v)
				}
			}
		})
	}
}

func TestUnpackEOF(t *testing.T) {
	t.Parallel()

	for _, tt := range unpackTests {
		tt := tt
		t.Run(tt.typ.String(), func(t *testing.T) {
			t.Parallel()

			for _, h := range tt.hs {
				p, err := hex.DecodeString(h)
				if err != nil {
					t.Fatalf("decode(%s) returned error %v", h, err)
				}

				for i := 1; i < len(p); i++ {
					d := NewDecoder(bytes.NewReader(p[:i]))
					err = d.Unpack()
					if err != io.ErrUnexpectedEOF {
						t.Fatalf("unpack(%s[:%d]) returned %v, want %v", h, i, err, io.ErrUnexpectedEOF)
					}
				}
			}
		})
	}
}
