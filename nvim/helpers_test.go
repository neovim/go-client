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
	v := newChildProcess(t)
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
