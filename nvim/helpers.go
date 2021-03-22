package nvim

import (
	"io"
)

type bufferReader struct {
	err   error
	v     *Nvim
	lines [][]byte
	b     Buffer
}

// compile time check whether the bufferReader implements io.Reader interface.
var _ io.Reader = (*bufferReader)(nil)

// NewBufferReader returns a reader for the specified buffer. If b = 0, then
// the current buffer is used.
func NewBufferReader(v *Nvim, b Buffer) io.Reader {
	return &bufferReader{v: v, b: b}
}

// Read implements io.Reader.
func (r *bufferReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.lines == nil {
		r.lines, r.err = r.v.BufferLines(r.b, 0, -1, true)
		if r.err != nil {
			return 0, r.err
		}
	}
	for {
		if len(r.lines) == 0 {
			r.err = io.EOF
			return n, r.err
		}
		if len(p) == 0 {
			return n, nil
		}

		line0 := r.lines[0]
		if len(line0) == 0 {
			p[0] = '\n'
			p = p[1:]
			n++
			r.lines = r.lines[1:]
			continue
		}

		nn := copy(p, line0)
		n += nn
		p = p[nn:]
		r.lines[0] = line0[nn:]
	}
}
