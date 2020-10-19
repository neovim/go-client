package nvim

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newChildProcess(tb testing.TB) (*Nvim, func()) {
	tb.Helper()

	v, err := NewChildProcess(
		ChildProcessArgs("-u", "NONE", "-n", "--embed", "--headless"),
		ChildProcessEnv([]string{}),
		ChildProcessLogf(tb.Logf))
	if err != nil {
		tb.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- v.Serve()
	}()

	return v, func() {
		if err := v.Close(); err != nil {
			tb.Fatal(err)
		}
	}
}

func TestAPI(t *testing.T) {
	t.Parallel()

	v, cleanup := newChildProcess(t)
	defer cleanup()

	t.Run("SimpleHandler", testSimpleHandler(t, v))
	t.Run("Buffer", testBuffer(t, v))
	t.Run("Window", testWindow(t, v))
	t.Run("Tabpage", testTabpage(t, v))
	t.Run("Lines", testLines(t, v))
	t.Run("Var", testVar(t, v))
	t.Run("StructValue", testStructValue(t, v))
	t.Run("Eval", testEval(t, v))
	t.Run("Batch", testBatch(t, v))
	t.Run("CallWithNoArgs", testCallWithNoArgs(t, v))
	t.Run("Mode", testMode(t, v))
	t.Run("ExecLua", testExecLua(t, v))
	t.Run("Highlight", testHighlight(t, v))
	t.Run("BufAttach", testBufAttach(t, v))
	t.Run("VirtualText", testVirtualText(t, v))
	t.Run("FloatingWindow", testFloatingWindow(t, v))
	t.Run("Context", testContext(t, v))
	t.Run("Extmarks", testExtmarks(t, v))
	t.Run("RuntimeFiles", testRuntimeFiles(t, v))
}

func testSimpleHandler(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		cid := v.ChannelID()
		if cid <= 0 {
			t.Fatal("could not get channel id")
		}

		helloHandler := func(s string) (string, error) {
			return "Hello, " + s, nil
		}
		errorHandler := func() error {
			return errors.New("ouch")
		}

		if err := v.RegisterHandler("hello", helloHandler); err != nil {
			t.Fatal(err)
		}
		if err := v.RegisterHandler("error", errorHandler); err != nil {
			t.Fatal(err)
		}
		var result string
		if err := v.Call("rpcrequest", &result, cid, "hello", "world"); err != nil {
			t.Fatal(err)
		}
		if expected := "Hello, world"; result != expected {
			t.Errorf("hello returned %q, want %q", result, expected)
		}

		// Test errors.
		if err := v.Call("execute", &result, fmt.Sprintf("silent! call rpcrequest(%d, 'error')", cid)); err != nil {
			t.Fatal(err)
		}
		if expected := "\nError invoking 'error' on channel 1:\nouch"; result != expected {
			t.Errorf("got error %q, want %q", result, expected)
		}
	}
}

func testBuffer(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		bufs, err := v.Buffers()
		if err != nil {
			t.Fatal(err)
		}
		if len(bufs) != 1 {
			t.Errorf("expected one buf, found %d bufs", len(bufs))
		}
		if bufs[0] == 0 {
			t.Errorf("bufs[0] == 0")
		}
		buf, err := v.CurrentBuffer()
		if err != nil {
			t.Fatal(err)
		}
		if buf != bufs[0] {
			t.Fatalf("buf %v != bufs[0] %v", buf, bufs[0])
		}
		err = v.SetCurrentBuffer(buf)
		if err != nil {
			t.Fatal(err)
		}

		err = v.SetBufferVar(buf, "bvar", "bval")
		if err != nil {
			t.Fatal(err)
		}

		var s string
		err = v.BufferVar(buf, "bvar", &s)
		if err != nil {
			t.Fatal(err)
		}
		if s != "bval" {
			t.Fatalf("expected bvar=bval, got %s", s)
		}

		err = v.DeleteBufferVar(buf, "bvar")
		if err != nil {
			t.Fatal(err)
		}

		s = ""
		err = v.BufferVar(buf, "bvar", &s)
		if err == nil {
			t.Errorf("expected key not found error")
		}
	}
}

