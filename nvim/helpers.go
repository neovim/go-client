package nvim

import (
	"io"
	"strconv"
)

// QuickfixError represents an item in a quickfix list.
type QuickfixError struct {
	// Buffer number
	Bufnr int `msgpack:"bufnr,omitempty"`

	// Line number in the file.
	LNum int `msgpack:"lnum,omitempty"`

	// Search pattern used to locate the error.
	Pattern string `msgpack:"pattern,omitempty"`

	// Column number (first column is 1).
	Col int `msgpack:"col,omitempty"`

	// When Vcol is != 0,  Col is visual column.
	VCol int `msgpack:"vcol,omitempty"`

	// Error number.
	Nr int `msgpack:"nr,omitempty"`

	// Description of the error.
	Text string `msgpack:"text,omitempty"`

	// Single-character error type, 'E', 'W', etc.
	Type string `msgpack:"type,omitempty"`

	// Name of a file; only used when bufnr is not present or it is invalid.
	FileName string `msgpack:"filename,omitempty"`

	// Valid is non-zero if this is a recognized error message.
	Valid int `msgpack:"valid,omitempty"`
}

// CommandCompletionArgs represents the arguments to a custom command line
// completion function.
//
//  :help :command-completion-custom
type CommandCompletionArgs struct {
	// ArgLead is the leading portion of the argument currently being completed
	// on.
	ArgLead string `msgpack:",array"`

	// CmdLine is the entire command line.
	CmdLine string

	// CursorPosString is decimal representation of the cursor position in
	// bytes.
	CursorPosString string
}

// CursorPos returns the cursor position.
func (a *CommandCompletionArgs) CursorPos() int {
	n, _ := strconv.Atoi(a.CursorPosString)
	return n
}

type bufferReader struct {
	v     *Nvim
	b     Buffer
	lines [][]byte
	err   error
}

// NewBufferReader returns a reader for the specified buffer. If b = 0, then
// the current buffer is used.
func NewBufferReader(v *Nvim, b Buffer) io.Reader {
	return &bufferReader{v: v, b: b}
}

var lineEnd = []byte{'\n'}

func (r *bufferReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.lines == nil {
		r.lines, r.err = r.v.BufferLines(r.b, 0, -1, true)
		if r.err != nil {
			return 0, r.err
		}
	}
	n := 0
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
