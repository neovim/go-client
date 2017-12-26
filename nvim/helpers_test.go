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
	"bytes"
	"io"
	"strings"
	"testing"
)

var readerData = []string{
	"\n",
	"hello\nworld\n",
	"blank\n\nline\n",
	"\n1\n22\n333\n",
	"333\n22\n1\n\n",
}

func TestBufferReader(t *testing.T) {
	v, cleanup := newChildProcess(t)
	defer cleanup()
	b, err := v.CurrentBuffer()
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range readerData {
		if err := v.SetBufferLines(b, 0, -1, true, bytes.Split([]byte(strings.TrimSuffix(d, "\n")), []byte{'\n'})); err != nil {
			t.Fatal(err)
		}
		for n := 1; n < 20; n++ {
			var buf bytes.Buffer
			r := NewBufferReader(v, b)
			_, err := io.CopyBuffer(struct{ io.Writer }{&buf}, r, make([]byte, n))
			if err != nil {
				t.Errorf("copy %q with buffer size %d returned error %v", d, n, err)
				continue
			}
			if d != buf.String() {
				t.Errorf("copy %q with buffer size %d = %q", d, n, buf.Bytes())
				continue
			}
		}
	}
}
