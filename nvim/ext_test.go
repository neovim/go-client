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

package nvim

import (
	"encoding/hex"
	"testing"
)

var decodeExtTests = []struct {
	n  int
	hs []string
}{
	{0x0, []string{"00", "d000", "d10000", "d200000000"}},
	{0x1, []string{"01", "d001", "d10001", "d200000001"}},
	{0x7f, []string{"7f", "d07f", "d1007f", "d20000007f"}},
	{0x80, []string{"d10080", "d200000080"}},
	{0x7fff, []string{"d17fff", "d200007fff"}},
	{0x8000, []string{"d200008000"}},
	{0x7fffffff, []string{"d27fffffff"}},
	{-0x1, []string{"ff", "d0ff", "d1ffff", "d2ffffffff"}},
	{-0x20, []string{"e0", "d0e0", "d1ffe0", "d2ffffffe0"}},
	{-0x21, []string{"d0df", "d1ffdf", "d2ffffffdf"}},
	{-0x80, []string{"d080", "d1ff80", "d2ffffff80"}},
	{-0x81, []string{"d1ff7f", "d2ffffff7f"}},
	{-0x8000, []string{"d18000", "d2ffff8000"}},
	{-0x8001, []string{"d2ffff7fff"}},
	{-0x80000000, []string{"d280000000"}},
	{0xff, []string{"ccff", "cd00ff", "ce000000ff"}},
	{0x100, []string{"cd0100", "ce00000100"}},
	{0xffff, []string{"cdffff", "ce0000ffff"}},
	{0x10000, []string{"ce00010000"}},
}

func TestDecodeExt(t *testing.T) {
	for _, tt := range decodeExtTests {
		for _, h := range tt.hs {
			p, err := hex.DecodeString(h)
			if err != nil {
				t.Errorf("hex.DecodeString(%s) returned error %v", h, err)
				continue
			}
			n, err := decodeExt(p)
			if err != nil {
				t.Errorf("decodeExt(%s) returned %v", h, err)
				continue
			}
			if n != tt.n {
				t.Errorf("decodeExt(%s) = %x, want %x", h, n, tt.n)
			}
		}
	}
}

func TestEncodeExt(t *testing.T) {
	for _, tt := range decodeExtTests {
		n, err := decodeExt(encodeExt(tt.n))
		if n != tt.n || err != nil {
			t.Errorf("decodeExt(encodeExt(%x)) returned %x, %v", tt.n, n, err)
		}
	}
}