func testWindow(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		wins, err := v.Windows()
		if err != nil {
			t.Fatal(err)
		}
		if len(wins) != 1 {
			t.Errorf("expected one win, found %d wins", len(wins))
		}
		if wins[0] == 0 {
			t.Errorf("wins[0] == 0")
		}
		win, err := v.CurrentWindow()
		if err != nil {
			t.Fatal(err)
		}
		if win != wins[0] {
			t.Fatalf("win %v != wins[0] %v", win, wins[0])
		}
		err = v.SetCurrentWindow(win)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testTabpage(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		pages, err := v.Tabpages()
		if err != nil {
			t.Fatal(err)
		}
		if len(pages) != 1 {
			t.Errorf("expected one page, found %d pages", len(pages))
		}
		if pages[0] == 0 {
			t.Errorf("pages[0] == 0")
		}
		page, err := v.CurrentTabpage()
		if err != nil {
			t.Fatal(err)
		}
		if page != pages[0] {
			t.Fatalf("page %v != pages[0] %v", page, pages[0])
		}
		err = v.SetCurrentTabpage(page)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testLines(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		buf, err := v.CurrentBuffer()
		if err != nil {
			t.Fatal(err)
		}
		lines := [][]byte{[]byte("hello"), []byte("world")}
		if err := v.SetBufferLines(buf, 0, -1, true, lines); err != nil {
			t.Fatal(err)
		}
		lines2, err := v.BufferLines(buf, 0, -1, true)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(lines2, lines) {
			t.Fatalf("lines = %+v, want %+v", lines2, lines)
		}
	}
}

func testVar(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		if err := v.SetVar("gvar", "gval"); err != nil {
			t.Fatal(err)
		}
		var value interface{}
		if err := v.Var("gvar", &value); err != nil {
			t.Fatal(err)
		}
		if value != "gval" {
			t.Errorf("got %v, want %q", value, "gval")
		}
		if err := v.SetVar("gvar", ""); err != nil {
			t.Fatal(err)
		}
		value = nil
		if err := v.Var("gvar", &value); err != nil {
			t.Fatal(err)
		}
		if value != "" {
			t.Errorf("got %v, want %q", value, "")
		}
	}
}

func testStructValue(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		var expected, actual struct {
			Str string
			Num int
		}
		expected.Str = "Hello"
		expected.Num = 42
		if err := v.SetVar("structvar", &expected); err != nil {
			t.Fatal(err)
		}
		if err := v.Var("structvar", &actual); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&actual, &expected) {
			t.Errorf("got %+v, want %+v", &actual, &expected)
		}
	}
}

func testEval(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		var a, b string
		if err := v.Eval(`["hello", "world"]`, []*string{&a, &b}); err != nil {
			t.Error(err)
		}
		if a != "hello" || b != "world" {
			t.Errorf("a=%q b=%q, want a=hello b=world", a, b)
		}
	}
}

func testBatch(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		b := v.NewBatch()
		results := make([]int, 128)

		for i := range results {
			b.SetVar(fmt.Sprintf("batch%d", i), i)
		}

		for i := range results {
			b.Var(fmt.Sprintf("batch%d", i), &results[i])
		}

		if err := b.Execute(); err != nil {
			t.Fatal(err)
		}

		for i := range results {
			if results[i] != i {
				t.Fatalf("result[i] = %d, want %d", results[i], i)
			}
		}

		// Reuse batch

		var i int
		b.Var("batch3", &i)
		if err := b.Execute(); err != nil {
			log.Fatal(err)
		}
		if i != 3 {
			t.Fatalf("i = %d, want %d", i, 3)
		}

		// Check for *BatchError

		const errorIndex = 3

		for i := range results {
			results[i] = -1
		}

		for i := range results {
			if i == errorIndex {
				b.Var("batch_bad_var", &results[i])
			} else {
				b.Var(fmt.Sprintf("batch%d", i), &results[i])
			}
		}
		err := b.Execute()
		if e, ok := err.(*BatchError); !ok || e.Index != errorIndex {
			t.Errorf("unxpected error %T %v", e, e)
		}
		// Expect results proceeding error.
		for i := 0; i < errorIndex; i++ {
			if results[i] != i {
				t.Errorf("result[i] = %d, want %d", results[i], i)
				break
			}
		}
		// No results after error.
		for i := errorIndex; i < len(results); i++ {
			if results[i] != -1 {
				t.Errorf("result[i] = %d, want %d", results[i], -1)
				break
			}
		}

		// Execute should return marshal error for argument that cannot be marshaled.
		b.SetVar("batch0", make(chan bool))
		err = b.Execute()
		if err == nil || !strings.Contains(err.Error(), "chan bool") {
			t.Errorf("err = nil, expect error containing text 'chan bool'")
		}

		// Test call with empty argument list.
		var buf Buffer
		b.CurrentBuffer(&buf)
		err = b.Execute()
		if err != nil {
			t.Errorf("GetCurrentBuffer returns err %s: %#v", err, err)
		}
	}
}

