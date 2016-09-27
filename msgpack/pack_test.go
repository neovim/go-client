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
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

var packTests = []struct {
	v interface{}
	h string
}{
	{int64(0x0), "00"},
	{int64(0x1), "01"},
	{int64(0x7f), "7f"},
	{int64(0x80), "cc80"},
	{int64(0x7fff), "cd7fff"},
	{int64(0x8000), "cd8000"},
	{int64(0x7fffffff), "ce7fffffff"},
	{int64(0x80000000), "ce80000000"},
	{int64(0x7fffffffffffffff), "cf7fffffffffffffff"},
	{int64(-0x1), "ff"},
	{int64(-0x20), "e0"},
	{int64(-0x21), "d0df"},
	{int64(-0x80), "d080"},
	{int64(-0x81), "d1ff7f"},
	{int64(-0x8000), "d18000"},
	{int64(-0x8001), "d2ffff7fff"},
	{int64(-0x80000000), "d280000000"},
	{int64(-0x80000001), "d3ffffffff7fffffff"},
	{int64(-0x8000000000000000), "d38000000000000000"},
	{uint64(0x0), "00"},
	{uint64(0x1), "01"},
	{uint64(0x7f), "7f"},
	{uint64(0xff), "ccff"},
	{uint64(0x100), "cd0100"},
	{uint64(0xffff), "cdffff"},
	{uint64(0x10000), "ce00010000"},
	{uint64(0xffffffff), "ceffffffff"},
	{uint64(0x100000000), "cf0000000100000000"},
	{uint64(0xffffffffffffffff), "cfffffffffffffffff"},
	{nil, "c0"},
	{true, "c3"},
	{false, "c2"},
	{float64(1.23456), "cb3ff3c0c1fc8f3238"},
	{mapLen(0x0), "80"},
	{mapLen(0x1), "81"},
	{mapLen(0xf), "8f"},
	{mapLen(0x10), "de0010"},
	{mapLen(0xffff), "deffff"},
	{mapLen(0x10000), "df00010000"},
	{mapLen(0xffffffff), "dfffffffff"},
	{arrayLen(0x0), "90"},
	{arrayLen(0x1), "91"},
	{arrayLen(0xf), "9f"},
	{arrayLen(0x10), "dc0010"},
	{arrayLen(0xffff), "dcffff"},
	{arrayLen(0x10000), "dd00010000"},
	{arrayLen(0xffffffff), "ddffffffff"},
	{"", "a0"},
	{"1", "a131"},
	{"1234567890123456789012345678901", "bf31323334353637383930313233343536373839303132333435363738393031"},
	{"12345678901234567890123456789012", "d9203132333435363738393031323334353637383930313233343536373839303132"},
	{[]byte(""), "c400"},
	{[]byte("1"), "c40131"},
	{extension{1, ""}, "c70001"},
	{extension{2, "1"}, "d40231"},
	{extension{3, "12"}, "d5033132"},
	{extension{4, "1234"}, "d60431323334"},
	{extension{5, "12345678"}, "d7053132333435363738"},
	{extension{6, "1234567890123456"}, "d80631323334353637383930313233343536"},
	{extension{7, "12345678901234567"}, "c711073132333435363738393031323334353637"},
}

func TestPack(t *testing.T) {
tests:
	for _, tt := range packTests {

		var arg string
		switch reflect.ValueOf(tt.v).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			arg = fmt.Sprintf("%T %x", tt.v, tt.v)
		default:
			arg = fmt.Sprintf("%T %v", tt.v, tt.v)
		}

		p, err := pack(tt.v)
		if err != nil {
			t.Errorf("pack %s returned error %v", arg, err)
			continue tests
		}
		h := hex.EncodeToString(p)
		if h != tt.h {
			t.Errorf("pack %s returned %s, want %s", arg, h, tt.h)
			continue tests
		}
	}
}
