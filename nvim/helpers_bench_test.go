package nvim

import (
	"bytes"
	"strings"
	"testing"
)

func BenchmarkBufferReader(b *testing.B) {
	v, cleanup := newChildProcess(b)
	defer cleanup()
	cbuf, err := v.CurrentBuffer()
	if err != nil {
		b.Fatal(err)
	}

	benchReaderData := []string{
		"\n",
		"hello\nworld\n",
		"blank\n\nline\n",
		"\n1\n22\n333\n",
		"333\n22\n1\n\n",
	}
	for _, data := range benchReaderData {
		if err := v.SetBufferLines(cbuf, 0, -1, true, bytes.Split([]byte(strings.TrimSuffix(data, "\n")), []byte{'\n'})); err != nil {
			b.Fatal(err)
		}
		b.Run(data, func(b *testing.B) {
			buf := make([]byte, len([]rune(data)))
			r := NewBufferReader(v, cbuf)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = r.Read(buf)
			}

			b.SetBytes(int64(len(buf)))
		})
	}
}