func testCallWithNoArgs(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		var wd string
		err := v.Call("getcwd", &wd)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testMode(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		m, err := v.Mode()
		if err != nil {
			t.Fatal(err)
		}
		if m.Mode != "n" {
			t.Errorf("Mode() returned %s, want n", m.Mode)
		}
	}
}

func testExecLua(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		var n int
		err := v.ExecLua("local a, b = ... return a + b", &n, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if n != 3 {
			t.Errorf("Mode() returned %v, want 3", n)
		}
	}
}

func testHighlight(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		cm, err := v.ColorMap()
		if err != nil {
			t.Fatal(err)
		}

		if err := v.Command("hi NewHighlight cterm=underline ctermbg=green guifg=red guibg=yellow guisp=blue gui=bold"); err != nil {
			t.Fatal(err)
		}

		cterm := &HLAttrs{Underline: true, Foreground: -1, Background: 10, Special: -1}
		gui := &HLAttrs{Bold: true, Foreground: cm["Red"], Background: cm["Yellow"], Special: cm["Blue"]}

		var id int
		if err := v.Eval("hlID('NewHighlight')", &id); err != nil {
			t.Fatal(err)
		}
		hl, err := v.HLByID(id, false)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(hl, cterm) {
			t.Errorf("HLByID(id, false)\n got %+v,\nwant %+v", hl, cterm)
		}
		hl, err = v.HLByID(id, true)
		if !reflect.DeepEqual(hl, gui) {
			t.Errorf("HLByID(id, true)\n got %+v,\nwant %+v", hl, gui)
		}
	}
}

func testBufAttach(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		clearBuffer(t, v, 0) // clear curret buffer text

		type ChangedtickEvent struct {
			Buffer     Buffer
			Changetick int64
		}
		type BufLinesEvent struct {
			Buffer      Buffer
			Changetick  int64
			FirstLine   int64
			LastLine    int64
			LineData    string
			IsMultipart bool
		}

		changedtickChan := make(chan *ChangedtickEvent)
		v.RegisterHandler("nvim_buf_changedtick_event", func(changedtickEvent ...interface{}) {
			ev := &ChangedtickEvent{
				Buffer:     changedtickEvent[0].(Buffer),
				Changetick: changedtickEvent[1].(int64),
			}
			changedtickChan <- ev
		})

		bufLinesChan := make(chan *BufLinesEvent)
		v.RegisterHandler("nvim_buf_lines_event", func(bufLinesEvent ...interface{}) {
			ev := &BufLinesEvent{
				Buffer:      bufLinesEvent[0].(Buffer),
				Changetick:  bufLinesEvent[1].(int64),
				FirstLine:   bufLinesEvent[2].(int64),
				LastLine:    bufLinesEvent[3].(int64),
				LineData:    fmt.Sprint(bufLinesEvent[4]),
				IsMultipart: bufLinesEvent[5].(bool),
			}
			bufLinesChan <- ev
		})

		ok, err := v.AttachBuffer(0, false, make(map[string]interface{})) // first 0 arg refers to the current buffer
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal(errors.New("could not attach buffer"))
		}

		changedtickExpected := &ChangedtickEvent{
			Buffer:     1,
			Changetick: 4,
		}
		bufLinesEventExpected := &BufLinesEvent{
			Buffer:      1,
			Changetick:  5,
			FirstLine:   0,
			LastLine:    1,
			LineData:    "[test]",
			IsMultipart: false,
		}

		var numEvent int64 // add and load should be atomically
		errc := make(chan error)
		done := make(chan struct{})
		go func() {
			for {
				select {
				default:
					if atomic.LoadInt64(&numEvent) == 2 { // end buf_attach test when handle 2 event
						done <- struct{}{}
						return
					}
				case changedtick := <-changedtickChan:
					if !reflect.DeepEqual(changedtick, changedtickExpected) {
						errc <- fmt.Errorf("changedtick = %+v, want %+v", changedtick, changedtickExpected)
					}
					atomic.AddInt64(&numEvent, 1)
				case bufLines := <-bufLinesChan:
					if expected := bufLinesEventExpected; !reflect.DeepEqual(bufLines, expected) {
						errc <- fmt.Errorf("bufLines = %+v, want %+v", bufLines, expected)
					}
					atomic.AddInt64(&numEvent, 1)
				}
			}
		}()

		go func() {
			<-done
			close(errc)
		}()

		test := []byte("test")
		if err := v.SetBufferLines(0, 0, -1, true, bytes.Fields(test)); err != nil { // first 0 arg refers to the current buffer
			t.Fatal(err)
		}

		for err := range errc {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func testVirtualText(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		clearBuffer(t, v, Buffer(0)) // clear curret buffer text

		nsID, err := v.CreateNamespace("test_virtual_text")
		if err != nil {
			t.Fatal(err)
		}

		lines := []byte("ping")
		if err := v.SetBufferLines(Buffer(0), 0, -1, true, bytes.Fields(lines)); err != nil {
			t.Fatal(err)
		}

		chunks := []VirtualTextChunk{
			{
				Text:    "pong",
				HLGroup: "String",
			},
		}
		nsID2, err := v.SetBufferVirtualText(Buffer(0), nsID, 0, chunks, make(map[string]interface{}))
		if err != nil {
			t.Fatal(err)
		}

		if got := nsID2; got != nsID {
			t.Fatalf("namespaceID: got %d, want %d", got, nsID)
		}

		if err := v.ClearBufferNamespace(Buffer(0), nsID, 0, -1); err != nil {
			t.Fatal(err)
		}
	}
}

