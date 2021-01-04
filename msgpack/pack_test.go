package msgpack

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestPack(t *testing.T) {
	t.Parallel()

	packTests := []struct {
		v  interface{}
		hs string
	}{
		{
			v:  int64(0x0),
			hs: "00",
		},
		{
			v:  int64(0x1),
			hs: "01",
		},
		{
			v:  int64(0x7f),
			hs: "7f",
		},
		{
			v:  int64(0x80),
			hs: "cc80",
		},
		{
			v:  int64(0x7fff),
			hs: "cd7fff",
		},
		{
			v:  int64(0x8000),
			hs: "cd8000",
		},
		{
			v:  int64(0x7fffffff),
			hs: "ce7fffffff",
		},
		{
			v:  int64(0x80000000),
			hs: "ce80000000",
		},
		{
			v:  int64(0x7fffffffffffffff),
			hs: "cf7fffffffffffffff",
		},
		{
			v:  int64(-0x1),
			hs: "ff",
		},
		{
			v:  int64(-0x20),
			hs: "e0",
		},
		{
			v:  int64(-0x21),
			hs: "d0df",
		},
		{
			v:  int64(-0x80),
			hs: "d080",
		},
		{
			v:  int64(-0x81),
			hs: "d1ff7f",
		},
		{
			v:  int64(-0x8000),
			hs: "d18000",
		},
		{
			v:  int64(-0x8001),
			hs: "d2ffff7fff",
		},
		{
			v:  int64(-0x80000000),
			hs: "d280000000",
		},
		{
			v:  int64(-0x80000001),
			hs: "d3ffffffff7fffffff",
		},
		{
			v:  int64(-0x8000000000000000),
			hs: "d38000000000000000",
		},
		{
			v:  uint64(0x0),
			hs: "00",
		},
		{
			v:  uint64(0x1),
			hs: "01",
		},
		{
			v:  uint64(0x7f),
			hs: "7f",
		},
		{
			v:  uint64(0xff),
			hs: "ccff",
		},
		{
			v:  uint64(0x100),
			hs: "cd0100",
		},
		{
			v:  uint64(0xffff),
			hs: "cdffff",
		},
		{
			v:  uint64(0x10000),
			hs: "ce00010000",
		},
		{
			v:  uint64(0xffffffff),
			hs: "ceffffffff",
		},
		{
			v:  uint64(0x100000000),
			hs: "cf0000000100000000",
		},
		{
			v:  uint64(0xffffffffffffffff),
			hs: "cfffffffffffffffff",
		},
		{
			v:  nil,
			hs: "c0",
		},
		{
			v:  true,
			hs: "c3",
		},
		{
			v:  false,
			hs: "c2",
		},
		{
			v:  float64(1.23456),
			hs: "cb3ff3c0c1fc8f3238",
		},
		{
			v:  mapLen(0x0),
			hs: "80",
		},
		{
			v:  mapLen(0x1),
			hs: "81",
		},
		{
			v:  mapLen(0xf),
			hs: "8f",
		},
		{
			v:  mapLen(0x10),
			hs: "de0010",
		},
		{
			v:  mapLen(0xffff),
			hs: "deffff",
		},
		{
			v:  mapLen(0x10000),
			hs: "df00010000",
		},
		{
			v:  mapLen(0xffffffff),
			hs: "dfffffffff",
		},
		{
			v:  arrayLen(0x0),
			hs: "90",
		},
		{
			v:  arrayLen(0x1),
			hs: "91",
		},
		{
			v:  arrayLen(0xf),
			hs: "9f",
		},
		{
			v:  arrayLen(0x10),
			hs: "dc0010",
		},
		{
			v:  arrayLen(0xffff),
			hs: "dcffff",
		},
		{
			v:  arrayLen(0x10000),
			hs: "dd00010000",
		},
		{
			v:  arrayLen(0xffffffff),
			hs: "ddffffffff",
		},
		{
			v:  "",
			hs: "a0",
		},
		{
			v:  "1",
			hs: "a131",
		},
		{
			v:  "1234567890123456789012345678901",
			hs: "bf31323334353637383930313233343536373839303132333435363738393031",
		},
		{
			v:  "12345678901234567890123456789012",
			hs: "d9203132333435363738393031323334353637383930313233343536373839303132",
		},
		{
			v:  []byte(""),
			hs: "c400",
		},
		{
			v:  []byte("1"),
			hs: "c40131",
		},
		{
			v:  extension{1, ""},
			hs: "c70001",
		},
		{
			v:  extension{2, "1"},
			hs: "d40231",
		},
		{
			v:  extension{3, "12"},
			hs: "d5033132",
		},
		{
			v:  extension{4, "1234"},
			hs: "d60431323334",
		},
		{
			v:  extension{5, "12345678"},
			hs: "d7053132333435363738",
		},
		{
			v:  extension{6, "1234567890123456"},
			hs: "d80631323334353637383930313233343536",
		},
		{
			v:  extension{7, "12345678901234567"},
			hs: "c711073132333435363738393031323334353637",
		},
	}
	for _, tt := range packTests {
		tt := tt
		t.Run(fmt.Sprintf("%s", tt.v), func(t *testing.T) {
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
