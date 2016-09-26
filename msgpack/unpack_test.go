// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msgpack

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
)

var unpackTests = []struct {
	// Expected type and value.
	typ Type
	v   interface{}

	// Hex encodings of typ, v
	hs []string
}{
	{Int, int64(0x0), []string{"00", "d000", "d10000", "d200000000", "d30000000000000000"}},
	{Int, int64(0x1), []string{"01", "d001", "d10001", "d200000001", "d30000000000000001"}},
	{Int, int64(0x7f), []string{"7f", "d07f", "d1007f", "d20000007f", "d3000000000000007f"}},
	{Int, int64(0x80), []string{"d10080", "d200000080", "d30000000000000080"}},
	{Int, int64(0x7fff), []string{"d17fff", "d200007fff", "d30000000000007fff"}},
	{Int, int64(0x8000), []string{"d200008000", "d30000000000008000"}},
	{Int, int64(0x7fffffff), []string{"d27fffffff", "d3000000007fffffff"}},
	{Int, int64(0x80000000), []string{"d30000000080000000"}},
	{Int, int64(0x7fffffffffffffff), []string{"d37fffffffffffffff"}},
	{Int, int64(-0x1), []string{"ff", "d0ff", "d1ffff", "d2ffffffff", "d3ffffffffffffffff"}},
	{Int, int64(-0x20), []string{"e0", "d0e0", "d1ffe0", "d2ffffffe0", "d3ffffffffffffffe0"}},
	{Int, int64(-0x21), []string{"d0df", "d1ffdf", "d2ffffffdf", "d3ffffffffffffffdf"}},
	{Int, int64(-0x80), []string{"d080", "d1ff80", "d2ffffff80", "d3ffffffffffffff80"}},
	{Int, int64(-0x81), []string{"d1ff7f", "d2ffffff7f", "d3ffffffffffffff7f"}},
	{Int, int64(-0x8000), []string{"d18000", "d2ffff8000", "d3ffffffffffff8000"}},
	{Int, int64(-0x8001), []string{"d2ffff7fff", "d3ffffffffffff7fff"}},
	{Int, int64(-0x80000000), []string{"d280000000", "d3ffffffff80000000"}},
	{Int, int64(-0x80000001), []string{"d3ffffffff7fffffff"}},
	{Int, int64(-0x8000000000000000), []string{"d38000000000000000"}},
	{Uint, uint64(0xff), []string{"ccff", "cd00ff", "ce000000ff", "cf00000000000000ff"}},
	{Uint, uint64(0x100), []string{"cd0100", "ce00000100", "cf0000000000000100"}},
	{Uint, uint64(0xffff), []string{"cdffff", "ce0000ffff", "cf000000000000ffff"}},
	{Uint, uint64(0x10000), []string{"ce00010000", "cf0000000000010000"}},
	{Uint, uint64(0xffffffff), []string{"ceffffffff", "cf00000000ffffffff"}},
	{Uint, uint64(0x100000000), []string{"cf0000000100000000"}},
	{Uint, uint64(0xffffffffffffffff), []string{"cfffffffffffffffff"}},
	{Nil, nil, []string{"c0"}},
	{Bool, true, []string{"c3"}},
	{Bool, false, []string{"c2"}},
	{Float, float64(123456), []string{"ca47f12000"}},
	{Float, float64(1.23456), []string{"cb3ff3c0c1fc8f3238"}},
	{MapLen, int64(0x0), []string{"80", "de0000", "df00000000"}},
	{MapLen, int64(0x1), []string{"81", "de0001", "df00000001"}},
	{MapLen, int64(0xf), []string{"8f", "de000f", "df0000000f"}},
	{MapLen, int64(0x10), []string{"de0010", "df00000010"}},
	{MapLen, int64(0xffff), []string{"deffff", "df0000ffff"}},
	{MapLen, int64(0x10000), []string{"df00010000"}},
	{MapLen, int64(0xffffffff), []string{"dfffffffff"}},
	{ArrayLen, int64(0x0), []string{"90", "dc0000", "dd00000000"}},
	{ArrayLen, int64(0x1), []string{"91", "dc0001", "dd00000001"}},
	{ArrayLen, int64(0xf), []string{"9f", "dc000f", "dd0000000f"}},
	{ArrayLen, int64(0x10), []string{"dc0010", "dd00000010"}},
	{ArrayLen, int64(0xffff), []string{"dcffff", "dd0000ffff"}},
	{ArrayLen, int64(0x10000), []string{"dd00010000"}},
	{ArrayLen, int64(0xffffffff), []string{"ddffffffff"}},
	{String, "", []string{"a0", "d900", "da0000", "db00000000"}},
	{String, "1", []string{"a131", "d90131", "da000131", "db0000000131"}},
	{String, "1234567890123456789012345678901", []string{
		"bf31323334353637383930313233343536373839303132333435363738393031",
		"d91f31323334353637383930313233343536373839303132333435363738393031",
		"da001f31323334353637383930313233343536373839303132333435363738393031",
		"db0000001f31323334353637383930313233343536373839303132333435363738393031"}},
	{String, "12345678901234567890123456789012", []string{
		"d9203132333435363738393031323334353637383930313233343536373839303132",
		"da00203132333435363738393031323334353637383930313233343536373839303132",
		"db000000203132333435363738393031323334353637383930313233343536373839303132"}},
	{Binary, "", []string{"c400", "c50000", "c600000000"}},
	{Binary, "1", []string{"c40131", "c5000131", "c60000000131"}},
	{Extension, extension{1, ""}, []string{"c70001", "c8000001", "c90000000001"}},
	{Extension, extension{2, "1"}, []string{"d40231", "c7010231", "c800010231", "c9000000010231"}},
	{Extension, extension{3, "12"}, []string{"d5033132", "c702033132", "c80002033132", "c900000002033132"}},
	{Extension, extension{4, "1234"}, []string{
		"d60431323334",
		"c7040431323334",
		"c800040431323334",
		"c9000000040431323334"}},
	{Extension, extension{5, "12345678"}, []string{
		"d7053132333435363738",
		"c708053132333435363738",
		"c80008053132333435363738",
		"c900000008053132333435363738"}},
	{Extension, extension{6, "1234567890123456"}, []string{
		"d80631323334353637383930313233343536",
		"c7100631323334353637383930313233343536",
		"c800100631323334353637383930313233343536",
		"c9000000100631323334353637383930313233343536"}},
	{Extension, extension{7, "12345678901234567"}, []string{
		"c711073132333435363738393031323334353637",
		"c80011073132333435363738393031323334353637",
		"c900000011073132333435363738393031323334353637"}},
	{Invalid, nil, []string{
		"c1",
	},
	},
}