func testFloatingWindow(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		clearBuffer(t, v, 0) // clear curret buffer text
		curwin, err := v.CurrentWindow()
		if err != nil {
			t.Fatal(err)
		}

		wantWidth := 40
		wantHeight := 20

		cfg := &WindowConfig{
			Relative:  "cursor",
			Anchor:    "NW",
			Width:     wantWidth,
			Height:    wantHeight,
			Row:       1,
			Col:       0,
			Focusable: true,
			Style:     "minimal",
		}
		w, err := v.OpenWindow(Buffer(0), true, cfg)
		if err != nil {
			t.Fatal(err)
		}
		if curwin == w {
			t.Fatal("same window number: floating window not focused")
		}

		gotWidth, err := v.WindowWidth(w)
		if err != nil {
			t.Fatal(err)
		}
		if gotWidth != wantWidth {
			t.Fatalf("got %d width but want %d", gotWidth, wantWidth)
		}

		gotHeight, err := v.WindowHeight(w)
		if err != nil {
			t.Fatal(err)
		}
		if gotHeight != wantHeight {
			t.Fatalf("got %d height but want %d", gotHeight, wantHeight)
		}

		batch := v.NewBatch()
		var (
			numberOpt         bool
			relativenumberOpt bool
			cursorlineOpt     bool
			cursorcolumnOpt   bool
			spellOpt          bool
			listOpt           bool
			signcolumnOpt     string
			colorcolumnOpt    string
		)
		batch.WindowOption(w, "number", &numberOpt)
		batch.WindowOption(w, "relativenumber", &relativenumberOpt)
		batch.WindowOption(w, "cursorline", &cursorlineOpt)
		batch.WindowOption(w, "cursorcolumn", &cursorcolumnOpt)
		batch.WindowOption(w, "spell", &spellOpt)
		batch.WindowOption(w, "list", &listOpt)
		batch.WindowOption(w, "signcolumn", &signcolumnOpt)
		batch.WindowOption(w, "colorcolumn", &colorcolumnOpt)
		if err := batch.Execute(); err != nil {
			t.Fatal(err)
		}
		if numberOpt || relativenumberOpt || cursorlineOpt || cursorcolumnOpt || spellOpt || listOpt || signcolumnOpt != "auto" || colorcolumnOpt != "" {
			t.Fatal("expected minimal style")
		}
	}
}

func testContext(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		ctxt, err := v.Context(make(map[string][]string))
		if err != nil {
			t.Fatal(err)
		}

		var result interface{}
		if err := v.LoadContext(ctxt, &result); err != nil {
			t.Fatal(err)
		}
		if result != nil {
			t.Fatal("expected result to nil")
		}
	}
}

