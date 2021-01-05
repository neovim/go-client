package msgpack

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestPack(t *testing.T) {
	t.Parallel()

	packTests := map[string]struct {
		// Expected value
		v interface{}
		// Hex encodings of typ, v
		hs string
	}{
		"Bool/True": {
			v:  true,
			hs: "c3",
		},
		"Bool/False": {
			v:  false,
			hs: "c2",
		},
		"Int64/0x0": {
			v:  int64(0x0),
			hs: "00",
		},
		"Int64/0x1": {
			v:  int64(0x1),
			hs: "01",
		},
		"Int64/0x7f": {
			v:  int64(0x7f),
			hs: "7f",
		},
		"Int64/0x80": {
			v:  int64(0x80),
			hs: "cc80",
		},
		"Int64/0x7fff": {
			v:  int64(0x7fff),
			hs: "cd7fff",
		},
		"Int64/0x8000": {
			v:  int64(0x8000),
			hs: "cd8000",
		},
		"Int64/0x7fffffff": {
			v:  int64(0x7fffffff),
			hs: "ce7fffffff",
		},
		"Int64/0x80000000": {
			v:  int64(0x80000000),
			hs: "ce80000000",
		},
		"Int64/0x7fffffffffffffff": {
			v:  int64(0x7fffffffffffffff),
			hs: "cf7fffffffffffffff",
		},
		"Int64/-0x1": {
			v:  int64(-0x1),
			hs: "ff",
		},
		"Int64/-0x20": {
			v:  int64(-0x20),
			hs: "e0",
		},
		"Int64/-0x21": {
			v:  int64(-0x21),
			hs: "d0df",
		},
		"Int64/-0x80": {
			v:  int64(-0x80),
			hs: "d080",
		},
		"Int64/-0x81": {
			v:  int64(-0x81),
			hs: "d1ff7f",
		},
		"Int64/-0x8000": {
			v:  int64(-0x8000),
			hs: "d18000",
		},
		"Int64/-0x8001": {
			v:  int64(-0x8001),
			hs: "d2ffff7fff",
		},
		"Int64/-0x80000000": {
			v:  int64(-0x80000000),
			hs: "d280000000",
		},
		"Int64/-0x80000001": {
			v:  int64(-0x80000001),
			hs: "d3ffffffff7fffffff",
		},
		"Int64/-0x8000000000000000": {
			v:  int64(-0x8000000000000000),
			hs: "d38000000000000000",
		},
		"Uint64/0x0": {
			v:  uint64(0x0),
			hs: "00",
		},
		"Uint64/0x1": {
			v:  uint64(0x1),
			hs: "01",
		},
		"Uint64/0x7f": {
			v:  uint64(0x7f),
			hs: "7f",
		},
		"Uint64/0xff": {
			v:  uint64(0xff),
			hs: "ccff",
		},
		"Uint64/0x100": {
			v:  uint64(0x100),
			hs: "cd0100",
		},
		"Uint64/0xffff": {
			v:  uint64(0xffff),
			hs: "cdffff",
		},
		"Uint64/0x10000": {
			v:  uint64(0x10000),
			hs: "ce00010000",
		},
		"Uint64/0xffffffff": {
			v:  uint64(0xffffffff),
			hs: "ceffffffff",
		},
		"Uint64/0x100000000": {
			v:  uint64(0x100000000),
			hs: "cf0000000100000000",
		},
		"Uint64/0xffffffffffffffff": {
			v:  uint64(0xffffffffffffffff),
			hs: "cfffffffffffffffff",
		},
		"Float64/1.23456": {
			v:  float64(1.23456),
			hs: "cb3ff3c0c1fc8f3238",
		},
		"String/Empty": {
			v:  string(""),
			hs: "a0",
		},
		"String/1": {
			v:  string("1"),
			hs: "a131",
		},
		"String/1234567890123456789012345678901": {
			v:  string("1234567890123456789012345678901"),
			hs: "bf31323334353637383930313233343536373839303132333435363738393031",
		},
		"String/12345678901234567890123456789012": {
			v:  string("12345678901234567890123456789012"),
			hs: "d9203132333435363738393031323334353637383930313233343536373839303132",
		},
		"Binary/Empty": {
			v:  []byte(""),
			hs: "c400",
		},
		"Binary/1": {
			v:  []byte("1"),
			hs: "c40131",
		},
		"MapLen/0x0": {
			v:  mapLen(0x0),
			hs: "80",
		},
		"MapLen/0x1": {
			v:  mapLen(0x1),
			hs: "81",
		},
		"MapLen/0xf": {
			v:  mapLen(0xf),
			hs: "8f",
		},
		"MapLen/0x10": {
			v:  mapLen(0x10),
			hs: "de0010",
		},
		"MapLen/0xffff": {
			v:  mapLen(0xffff),
			hs: "deffff",
		},
		"MapLen/0x10000": {
			v:  mapLen(0x10000),
			hs: "df00010000",
		},
		"MapLen/0xffffffff": {
			v:  mapLen(0xffffffff),
			hs: "dfffffffff",
		},
		"ArrayLen/0x0": {
			v:  arrayLen(0x0),
			hs: "90",
		},
		"ArrayLen/0x1": {
			v:  arrayLen(0x1),
			hs: "91",
		},
		"ArrayLen/0xf": {
			v:  arrayLen(0xf),
			hs: "9f",
		},
		"ArrayLen/0x10": {
			v:  arrayLen(0x10),
			hs: "dc0010",
		},
		"ArrayLen/0xffff": {
			v:  arrayLen(0xffff),
			hs: "dcffff",
		},
		"ArrayLen/0x10000": {
			v:  arrayLen(0x10000),
			hs: "dd00010000",
		},
		"ArrayLen/0xffffffff": {
			v:  arrayLen(0xffffffff),
			hs: "ddffffffff",
		},
		"Extension/1/Empty": {
			v:  extension{1, ""},
			hs: "c70001",
		},
		"Extension/2/1": {
			v:  extension{2, "1"},
			hs: "d40231",
		},
		"Extension/3/12": {
			v:  extension{3, "12"},
			hs: "d5033132",
		},
		"Extension/4/1234": {
			v:  extension{4, "1234"},
			hs: "d60431323334",
		},
		"Extension/5/12345678": {
			v:  extension{5, "12345678"},
			hs: "d7053132333435363738",
		},
		"Extension/6/1234567890123456": {
			v:  extension{6, "1234567890123456"},
			hs: "d80631323334353637383930313233343536",
		},
		"Extension/7/12345678901234567": {
			v:  extension{7, "12345678901234567"},
			hs: "c711073132333435363738393031323334353637",
		},
		"Nil": {
			v:  nil,
			hs: "c0",
		},
	}
	for name, tt := range packTests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var arg string
			switch reflect.ValueOf(tt.v).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				arg = fmt.Sprintf("%T %x", tt.v, tt.v)
			default:
				arg = fmt.Sprintf("%T %v", tt.v, tt.v)
			}

			p, err := pack(tt.v)
			if err != nil {
				t.Fatalf("pack %s returned error %v", arg, err)
			}

			h := hex.EncodeToString(p)
			if h != tt.hs {
				t.Fatalf("pack %s returned %s, want %s", arg, h, tt.hs)
			}
		})
	}
}