func TestUnpack(t *testing.T) {
	for _, tt := range unpackTests {
	tests:
		for _, h := range tt.hs {
			p, err := hex.DecodeString(h)
			if err != nil {
				t.Errorf("decode(%s) returned error %v", h, err)
				continue tests
			}
			d := NewDecoder(bytes.NewReader(p))
			err = d.Unpack()
			if err != nil && tt.typ != Invalid {
				t.Errorf("unpack(%s) returned %v", h, err)
				continue tests
			}
			if d.Type() != tt.typ {
				t.Errorf("unpack(%s) returned type %d, want %d", h, d.Type(), tt.typ)
				continue tests
			}
			switch v := tt.v.(type) {
			case int64:
				if d.Int() != v {
					t.Errorf("unpack(%s) returned %x, want %x", h, d.Int(), v)
					continue tests
				}
			case uint64:
				if d.Uint() != v {
					t.Errorf("unpack(%s) returned %x, want %x", h, d.Uint(), v)
					continue tests
				}
			case bool:
				if d.Bool() != v {
					t.Errorf("unpack(%s) returned %v, want %v", h, d.Bool(), v)
					continue tests
				}
			case float64:
				if d.Float() != v {
					t.Errorf("unpack(%s) returned %v, want %v", h, d.Float(), v)
					continue tests
				}
			case string:
				if d.String() != v {
					t.Errorf("unpack(%s) returned %q, want %q", h, d.String(), v)
				}
			case extension:
				k, d := d.Extension(), d.String()
				if k != v.k || d != v.d {
					t.Errorf("unpack(%s) returned (%d, %q) want (%d, %q)", h, k, d, v.k, v.d)
				}
			case nil:
				// do nothing
			default:
				t.Errorf("no check for %T", v)
				continue tests
			}
		}
	}
}

func TestUnpackEOF(t *testing.T) {
	for _, tt := range unpackTests {
		for _, h := range tt.hs {
			p, err := hex.DecodeString(h)
			if err != nil {
				t.Errorf("decode(%s) returned error %v", h, err)
				continue
			}
			for i := 1; i < len(p); i++ {
				d := NewDecoder(bytes.NewReader(p[:i]))
				err = d.Unpack()
				if err != io.ErrUnexpectedEOF {
					t.Errorf("unpack(%s[:%d]) returned %v, want %v", h, i, err, io.ErrUnexpectedEOF)
				}
			}
		}
	}
}