func testExtmarks(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		clearBuffer(t, v, 0) // clear curret buffer text

		lines := [][]byte{[]byte("hello"), []byte("world")}
		if err := v.SetBufferLines(Buffer(0), 0, -1, true, lines); err != nil {
			t.Fatal(err)
		}

		nsID, err := v.CreateNamespace("test_extmarks")
		if err != nil {
			t.Fatal(err)
		}
		const (
			extMarkID = 1
			wantLine  = 1
			wantCol   = 3
		)
		gotExtMarkID, err := v.SetBufferExtmark(Buffer(0), nsID, wantLine, wantCol, make(map[string]interface{}))
		if err != nil {
			t.Fatal(err)
		}
		if gotExtMarkID != extMarkID {
			t.Fatalf("got %d extMarkID but want %d", gotExtMarkID, extMarkID)
		}

		extmarks, err := v.BufferExtmarks(Buffer(0), nsID, 0, -1, make(map[string]interface{}))
		if err != nil {
			t.Fatal(err)
		}
		if len(extmarks) > 1 {
			t.Fatalf("expected extmarks length to 1 but %d", len(extmarks))
		}
		if extmarks[0].ExtmarkID != gotExtMarkID {
			t.Fatalf("got %d extMarkID but want %d", extmarks[0].ExtmarkID, extMarkID)
		}
		if extmarks[0].Row != wantLine {
			t.Fatalf("got %d extmarks Row but want %d", extmarks[0].Row, wantLine)
		}
		if extmarks[0].Col != wantCol {
			t.Fatalf("got %d extmarks Col but want %d", extmarks[0].Col, wantCol)
		}

		pos, err := v.BufferExtmarkByID(Buffer(0), nsID, gotExtMarkID, make(map[string]interface{}))
		if err != nil {
			t.Fatal(err)
		}
		if pos[0] != wantLine {
			t.Fatalf("got %d extMark line but want %d", pos[0], wantLine)
		}
		if pos[1] != wantCol {
			t.Fatalf("got %d extMark col but want %d", pos[1], wantCol)
		}

		if err := v.ClearBufferNamespace(Buffer(0), nsID, 0, -1); err != nil {
			t.Fatal(err)
		}
	}
}

func testRuntimeFiles(t *testing.T, v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		files, err := v.RuntimeFiles("doc/*_diff.txt", true)
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(files)
		if len(files) != 2 {
			t.Fatalf("expected 2 length but got %d", len(files))
		}

		var runtimePath string
		if err := v.Eval("$VIMRUNTIME", &runtimePath); err != nil {
			t.Fatal(err)
		}

		viDiff := filepath.Join(runtimePath, "doc", "vi_diff.txt")
		vimDiff := filepath.Join(runtimePath, "doc", "vim_diff.txt")
		want := fmt.Sprintf("%s,%s", viDiff, vimDiff)
		if got := strings.Join(files, ","); !strings.EqualFold(got, want) {
			t.Fatalf("got %s but want %s", got, want)
		}
	}
}

// clearBuffer clears the buffer text.
func clearBuffer(tb testing.TB, v *Nvim, buffer Buffer) {
	tb.Helper()

	if err := v.SetBufferLines(buffer, 0, -1, true, bytes.Fields(nil)); err != nil {
		tb.Fatal(err)
	}
}

func TestDial(t *testing.T) {
	t.Parallel()

	v1, cleanup := newChildProcess(t)
	defer cleanup()

	var addr string
	if err := v1.Eval("$NVIM_LISTEN_ADDRESS", &addr); err != nil {
		t.Fatal(err)
	}

	v2, err := Dial(addr, DialLogf(t.Logf))
	if err != nil {
		t.Fatal(err)
	}
	defer v2.Close()

	if err := v2.SetVar("dial_test", "Hello"); err != nil {
		t.Fatal(err)
	}

	var result string
	if err := v1.Var("dial_test", &result); err != nil {
		t.Fatal(err)
	}

	if expected := "Hello"; result != expected {
		t.Fatalf("got %s, want %s", result, expected)
	}

	if err := v2.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestEmbedded(t *testing.T) {
	t.Parallel()

	v, err := NewEmbedded(&EmbedOptions{
		Args: []string{"-u", "NONE", "-n"},
		Env:  []string{},
		Logf: t.Logf,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer v.Close()

	done := make(chan error, 1)
	go func() {
		done <- v.Serve()
	}()

	var n int
	if err := v.Eval("1+2", &n); err != nil {
		log.Fatal(err)
	}

	if want := 3; n != want {
		log.Fatalf("got %d, want %d", n, want)
	}

	if err := v.Close(); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for serve to exit")
	}
}
