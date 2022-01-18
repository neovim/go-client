package nvim

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

type version struct {
	Major int64
	Minor int64
	Patch int64
}

var nvimVersion version

func parseVersion(tb testing.TB, version string) (major, minor, patch int64) {
	tb.Helper()

	version = strings.TrimPrefix(version, "v")
	vpair := strings.Split(version, ".")
	if len(vpair) != 3 {
		tb.Fatalf("could not parse neovim version: %s", version)
	}

	var err error
	major, err = strconv.ParseInt(vpair[0], 10, 0)
	if err != nil {
		tb.Fatal(err)
	}
	minor, err = strconv.ParseInt(vpair[1], 10, 0)
	if err != nil {
		tb.Fatal(err)
	}
	patch, err = strconv.ParseInt(vpair[2], 10, 0)
	if err != nil {
		tb.Fatal(err)
	}

	return major, minor, patch
}

func skipVersion(tb testing.TB, version string) {
	major, minor, patch := parseVersion(tb, version)

	const skipFmt = "SKIP: current neovim version v%d.%d.%d but expected version %s"
	if nvimVersion.Major < major || nvimVersion.Minor < minor || nvimVersion.Patch < patch {
		tb.Skipf(skipFmt, nvimVersion.Major, nvimVersion.Minor, nvimVersion.Patch, version)
	}
}

// clearBuffer clears the buffer lines.
func clearBuffer(tb testing.TB, v *Nvim, buffer Buffer) {
	tb.Helper()

	if err := v.SetBufferLines(buffer, 0, -1, true, bytes.Fields(nil)); err != nil {
		tb.Fatal(err)
	}
}

var channelID int64

func TestAPI(t *testing.T) {
	t.Parallel()

	v := newChildProcess(t)

	apiInfo, err := v.APIInfo()
	if err != nil {
		t.Fatal(err)
	}
	if len(apiInfo) != 2 {
		t.Fatalf("unknown APIInfo: %#v", apiInfo)
	}

	var ok bool
	channelID, ok = apiInfo[0].(int64)
	if !ok {
		t.Fatalf("apiInfo[0] is not int64 type: %T", apiInfo[0])
	}

	info, ok := apiInfo[1].(map[string]interface{})
	if !ok {
		t.Fatalf("apiInfo[1] is not map[string]interface{} type: %T", apiInfo[1])
	}
	infoV := info["version"].(map[string]interface{})
	nvimVersion.Major = infoV["major"].(int64)
	nvimVersion.Minor = infoV["minor"].(int64)
	nvimVersion.Patch = infoV["patch"].(int64)

	t.Run("BufAttach", testBufAttach(v))
	t.Run("APIInfo", testAPIInfo(v))
	t.Run("SimpleHandler", testSimpleHandler(v))
	t.Run("Buffer", testBuffer(v))
	t.Run("Window", testWindow(v))
	t.Run("Tabpage", testTabpage(v))
	t.Run("Lines", testLines(v))
	t.Run("Command", testCommand(v))
	t.Run("Var", testVar(v))
	t.Run("Message", testMessage(v))
	t.Run("Key", testKey(v))
	t.Run("Eval", testEval(v))
	t.Run("Batch", testBatch(v))
	t.Run("Mode", testMode(v))
	t.Run("ExecLua", testExecLua(v))
	t.Run("Highlight", testHighlight(v))
	t.Run("VirtualText", testVirtualText(v))
	t.Run("FloatingWindow", testFloatingWindow(v))
	t.Run("Context", testContext(v))
	t.Run("Extmarks", testExtmarks(v))
	t.Run("Runtime", testRuntime(v))
	t.Run("Namespace", testNamespace(v))
	t.Run("PutPaste", testPutPaste(v))
	t.Run("Options", testOptions(v))
	t.Run("AllOptionsInfo", testAllOptionsInfo(v))
	t.Run("OptionsInfo", testOptionsInfo(v))
	t.Run("OptionsValue", testOptionsValue(v))
	t.Run("OpenTerm", testTerm(v))
	t.Run("ChannelClientInfo", testChannelClientInfo(v))
	t.Run("UI", testUI(v))
	t.Run("Proc", testProc(v))
	t.Run("Mark", testMark(v))
	t.Run("StatusLine", testStatusLine(v))
}

func testBufAttach(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			clearBuffer(t, v, 0) // clear curret buffer text

			changedtickChan := make(chan *ChangedtickEvent)
			v.RegisterHandler(EventBufChangedtick, func(changedtickEvent ...interface{}) {
				ev := &ChangedtickEvent{
					Buffer:     changedtickEvent[0].(Buffer),
					Changetick: changedtickEvent[1].(int64),
				}
				changedtickChan <- ev
			})

			bufLinesChan := make(chan *BufLinesEvent)
			v.RegisterHandler(EventBufLines, func(bufLinesEvent ...interface{}) {
				ev := &BufLinesEvent{
					Buffer:      bufLinesEvent[0].(Buffer),
					Changetick:  bufLinesEvent[1].(int64),
					FirstLine:   bufLinesEvent[2].(int64),
					LastLine:    bufLinesEvent[3].(int64),
					IsMultipart: bufLinesEvent[5].(bool),
				}
				for _, line := range bufLinesEvent[4].([]interface{}) {
					ev.LineData = append(ev.LineData, line.(string))
				}
				bufLinesChan <- ev
			})

			bufDetachChan := make(chan *BufDetachEvent)
			v.RegisterHandler(EventBufDetach, func(bufDetachEvent ...interface{}) {
				ev := &BufDetachEvent{
					Buffer: bufDetachEvent[0].(Buffer),
				}
				bufDetachChan <- ev
			})

			ok, err := v.AttachBuffer(0, false, make(map[string]interface{})) // first 0 arg refers to the current buffer
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Fatal("could not attach buffer")
			}

			changedtickExpected := &ChangedtickEvent{
				Buffer:     Buffer(1),
				Changetick: 3,
			}
			bufLinesEventExpected := &BufLinesEvent{
				Buffer:      Buffer(1),
				Changetick:  4,
				FirstLine:   0,
				LastLine:    1,
				LineData:    []string{"foo", "bar", "baz", "qux", "quux", "quuz"},
				IsMultipart: false,
			}
			bufDetachEventExpected := &BufDetachEvent{
				Buffer: Buffer(1),
			}

			var numEvent int64 // add and load should be atomically
			errc := make(chan error)
			done := make(chan struct{})
			go func() {
				for {
					select {
					default:
						if atomic.LoadInt64(&numEvent) == 3 { // end buf_attach test when handle 2 event
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
					case detach := <-bufDetachChan:
						if expected := bufDetachEventExpected; !reflect.DeepEqual(detach, expected) {
							errc <- fmt.Errorf("bufDetach = %+v, want %+v", detach, expected)
						}
						atomic.AddInt64(&numEvent, 1)
					}
				}
			}()

			go func() {
				<-done
				close(errc)
			}()

			test := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz"), []byte("qux"), []byte("quux"), []byte("quuz")}
			if err := v.SetBufferLines(Buffer(0), 0, -1, true, test); err != nil { // first 0 arg refers to the current buffer
				t.Fatal(err)
			}

			if detached, err := v.DetachBuffer(Buffer(0)); err != nil || !detached {
				t.Fatal(err)
			}

			for err := range errc {
				if err != nil {
					t.Fatal(err)
				}
			}
		})

		t.Run("Batch", func(t *testing.T) {
			clearBuffer(t, v, 0) // clear curret buffer text

			changedtickChan := make(chan *ChangedtickEvent)
			v.RegisterHandler(EventBufChangedtick, func(changedtickEvent ...interface{}) {
				ev := &ChangedtickEvent{
					Buffer:     changedtickEvent[0].(Buffer),
					Changetick: changedtickEvent[1].(int64),
				}
				changedtickChan <- ev
			})

			bufLinesChan := make(chan *BufLinesEvent)
			v.RegisterHandler(EventBufLines, func(bufLinesEvent ...interface{}) {
				ev := &BufLinesEvent{
					Buffer:      bufLinesEvent[0].(Buffer),
					Changetick:  bufLinesEvent[1].(int64),
					FirstLine:   bufLinesEvent[2].(int64),
					LastLine:    bufLinesEvent[3].(int64),
					IsMultipart: bufLinesEvent[5].(bool),
				}
				for _, line := range bufLinesEvent[4].([]interface{}) {
					ev.LineData = append(ev.LineData, line.(string))
				}
				bufLinesChan <- ev
			})

			bufDetachChan := make(chan *BufDetachEvent)
			v.RegisterHandler(EventBufDetach, func(bufDetachEvent ...interface{}) {
				ev := &BufDetachEvent{
					Buffer: bufDetachEvent[0].(Buffer),
				}
				bufDetachChan <- ev
			})

			b := v.NewBatch()

			var attached bool
			b.AttachBuffer(0, false, make(map[string]interface{}), &attached) // first 0 arg refers to the current buffer
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !attached {
				t.Fatal("could not attach buffer")
			}

			changedtickExpected := &ChangedtickEvent{
				Buffer:     Buffer(1),
				Changetick: 5,
			}
			bufLinesEventExpected := &BufLinesEvent{
				Buffer:      Buffer(1),
				Changetick:  6,
				FirstLine:   0,
				LastLine:    1,
				LineData:    []string{"foo", "bar", "baz", "qux", "quux", "quuz"},
				IsMultipart: false,
			}
			bufDetachEventExpected := &BufDetachEvent{
				Buffer: Buffer(1),
			}

			var numEvent int64 // add and load should be atomically
			errc := make(chan error)
			done := make(chan struct{})
			go func() {
				for {
					select {
					default:
						if atomic.LoadInt64(&numEvent) == 3 { // end buf_attach test when handle 2 event
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
					case detach := <-bufDetachChan:
						if expected := bufDetachEventExpected; !reflect.DeepEqual(detach, expected) {
							errc <- fmt.Errorf("bufDetach = %+v, want %+v", detach, expected)
						}
						atomic.AddInt64(&numEvent, 1)
					}
				}
			}()

			go func() {
				<-done
				close(errc)
			}()

			test := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz"), []byte("qux"), []byte("quux"), []byte("quuz")}
			if err := v.SetBufferLines(Buffer(0), 0, -1, true, test); err != nil { // first 0 arg refers to the current buffer
				t.Fatal(err)
			}

			var detached bool
			b.DetachBuffer(Buffer(0), &detached)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			for err := range errc {
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func testAPIInfo(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			apiinfo, err := v.APIInfo()
			if err != nil {
				t.Fatal(err)
			}
			if len(apiinfo) == 0 {
				t.Fatal("expected apiinfo is non-nil")
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var apiinfo []interface{}
			b.APIInfo(&apiinfo)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if len(apiinfo) == 0 {
				t.Fatal("expected apiinfo is non-nil")
			}
		})
	}
}

func testSimpleHandler(v *Nvim) func(*testing.T) {
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
			t.Fatalf("hello returned %q, want %q", result, expected)
		}

		// Test errors.
		if err := v.Call("execute", &result, fmt.Sprintf("silent! call rpcrequest(%d, 'error')", cid)); err != nil {
			t.Fatal(err)
		}
		if expected := "\nError invoking 'error' on channel 1:\nouch"; result != expected {
			t.Fatalf("got error %q, want %q", result, expected)
		}
	}
}

func testBuffer(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Run("BufferName", func(t *testing.T) {
				cwd, err := os.Getwd() // buffer name is full path
				if err != nil {
					t.Fatal(err)
				}
				wantBufName := filepath.Join(cwd, "test_buffer")
				if err := v.SetBufferName(Buffer(0), wantBufName); err != nil {
					t.Fatal(err)
				}

				bufName, err := v.BufferName(Buffer(0))
				if err != nil {
					t.Fatal(err)
				}

				if bufName != wantBufName {
					t.Fatalf("want %s buffer name but got %s", wantBufName, bufName)
				}

				t.Cleanup(func() {
					// cleanup buffer name
					if err := v.SetBufferName(Buffer(0), ""); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("Buffers", func(t *testing.T) {
				bufs, err := v.Buffers()
				if err != nil {
					t.Fatal(err)
				}
				if len(bufs) != 2 {
					t.Fatalf("expected one buf, found %d bufs", len(bufs))
				}
				if bufs[0] == 0 {
					t.Fatalf("bufs[0] is not %q: %q", bufs[0], Buffer(0))
				}

				buf, err := v.CurrentBuffer()
				if err != nil {
					t.Fatal(err)
				}
				if buf != bufs[0] {
					t.Fatalf("buf is not bufs[0]: buf %v, bufs[0]: %v", buf, bufs[0])
				}

				const want = "Buffer:1"
				if got := buf.String(); got != want {
					t.Fatalf("buf.String() = %s, want %s", got, want)
				}

				if err := v.SetCurrentBuffer(buf); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("Var", func(t *testing.T) {
				buf, err := v.CurrentBuffer()
				if err != nil {
					t.Fatal(err)
				}

				const (
					varkey = "bvar"
					varVal = "bval"
				)
				if err := v.SetBufferVar(buf, varkey, varVal); err != nil {
					t.Fatal(err)
				}

				var s string
				if err := v.BufferVar(buf, varkey, &s); err != nil {
					t.Fatal(err)
				}
				if s != "bval" {
					t.Fatalf("expected %s=%s, got %s", s, varkey, varVal)
				}

				if err := v.DeleteBufferVar(buf, varkey); err != nil {
					t.Fatal(err)
				}

				s = "" // reuse
				if err := v.BufferVar(buf, varkey, &s); err == nil {
					t.Fatalf("expected %s not found but error is nil: err: %#v", varkey, err)
				}
			})

			t.Run("Delete", func(t *testing.T) {
				buf, err := v.CreateBuffer(true, true)
				if err != nil {
					t.Fatal(err)
				}

				deleteBufferOpts := map[string]bool{
					"force":  true,
					"unload": false,
				}
				if err := v.DeleteBuffer(buf, deleteBufferOpts); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("ChangeTick", func(t *testing.T) {
				buf, err := v.CreateBuffer(true, true)
				if err != nil {
					t.Fatal(err)
				}

				// 1 changedtick
				lines := [][]byte{[]byte("hello"), []byte("world")}
				if err := v.SetBufferLines(buf, 0, -1, true, lines); err != nil {
					t.Fatal(err)
				}

				// 2 changedtick
				const wantChangedTick = 2
				changedTick, err := v.BufferChangedTick(buf)
				if err != nil {
					t.Fatal(err)
				}
				if changedTick != wantChangedTick {
					t.Fatalf("got %d changedTick but want %d", changedTick, wantChangedTick)
				}

				// cleanup buffer
				deleteBufferOpts := map[string]bool{
					"force":  true,
					"unload": false,
				}
				if err := v.DeleteBuffer(buf, deleteBufferOpts); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("SetCurrentDirectory", func(t *testing.T) {
				wantDir, err := os.UserHomeDir()
				if err != nil {
					t.Fatal(err)
				}

				if err := v.SetCurrentDirectory(wantDir); err != nil {
					t.Fatal(err)
				}

				var got string
				if err := v.Eval(`getcwd()`, &got); err != nil {
					t.Fatal(err)
				}

				if got != wantDir {
					t.Fatalf("SetCurrentDirectory(%s) = %s, want: %s", wantDir, got, wantDir)
				}
			})

			t.Run("BufferCommands", func(t *testing.T) {
				commands, err := v.BufferCommands(Buffer(0), make(map[string]interface{}))
				if err != nil {
					t.Fatal(err)
				}

				if len(commands) > 0 {
					t.Fatalf("expected commands empty but non-zero: %#v", commands)
				}
			})

			t.Run("BufferOption", func(t *testing.T) {
				var cindent bool
				if err := v.BufferOption(Buffer(0), "cindent", &cindent); err != nil {
					t.Fatal(err)
				}

				if cindent {
					t.Fatalf("expected cindent is false but got %t", cindent)
				}

				if err := v.SetBufferOption(Buffer(0), "cindent", true); err != nil {
					t.Fatal(err)
				}

				if err := v.BufferOption(Buffer(0), "cindent", &cindent); err != nil {
					t.Fatal(err)
				}

				if !cindent {
					t.Fatalf("expected cindent is true but got %t", cindent)
				}

				t.Cleanup(func() {
					// cleanup cindent option
					if err := v.SetBufferOption(Buffer(0), "cindent", false); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("IsBufferLoaded", func(t *testing.T) {
				loaded, err := v.IsBufferLoaded(Buffer(0))
				if err != nil {
					t.Fatal(err)
				}

				if !loaded {
					t.Fatalf("expected buffer is loaded but got %t", loaded)
				}
			})

			t.Run("IsBufferValid", func(t *testing.T) {
				valid, err := v.IsBufferValid(Buffer(0))
				if err != nil {
					t.Fatal(err)
				}

				if !valid {
					t.Fatalf("expected buffer is valid but got %t", valid)
				}
			})

			t.Run("BufferMark", func(t *testing.T) {
				lines := [][]byte{
					[]byte("a"),
					[]byte("bit of"),
					[]byte("text"),
				}
				if err := v.SetBufferLines(Buffer(0), -1, -1, true, lines); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					clearBuffer(t, v, Buffer(0))
				})

				if err := v.SetWindowCursor(Window(0), [2]int{3, 4}); err != nil {
					t.Fatal(err)
				}

				const (
					mark     = "V"
					wantLine = 3
					wantCol  = 0
				)
				set, err := v.SetBufferMark(Buffer(0), mark, wantLine, wantCol, make(map[string]interface{}))
				if err != nil {
					t.Fatal(err)
				}
				if !set {
					t.Fatalf("could not set %s mark", mark)
				}

				pos, err := v.BufferMark(Buffer(0), mark)
				if err != nil {
					t.Fatal(err)
				}
				if pos[0] != wantLine {
					t.Fatalf("got %d extMark line but want %d", pos[0], wantLine)
				}
				if pos[1] != wantCol {
					t.Fatalf("got %d extMark col but want %d", pos[1], wantCol)
				}

				deleted, err := v.DeleteBufferMark(Buffer(0), mark)
				if err != nil {
					t.Fatal(err)
				}
				if !deleted {
					t.Fatalf("could not delete %s mark", mark)
				}

				pos2, err := v.BufferMark(Buffer(0), mark)
				if err != nil {
					t.Fatal(err)
				}
				if pos2[0] != 0 || pos2[1] != 0 {
					t.Fatalf("got %d mark but want zero", pos2)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("BufferName", func(t *testing.T) {
				b := v.NewBatch()

				cwd, err := os.Getwd() // buffer name is full path
				if err != nil {
					t.Fatal(err)
				}
				wantBufName := filepath.Join(cwd, "test_buffer")
				b.SetBufferName(Buffer(0), wantBufName)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var bufName string
				b.BufferName(Buffer(0), &bufName)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if bufName != wantBufName {
					t.Fatalf("want %s buffer name but got %s", wantBufName, bufName)
				}

				t.Cleanup(func() {
					// cleanup buffer name
					b.SetBufferName(Buffer(0), "")
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("Buffers", func(t *testing.T) {
				b := v.NewBatch()

				var bufs []Buffer
				b.Buffers(&bufs)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if len(bufs) != 2 {
					t.Fatalf("expected one buf, found %d bufs", len(bufs))
				}
				if bufs[0] == Buffer(0) {
					t.Fatalf("bufs[0] is not %q: %q", bufs[0], Buffer(0))
				}

				var buf Buffer
				b.CurrentBuffer(&buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if buf != bufs[0] {
					t.Fatalf("buf is not bufs[0]: buf %v, bufs[0]: %v", buf, bufs[0])
				}

				const want = "Buffer:1"
				if got := buf.String(); got != want {
					t.Fatalf("buf.String() = %s, want %s", got, want)
				}

				b.SetCurrentBuffer(buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("Var", func(t *testing.T) {
				b := v.NewBatch()

				var buf Buffer
				b.CurrentBuffer(&buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				const (
					varkey = "bvar"
					varVal = "bval"
				)
				b.SetBufferVar(buf, varkey, varVal)
				var s string
				b.BufferVar(buf, varkey, &s)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if s != varVal {
					t.Fatalf("expected bvar=bval, got %s", s)
				}

				b.DeleteBufferVar(buf, varkey)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				s = "" // reuse
				b.BufferVar(buf, varkey, &s)
				if err := b.Execute(); err == nil {
					t.Fatalf("expected %s not found but error is nil: err: %#v", varkey, err)
				}
			})

			t.Run("Delete", func(t *testing.T) {
				b := v.NewBatch()

				var buf Buffer
				b.CreateBuffer(true, true, &buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				deleteBufferOpts := map[string]bool{
					"force":  true,
					"unload": false,
				}
				b.DeleteBuffer(buf, deleteBufferOpts)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("ChangeTick", func(t *testing.T) {
				b := v.NewBatch()

				var buf Buffer
				b.CreateBuffer(true, true, &buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// 1 changedtick
				lines := [][]byte{[]byte("hello"), []byte("world")}
				b.SetBufferLines(buf, 0, -1, true, lines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// 2 changedtick
				const wantChangedTick = 2
				var changedTick int
				b.BufferChangedTick(buf, &changedTick)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if changedTick != wantChangedTick {
					t.Fatalf("got %d changedTick but want %d", changedTick, wantChangedTick)
				}

				// cleanup buffer
				deleteBufferOpts := map[string]bool{
					"force":  true,
					"unload": false,
				}
				b.DeleteBuffer(buf, deleteBufferOpts)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("SetCurrentDirectory", func(t *testing.T) {
				wantDir, err := os.UserHomeDir()
				if err != nil {
					t.Fatal(err)
				}

				b := v.NewBatch()
				b.SetCurrentDirectory(wantDir)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var got string
				if err := v.Eval(`getcwd()`, &got); err != nil {
					t.Fatal(err)
				}

				if got != wantDir {
					t.Fatalf("SetCurrentDirectory(%s) = %s, want: %s", wantDir, got, wantDir)
				}
			})

			t.Run("BufferCommands", func(t *testing.T) {
				b := v.NewBatch()

				var commands map[string]*Command
				b.BufferCommands(Buffer(0), make(map[string]interface{}), &commands)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if len(commands) > 0 {
					t.Fatalf("expected commands empty but non-zero: %#v", commands)
				}
			})

			t.Run("BufferOption", func(t *testing.T) {
				b := v.NewBatch()

				var cindent bool
				b.BufferOption(Buffer(0), "cindent", &cindent)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if cindent {
					t.Fatalf("expected cindent is false but got %t", cindent)
				}

				b.SetBufferOption(Buffer(0), "cindent", true)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				b.BufferOption(Buffer(0), "cindent", &cindent)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !cindent {
					t.Fatalf("expected cindent is true but got %t", cindent)
				}

				t.Cleanup(func() {
					// cleanup cindent option
					b.SetBufferOption(Buffer(0), "cindent", false)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("IsBufferLoaded", func(t *testing.T) {
				b := v.NewBatch()

				var loaded bool
				b.IsBufferLoaded(Buffer(0), &loaded)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !loaded {
					t.Fatalf("expected buffer is loaded but got %t", loaded)
				}
			})

			t.Run("IsBufferValid", func(t *testing.T) {
				b := v.NewBatch()

				var valid bool
				b.IsBufferValid(Buffer(0), &valid)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !valid {
					t.Fatalf("expected buffer is valid but got %t", valid)
				}
			})

			t.Run("BufferMark", func(t *testing.T) {
				b := v.NewBatch()

				lines := [][]byte{
					[]byte("a"),
					[]byte("bit of"),
					[]byte("text"),
				}
				b.SetBufferLines(Buffer(0), -1, -1, true, lines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					clearBuffer(t, v, Buffer(0))
				})

				b.SetWindowCursor(Window(0), [2]int{3, 4})

				const (
					mark     = "V"
					wantLine = 3
					wantCol  = 0
				)
				var set bool
				b.SetBufferMark(Buffer(0), mark, wantLine, wantCol, make(map[string]interface{}), &set)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !set {
					t.Fatalf("could not set %s mark", mark)
				}

				var pos [2]int
				b.BufferMark(Buffer(0), mark, &pos)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if pos[0] != wantLine {
					t.Fatalf("got %d extMark line but want %d", pos[0], wantLine)
				}
				if pos[1] != wantCol {
					t.Fatalf("got %d extMark col but want %d", pos[1], wantCol)
				}

				var deleted bool
				b.DeleteBufferMark(Buffer(0), mark, &deleted)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !deleted {
					t.Fatalf("could not delete %s mark", mark)
				}

				var pos2 [2]int
				b.BufferMark(Buffer(0), mark, &pos2)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if pos2[0] != 0 || pos2[1] != 0 {
					t.Fatalf("got %d mark but want zero", pos2)
				}
			})
		})
	}
}

func testWindow(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			wins, err := v.Windows()
			if err != nil {
				t.Fatal(err)
			}
			if len(wins) != 1 {
				for i := 0; i < len(wins); i++ {
					t.Logf("wins[%d]: %v", i, wins[i])
				}
				t.Fatalf("expected one win, found %d wins", len(wins))
			}
			if wins[0] == 0 {
				t.Fatalf("wins[0] == 0")
			}

			win, err := v.CurrentWindow()
			if err != nil {
				t.Fatal(err)
			}
			if win != wins[0] {
				t.Fatalf("win is not wins[0]: win: %v wins[0]: %v", win, wins[0])
			}

			const want = "Window:1000"
			if got := win.String(); got != want {
				t.Fatalf("got %s but want %s", got, want)
			}

			win, err = v.CurrentWindow()
			if err != nil {
				t.Fatal(err)
			}
			if err := v.Command("split"); err != nil {
				t.Fatal(err)
			}
			win2, err := v.CurrentWindow()
			if err != nil {
				t.Fatal(err)
			}

			if err := v.SetCurrentWindow(win); err != nil {
				t.Fatal(err)
			}

			gotwin, err := v.CurrentWindow()
			if err != nil {
				t.Fatal(err)
			}
			if gotwin != win {
				t.Fatalf("expected current window %s but got %s", win, gotwin)
			}

			if err := v.HideWindow(win2); err != nil {
				t.Fatalf("failed to HideWindow(%v)", win2)
			}
			wins2, err := v.Windows()
			if err != nil {
				t.Fatal(err)
			}
			if len(wins2) != 1 {
				for i := 0; i < len(wins2); i++ {
					t.Logf("wins[%d]: %v", i, wins2[i])
				}
				t.Fatalf("expected one win, found %d wins", len(wins2))
			}
			if wins2[0] == 0 {
				t.Fatalf("wins[0] == 0")
			}
			if win != wins2[0] {
				t.Fatalf("win2 is not wins2[0]: want: %v, win2: %v ", wins2[0], win)
			}

			t.Run("WindowBuffer", func(t *testing.T) {
				skipVersion(t, "v0.6.0")

				gotBuf, err := v.WindowBuffer(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				wantBuffer := Buffer(1)
				if gotBuf != wantBuffer {
					t.Fatalf("want %s buffer but got %s", wantBuffer, gotBuf)
				}

				buf, err := v.CreateBuffer(true, true)
				if err != nil {
					t.Fatal(err)
				}

				if err := v.SetBufferToWindow(Window(0), buf); err != nil {
					t.Fatal(err)
				}

				gotBuf2, err := v.WindowBuffer(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				if gotBuf2 != buf {
					t.Fatalf("want %s buffer but got %s", buf, gotBuf2)
				}

				t.Cleanup(func() {
					if err := v.SetBufferToWindow(Window(0), gotBuf); err != nil {
						t.Fatal(err)
					}

					deleteBufferOpts := map[string]bool{
						"force":  true,
						"unload": false,
					}
					if err := v.DeleteBuffer(buf, deleteBufferOpts); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("WindowCursor", func(t *testing.T) {
				wantLine := []byte("hello world")
				if err := v.SetCurrentLine(wantLine); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := v.DeleteCurrentLine(); err != nil {
						t.Fatal(err)
					}
				})

				wantPos := [2]int{1, 5}
				if err := v.SetWindowCursor(Window(0), wantPos); err != nil {
					t.Fatal(err)
				}

				gotPos, err := v.WindowCursor(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				if wantPos != gotPos {
					t.Fatalf("want %#v position buf got %#v", wantPos, gotPos)
				}
			})

			t.Run("WindowVar", func(t *testing.T) {
				wantValue := []int{1, 2}
				if err := v.SetWindowVar(Window(0), "lua", wantValue); err != nil {
					t.Fatal(err)
				}

				var gotValue []int
				if err := v.WindowVar(Window(0), "lua", &gotValue); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(gotValue, wantValue) {
					t.Fatalf("want %#v but got %#v", wantValue, gotValue)
				}

				if err := v.DeleteWindowVar(Window(0), "lua"); err != nil {
					t.Fatal(err)
				}

				if err := v.WindowVar(Window(0), "lua", nil); err == nil {
					t.Fatalf("expect Key not found but fonud key")
				}
			})

			t.Run("WindowOption", func(t *testing.T) {
				wantValue := "+1"
				if err := v.SetWindowOption(Window(0), "colorcolumn", &wantValue); err != nil {
					t.Fatal(err)
				}

				var gotValue string
				if err := v.WindowOption(Window(0), "colorcolumn", &gotValue); err != nil {
					t.Fatal(err)
				}

				if gotValue != wantValue {
					t.Fatalf("expected %s but got %s", wantValue, gotValue)
				}

				t.Cleanup(func() {
					if err := v.SetWindowOption(Window(0), "colorcolumn", ""); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("WindowPosition", func(t *testing.T) {
				gotPos, err := v.WindowPosition(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				wantPos := [2]int{0, 0}
				if gotPos != wantPos {
					t.Fatalf("expected %v but got %v", wantPos, gotPos)
				}
			})

			t.Run("WindowTabpage", func(t *testing.T) {
				gotTabpage, err := v.WindowTabpage(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				wantTabpage := Tabpage(1)
				if gotTabpage != wantTabpage {
					t.Fatalf("expected %v but got %v", wantTabpage, gotTabpage)
				}
			})

			t.Run("WindowNumber", func(t *testing.T) {
				gotWinNum, err := v.WindowNumber(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				wantWinNum := 1
				if gotWinNum != wantWinNum {
					t.Fatalf("expected %v but got %v", wantWinNum, gotWinNum)
				}
			})

			t.Run("IsWindowValid", func(t *testing.T) {
				valid, err := v.IsWindowValid(Window(0))
				if err != nil {
					t.Fatal(err)
				}

				if !valid {
					t.Fatalf("expected valid but got %t", valid)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var wins []Window
			b.Windows(&wins)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if len(wins) != 1 {
				t.Fatalf("expected one win, found %d wins", len(wins))
			}
			if wins[0] == 0 {
				t.Fatalf("wins[0] == 0")
			}

			var win Window
			b.CurrentWindow(&win)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if win != wins[0] {
				t.Fatalf("win is not wins[0]: win: %v wins[0]: %v", win, wins[0])
			}

			const want = "Window:1000"
			if got := win.String(); got != want {
				t.Fatalf("got %s but want %s", got, want)
			}

			b.CurrentWindow(&win)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			b.Command("split")
			var win2 Window
			b.CurrentWindow(&win2)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			b.SetCurrentWindow(win)

			var gotwin Window
			b.CurrentWindow(&gotwin)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if gotwin != win {
				t.Fatalf("expected current window %s but got %s", win, gotwin)
			}

			b.HideWindow(win2)
			var wins2 []Window
			b.Windows(&wins2)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if len(wins2) != 1 {
				for i := 0; i < len(wins2); i++ {
					t.Logf("wins[%d]: %v", i, wins2[i])
				}
				t.Fatalf("expected one win, found %d wins", len(wins2))
			}
			if wins2[0] == 0 {
				t.Fatalf("wins[0] == 0")
			}
			if win != wins2[0] {
				t.Fatalf("win2 is not wins2[0]: want: %v, win2: %v ", wins2[0], win)
			}

			t.Run("WindowBuffer", func(t *testing.T) {
				skipVersion(t, "v0.6.0")

				b := v.NewBatch()

				var gotBuf Buffer
				b.WindowBuffer(Window(0), &gotBuf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantBuffer := Buffer(1)
				if gotBuf != wantBuffer {
					t.Fatalf("want %s buffer but got %s", wantBuffer, gotBuf)
				}

				var buf Buffer
				b.CreateBuffer(true, true, &buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				b.SetBufferToWindow(Window(0), buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var gotBuf2 Buffer
				b.WindowBuffer(Window(0), &gotBuf2)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if gotBuf2 != buf {
					t.Fatalf("want %s buffer but got %s", buf, gotBuf2)
				}

				t.Cleanup(func() {
					b.SetBufferToWindow(Window(0), gotBuf)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}

					deleteBufferOpts := map[string]bool{
						"force":  true,
						"unload": false,
					}
					b.DeleteBuffer(buf, deleteBufferOpts)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("WindowCursor", func(t *testing.T) {
				b := v.NewBatch()

				wantLine := []byte("hello world")
				b.SetCurrentLine(wantLine)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					b.DeleteCurrentLine()
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
				})

				wantPos := [2]int{1, 5}
				b.SetWindowCursor(Window(0), wantPos)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var gotPos [2]int
				b.WindowCursor(Window(0), &gotPos)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if wantPos != gotPos {
					t.Fatalf("want %#v position buf got %#v", wantPos, gotPos)
				}
			})

			t.Run("WindowVar", func(t *testing.T) {
				b := v.NewBatch()

				wantValue := []int{1, 2}
				b.SetWindowVar(Window(0), "lua", wantValue)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var gotValue []int
				b.WindowVar(Window(0), "lua", &gotValue)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(gotValue, wantValue) {
					t.Fatalf("want %#v but got %#v", wantValue, gotValue)
				}

				b.DeleteWindowVar(Window(0), "lua")
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				b.WindowVar(Window(0), "lua", nil)
				if err := b.Execute(); err == nil {
					t.Fatalf("expect Key not found but fonud key")
				}
			})

			t.Run("WindowOption", func(t *testing.T) {
				b := v.NewBatch()

				wantValue := "+1"
				b.SetWindowOption(Window(0), "colorcolumn", &wantValue)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var gotValue string
				b.WindowOption(Window(0), "colorcolumn", &gotValue)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if gotValue != wantValue {
					t.Fatalf("expected %s but got %s", wantValue, gotValue)
				}

				t.Cleanup(func() {
					b.SetWindowOption(Window(0), "colorcolumn", "")
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
				})
			})

			t.Run("WindowPosition", func(t *testing.T) {
				b := v.NewBatch()

				var gotPos [2]int
				b.WindowPosition(Window(0), &gotPos)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantPos := [2]int{0, 0}
				if gotPos != wantPos {
					t.Fatalf("expected %v but got %v", wantPos, gotPos)
				}
			})

			t.Run("WindowTabpage", func(t *testing.T) {
				b := v.NewBatch()

				var gotTabpage Tabpage
				b.WindowTabpage(Window(0), &gotTabpage)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantTabpage := Tabpage(1)
				if gotTabpage != wantTabpage {
					t.Fatalf("expected %v but got %v", wantTabpage, gotTabpage)
				}
			})

			t.Run("WindowNumber", func(t *testing.T) {
				b := v.NewBatch()

				var gotWinNum int
				b.WindowNumber(Window(0), &gotWinNum)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantWinNum := 1
				if gotWinNum != wantWinNum {
					t.Fatalf("expected %v but got %v", wantWinNum, gotWinNum)
				}
			})

			t.Run("IsWindowValid", func(t *testing.T) {
				b := v.NewBatch()

				var valid bool
				b.IsWindowValid(Window(0), &valid)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !valid {
					t.Fatalf("expected valid but got %t", valid)
				}
			})
		})
	}
}

func testTabpage(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			pages, err := v.Tabpages()
			if err != nil {
				t.Fatal(err)
			}
			if len(pages) != 1 {
				t.Fatalf("expected one page, found %d pages", len(pages))
			}
			if pages[0] == 0 {
				t.Fatalf("pages[0] is not 0: %d", pages[0])
			}

			page, err := v.CurrentTabpage()
			if err != nil {
				t.Fatal(err)
			}
			if page != pages[0] {
				t.Fatalf("page is not pages[0]: page: %v pages[0]: %v", page, pages[0])
			}

			const want = "Tabpage:1"
			if got := page.String(); got != want {
				t.Fatalf("got %s but want %s", got, want)
			}

			if err := v.SetCurrentTabpage(page); err != nil {
				t.Fatal(err)
			}

			t.Run("TabpageWindow", func(t *testing.T) {
				gotWin, err := v.TabpageWindow(Tabpage(0))
				if err != nil {
					t.Fatal(err)
				}

				wantWin := Window(1000)
				if !reflect.DeepEqual(gotWin, wantWin) {
					t.Fatalf("expected %v but got %v", wantWin, gotWin)
				}
			})

			t.Run("TabpageWindows", func(t *testing.T) {
				gotWins, err := v.TabpageWindows(Tabpage(0))
				if err != nil {
					t.Fatal(err)
				}

				wantWins := []Window{Window(1000)}
				if !reflect.DeepEqual(gotWins, wantWins) {
					t.Fatalf("expected %v but got %v", wantWins, gotWins)
				}
			})

			t.Run("TabpageNumber", func(t *testing.T) {
				gotTabpageNum, err := v.TabpageNumber(Tabpage(0))
				if err != nil {
					t.Fatal(err)
				}

				wantTabpageNum := 1
				if gotTabpageNum != wantTabpageNum {
					t.Fatalf("expected %v but got %v", wantTabpageNum, gotTabpageNum)
				}
			})

			t.Run("IsTabpageValid", func(t *testing.T) {
				valid, err := v.IsTabpageValid(Tabpage(0))
				if err != nil {
					t.Fatal(err)
				}

				if !valid {
					t.Fatalf("expected valid but got %t", valid)
				}
			})

			t.Run("TabpageVar", func(t *testing.T) {
				wantValue := []int{1, 2}
				if err := v.SetTabpageVar(Tabpage(0), "lua", wantValue); err != nil {
					t.Fatal(err)
				}

				var gotValue []int
				if err := v.TabpageVar(Tabpage(0), "lua", &gotValue); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(gotValue, wantValue) {
					t.Fatalf("want %#v but got %#v", wantValue, gotValue)
				}

				if err := v.DeleteTabpageVar(Tabpage(0), "lua"); err != nil {
					t.Fatal(err)
				}

				if err := v.TabpageVar(Tabpage(0), "lua", nil); err == nil {
					t.Fatalf("expect Key not found but fonud key")
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var pages []Tabpage
			b.Tabpages(&pages)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if len(pages) != 1 {
				t.Fatalf("expected one page, found %d pages", len(pages))
			}
			if pages[0] == 0 {
				t.Fatalf("pages[0] is not 0: %d", pages[0])
			}

			var page Tabpage
			b.CurrentTabpage(&page)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if page != pages[0] {
				t.Fatalf("page is not pages[0]: page: %v pages[0]: %v", page, pages[0])
			}

			const want = "Tabpage:1"
			if got := page.String(); got != want {
				t.Fatalf("got %s but want %s", got, want)
			}

			b.SetCurrentTabpage(page)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			t.Run("TabpageWindow", func(t *testing.T) {
				b := v.NewBatch()

				var gotWin Window
				b.TabpageWindow(Tabpage(0), &gotWin)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantWin := Window(1000)
				if gotWin != wantWin {
					t.Fatalf("expected %v but got %v", wantWin, gotWin)
				}
			})

			t.Run("TabpageWindows", func(t *testing.T) {
				b := v.NewBatch()

				var gotWins []Window
				b.TabpageWindows(Tabpage(0), &gotWins)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantWins := []Window{Window(1000)}
				if !reflect.DeepEqual(gotWins, wantWins) {
					t.Fatalf("expected %v but got %v", wantWins, gotWins)
				}
			})

			t.Run("TabpageNumber", func(t *testing.T) {
				b := v.NewBatch()

				var gotTabpageNum int
				b.TabpageNumber(Tabpage(0), &gotTabpageNum)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantTabpageNum := 1
				if gotTabpageNum != wantTabpageNum {
					t.Fatalf("expected %v but got %v", wantTabpageNum, gotTabpageNum)
				}
			})

			t.Run("IsWindowValid", func(t *testing.T) {
				b := v.NewBatch()

				var valid bool
				b.IsTabpageValid(Tabpage(0), &valid)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !valid {
					t.Fatalf("expected valid but got %t", valid)
				}
			})

			t.Run("TabpageVar", func(t *testing.T) {
				b := v.NewBatch()

				wantValue := []int{1, 2}
				b.SetTabpageVar(Tabpage(0), "lua", wantValue)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var gotValue []int
				b.TabpageVar(Tabpage(0), "lua", &gotValue)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(gotValue, wantValue) {
					t.Fatalf("want %#v but got %#v", wantValue, gotValue)
				}

				b.DeleteTabpageVar(Tabpage(0), "lua")
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				b.TabpageVar(Tabpage(0), "lua", nil)
				if err := b.Execute(); err == nil {
					t.Fatalf("expect Key not found but fonud key")
				}
			})
		})
	}
}

func testLines(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Run("CurrentLine", func(t *testing.T) {
				clearBuffer(t, v, Buffer(0))

				beforeLine, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}

				wantLine := []byte("hello world")
				if err := v.SetCurrentLine(wantLine); err != nil {
					t.Fatal(err)
				}

				afterLine, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}
				if bytes.EqualFold(beforeLine, afterLine) {
					t.Fatalf("current line not change: before: %v, after: %v", beforeLine, afterLine)
				}

				if err := v.DeleteCurrentLine(); err != nil {
					t.Fatal(err)
				}
				deletedLine, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}
				if len(deletedLine) != 0 {
					t.Fatal("DeleteCurrentLine not deleted")
				}
			})

			t.Run("BufferLines", func(t *testing.T) {
				buf, err := v.CurrentBuffer()
				if err != nil {
					t.Fatal(err)
				}
				defer clearBuffer(t, v, buf) // clear buffer after run sub-test.

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

				const wantCount = 2
				count, err := v.BufferLineCount(buf)
				if err != nil {
					t.Fatal(err)
				}
				if count != wantCount {
					t.Fatalf("got count %d but want %d", count, wantCount)
				}

				const wantOffset = 12 // [][]byte{[]byte("hello"), []byte("\n"), []byte("world"), []byte("\n")}
				offset, err := v.BufferOffset(buf, count)
				if err != nil {
					t.Fatal(err)
				}
				if offset != wantOffset {
					t.Fatalf("got offset %d but want %d", offset, wantOffset)
				}
			})

			t.Run("SetBufferText", func(t *testing.T) {
				buf, err := v.CurrentBuffer()
				if err != nil {
					t.Fatal(err)
				}
				defer clearBuffer(t, v, buf) // clear buffer after run sub-test.

				// sets test buffer text.
				lines := [][]byte{[]byte("Vim is the"), []byte("Nvim-fork? focused on extensibility and usability")}
				if err := v.SetBufferLines(buf, 0, -1, true, lines); err != nil {
					t.Fatal(err)
				}

				// Replace `Vim is the` to `Neovim is the`
				if err := v.SetBufferText(buf, 0, 0, 0, 3, [][]byte{[]byte("Neovim")}); err != nil {
					t.Fatal(err)
				}
				// Replace `Nvim-fork?` to `Vim-fork`
				if err := v.SetBufferText(buf, 1, 0, 1, 10, [][]byte{[]byte("Vim-fork")}); err != nil {
					t.Fatal(err)
				}

				want := [2][]byte{
					[]byte("Neovim is the"),
					[]byte("Vim-fork focused on extensibility and usability"),
				}
				got, err := v.BufferLines(buf, 0, -1, true)
				if err != nil {
					t.Fatal(err)
				}

				// assert buffer lines count.
				const wantCount = 2
				if len(got) != wantCount {
					t.Fatalf("expected buffer lines rows is %d: got %d", wantCount, len(got))
				}

				// assert row 1 buffer text.
				if !bytes.EqualFold(want[0], got[0]) {
					t.Fatalf("row 1 is not equal: want: %q, got: %q", string(want[0]), string(got[0]))
				}

				// assert row 2 buffer text.
				if !bytes.EqualFold(want[1], got[1]) {
					t.Fatalf("row 2 is not equal: want: %q, got: %q", string(want[1]), string(got[1]))
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("CurrentLine", func(t *testing.T) {
				clearBuffer(t, v, Buffer(0))

				b := v.NewBatch()

				var beforeLine []byte
				b.CurrentLine(&beforeLine)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantLine := []byte("hello world")
				b.SetCurrentLine(wantLine)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var afterLine []byte
				b.CurrentLine(&afterLine)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if bytes.EqualFold(beforeLine, afterLine) {
					t.Fatalf("current line not change: before: %v, after: %v", beforeLine, afterLine)
				}

				b.DeleteCurrentLine()
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				var deletedLine []byte
				b.CurrentLine(&deletedLine)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if len(deletedLine) != 0 {
					t.Fatal("DeleteCurrentLine not deleted")
				}
			})

			t.Run("BufferLines", func(t *testing.T) {
				b := v.NewBatch()

				var buf Buffer
				b.CurrentBuffer(&buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				defer clearBuffer(t, v, buf) // clear buffer after run sub-test.

				lines := [][]byte{[]byte("hello"), []byte("world")}
				b.SetBufferLines(buf, 0, -1, true, lines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var lines2 [][]byte
				b.BufferLines(buf, 0, -1, true, &lines2)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(lines2, lines) {
					t.Fatalf("lines = %+v, want %+v", lines2, lines)
				}

				const wantCount = 2
				var count int
				b.BufferLineCount(buf, &count)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if count != wantCount {
					t.Fatalf("count is not 2 %d", count)
				}

				const wantOffset = 12 // [][]byte{[]byte("hello"), []byte("\n"), []byte("world"), []byte("\n")}
				var offset int
				b.BufferOffset(buf, count, &offset)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if offset != wantOffset {
					t.Fatalf("got offset %d but want %d", offset, wantOffset)
				}
			})

			t.Run("SetBufferText", func(t *testing.T) {
				b := v.NewBatch()

				var buf Buffer
				b.CurrentBuffer(&buf)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				defer clearBuffer(t, v, buf) // clear buffer after run sub-test.

				// sets test buffer text.
				lines := [][]byte{[]byte("Vim is the"), []byte("Nvim-fork? focused on extensibility and usability")}
				b.SetBufferLines(buf, 0, -1, true, lines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// Replace `Vim is the` to `Neovim is the`
				b.SetBufferText(buf, 0, 0, 0, 3, [][]byte{[]byte("Neovim")})
				// Replace `Nvim-fork?` to `Vim-fork`
				b.SetBufferText(buf, 1, 0, 1, 10, [][]byte{[]byte("Vim-fork")})
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				want := [2][]byte{
					[]byte("Neovim is the"),
					[]byte("Vim-fork focused on extensibility and usability"),
				}
				var got [][]byte
				b.BufferLines(buf, 0, -1, true, &got)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// assert buffer lines count.
				const wantCount = 2
				if len(got) != wantCount {
					t.Fatalf("expected buffer lines rows is %d: got %d", wantCount, len(got))
				}

				// assert row 1 buffer text.
				if !bytes.EqualFold(want[0], got[0]) {
					t.Fatalf("row 1 is not equal: want: %q, got: %q", string(want[0]), string(got[0]))
				}

				// assert row 1 buffer text.
				if !bytes.EqualFold(want[1], got[1]) {
					t.Fatalf("row 2 is not equal: want: %q, got: %q", string(want[1]), string(got[1]))
				}
			})
		})
	}
}

func testCommand(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Commands", func(t *testing.T) {
			t.Run("Nvim", func(t *testing.T) {
				opts := map[string]interface{}{
					"builtin": false,
				}
				cmds, err := v.Commands(opts)
				if err != nil {
					t.Fatal(err)
				}
				if len(cmds) > 0 {
					t.Fatalf("expected 0 length but got %#v", cmds)
				}
			})

			t.Run("Batch", func(t *testing.T) {
				b := v.NewBatch()

				opts := map[string]interface{}{
					"builtin": false,
				}
				var cmds map[string]*Command
				b.Commands(opts, &cmds)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if len(cmds) > 0 {
					t.Fatalf("expected 0 length but got %#v", cmds)
				}
			})
		})

		t.Run("UserCommand", func(t *testing.T) {
			tests := map[string]struct {
				name    string
				command UserCommand
				opts    map[string]interface{}
				want    string
			}{
				"SayHello": {
					name:    "SayHello",
					command: UserVimCommand(`echo "Hello world!"`),
					opts: map[string]interface{}{
						"force": false,
					},
					want: "Hello world!",
				},
			}
			for name, tt := range tests {
				t.Run(path.Join(name, "Nvim"), func(t *testing.T) {
					skipVersion(t, "v0.7.0")

					if err := v.AddUserCommand(tt.name, tt.command, tt.opts); err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						if err := v.DeleteUserCommand(tt.name); err != nil {
							t.Fatal(err)
						}
					})

					got, err := v.Exec(tt.name, true)
					if err != nil {
						t.Fatal(err)
					}
					if !strings.EqualFold(tt.want, got) {
						t.Fatalf("expected %s but got %s", tt.want, got)
					}
				})

				t.Run(path.Join(name, "Batch"), func(t *testing.T) {
					skipVersion(t, "v0.7.0")

					b := v.NewBatch()

					b.AddUserCommand(tt.name, tt.command, tt.opts)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						b.DeleteUserCommand(tt.name)
						if err := b.Execute(); err != nil {
							t.Fatal(err)
						}
					})

					var got string
					b.Exec(tt.name, true, &got)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					if !strings.EqualFold(tt.want, got) {
						t.Fatalf("expected %s but got %s", tt.want, got)
					}
				})
			}
		})

		t.Run("BufferUserCommand", func(t *testing.T) {
			tests := map[string]struct {
				name    string
				command UserCommand
				opts    map[string]interface{}
				want    string
			}{
				"SayHello": {
					name:    "SayHello",
					command: UserVimCommand(`echo "Hello world!"`),
					opts: map[string]interface{}{
						"force": false,
					},
					want: "Hello world!",
				},
			}
			for name, tt := range tests {
				t.Run(path.Join(name, "Nvim"), func(t *testing.T) {
					skipVersion(t, "v0.7.0")

					if err := v.AddBufferUserCommand(Buffer(0), tt.name, tt.command, tt.opts); err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						if err := v.DeleteBufferUserCommand(Buffer(0), tt.name); err != nil {
							t.Fatal(err)
						}
					})

					got, err := v.Exec(tt.name, true)
					if err != nil {
						t.Fatal(err)
					}
					if !strings.EqualFold(tt.want, got) {
						t.Fatalf("expected %s but got %s", tt.want, got)
					}
				})

				t.Run(path.Join(name, "Batch"), func(t *testing.T) {
					skipVersion(t, "v0.7.0")

					b := v.NewBatch()

					b.AddBufferUserCommand(Buffer(0), tt.name, tt.command, tt.opts)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() {
						b.DeleteBufferUserCommand(Buffer(0), tt.name)
						if err := b.Execute(); err != nil {
							t.Fatal(err)
						}
					})

					var got string
					b.Exec(tt.name, true, &got)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					if !strings.EqualFold(tt.want, got) {
						t.Fatalf("expected %s but got %s", tt.want, got)
					}
				})
			}
		})
	}
}

func testVar(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		const (
			varKey = `gvar`
			varVal = `gval`
		)

		const (
			vvarKey  = "statusmsg"
			wantVvar = "test"
		)

		t.Run("Nvim", func(t *testing.T) {
			t.Run("Var", func(t *testing.T) {
				if err := v.SetVar(varKey, varVal); err != nil {
					t.Fatal(err)
				}

				var value interface{}
				if err := v.Var(varKey, &value); err != nil {
					t.Fatal(err)
				}
				if value != varVal {
					t.Fatalf("got %v, want %q", value, varVal)
				}

				if err := v.SetVar(varKey, ""); err != nil {
					t.Fatal(err)
				}

				value = nil
				if err := v.Var(varKey, &value); err != nil {
					t.Fatal(err)
				}
				if value != "" {
					t.Fatalf("got %v, want %q", value, "")
				}

				if err := v.DeleteVar(varKey); err != nil && !strings.Contains(err.Error(), "Key not found") {
					t.Fatal(err)
				}
			})

			t.Run("VVar", func(t *testing.T) {
				if err := v.SetVVar(vvarKey, wantVvar); err != nil {
					t.Fatalf("failed to SetVVar: %v", err)
				}

				var vvar string
				if err := v.VVar(vvarKey, &vvar); err != nil {
					t.Fatalf("failed to SetVVar: %v", err)
				}
				if vvar != wantVvar {
					t.Fatalf("VVar(%s, %s) = %s, want: %s", vvarKey, wantVvar, vvar, wantVvar)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("Var", func(t *testing.T) {
				b := v.NewBatch()

				b.SetVar(varKey, varVal)

				var value interface{}
				b.Var(varKey, &value)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if value != varVal {
					t.Fatalf("got %v, want %q", value, varVal)
				}

				b.SetVar(varKey, "")

				value = nil
				b.Var(varKey, &value)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if value != "" {
					t.Fatalf("got %v, want %q", value, "")
				}

				b.DeleteVar(varKey)
				if err := b.Execute(); err != nil && !strings.Contains(err.Error(), "Key not found") {
					t.Fatal(err)
				}
			})

			t.Run("VVar", func(t *testing.T) {
				b := v.NewBatch()

				b.SetVVar(vvarKey, wantVvar)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var vvar string
				b.VVar(vvarKey, &vvar)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if vvar != wantVvar {
					t.Fatalf("VVar(%s, %s) = %s, want: %s", vvarKey, wantVvar, vvar, wantVvar)
				}
			})
		})
	}
}

func testMessage(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Run("Echo", func(t *testing.T) {
				const wantEcho = `hello Echo`
				chunk := []TextChunk{
					{
						Text: wantEcho,
					},
				}
				if err := v.Echo(chunk, true, make(map[string]interface{})); err != nil {
					t.Fatalf("failed to Echo: %v", err)
				}

				gotEcho, err := v.Exec("message", true)
				if err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotEcho != wantEcho {
					t.Fatalf("Echo(%q) = %q, want: %q", wantEcho, gotEcho, wantEcho)
				}
			})

			t.Run("WriteOut", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				const wantWriteOut = `hello WriteOut`
				if err := v.WriteOut(wantWriteOut + "\n"); err != nil {
					t.Fatalf("failed to WriteOut: %v", err)
				}

				var gotWriteOut string
				if err := v.VVar("statusmsg", &gotWriteOut); err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotWriteOut != wantWriteOut {
					t.Fatalf("WriteOut(%q) = %q, want: %q", wantWriteOut, gotWriteOut, wantWriteOut)
				}
			})

			t.Run("WriteErr", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				const wantWriteErr = `hello WriteErr`
				if err := v.WriteErr(wantWriteErr + "\n"); err != nil {
					t.Fatalf("failed to WriteErr: %v", err)
				}

				var gotWriteErr string
				if err := v.VVar("errmsg", &gotWriteErr); err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotWriteErr != wantWriteErr {
					t.Fatalf("WriteErr(%q) = %q, want: %q", wantWriteErr, gotWriteErr, wantWriteErr)
				}
			})

			t.Run("WritelnErr", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				const wantWritelnErr = `hello WritelnErr`
				if err := v.WritelnErr(wantWritelnErr); err != nil {
					t.Fatalf("failed to WriteErr: %v", err)
				}

				var gotWritelnErr string
				if err := v.VVar("errmsg", &gotWritelnErr); err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotWritelnErr != wantWritelnErr {
					t.Fatalf("WritelnErr(%q) = %q, want: %q", wantWritelnErr, gotWritelnErr, wantWritelnErr)
				}
			})

			t.Run("Notify", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				const wantNotifyMsg = `hello Notify`
				if err := v.Notify(wantNotifyMsg, LogInfoLevel, make(map[string]interface{})); err != nil {
					t.Fatalf("failed to Notify: %v", err)
				}
				gotNotifyMsg, err := v.Exec(":messages", true)
				if err != nil {
					t.Fatalf("failed to messages command: %v", err)
				}
				if wantNotifyMsg != gotNotifyMsg {
					t.Fatalf("Notify(%[1]q, %[2]q) = %[3]q, want: %[1]q", wantNotifyMsg, LogInfoLevel, gotNotifyMsg)
				}
			})

			t.Run("Notify/Error", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				const wantNotifyErr = `hello Notify Error`
				if err := v.Notify(wantNotifyErr, LogErrorLevel, make(map[string]interface{})); err != nil {
					t.Fatalf("failed to Notify: %v", err)
				}
				var gotNotifyErr string
				if err := v.VVar("errmsg", &gotNotifyErr); err != nil {
					t.Fatalf("could not get v:errmsg nvim variable: %v", err)
				}
				if wantNotifyErr != gotNotifyErr {
					t.Fatalf("Notify(%[1]q, %[2]q) = %[3]q, want: %[1]q", wantNotifyErr, LogErrorLevel, gotNotifyErr)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("Echo", func(t *testing.T) {
				b := v.NewBatch()

				const wantEcho = `hello Echo`
				chunk := []TextChunk{
					{
						Text: wantEcho,
					},
				}
				b.Echo(chunk, true, make(map[string]interface{}))

				var gotEcho string
				b.Exec("message", true, &gotEcho)
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to Execute: %v", err)
				}
				if gotEcho != wantEcho {
					t.Fatalf("Echo(%q) = %q, want: %q", wantEcho, gotEcho, wantEcho)
				}
			})
			t.Run("WriteOut", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				b := v.NewBatch()

				const wantWriteOut = `hello WriteOut`
				b.WriteOut(wantWriteOut + "\n")
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to WriteOut: %v", err)
				}

				var gotWriteOut string
				b.VVar("statusmsg", &gotWriteOut)
				if err := b.Execute(); err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotWriteOut != wantWriteOut {
					t.Fatalf("b.WriteOut(%q) = %q, want: %q", wantWriteOut, gotWriteOut, wantWriteOut)
				}
			})

			t.Run("WriteErr", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				b := v.NewBatch()

				const wantWriteErr = `hello WriteErr`
				b.WriteErr(wantWriteErr + "\n")
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to WriteErr: %v", err)
				}
				var gotWriteErr string
				b.VVar("errmsg", &gotWriteErr)
				if err := b.Execute(); err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotWriteErr != wantWriteErr {
					t.Fatalf("b.WriteErr(%q) = %q, want: %q", wantWriteErr, gotWriteErr, wantWriteErr)
				}
			})

			t.Run("WritelnErr", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				b := v.NewBatch()

				const wantWritelnErr = `hello WritelnErr`
				b.WritelnErr(wantWritelnErr)
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to WriteErr: %v", err)
				}
				var gotWritelnErr string
				b.VVar("errmsg", &gotWritelnErr)
				if err := b.Execute(); err != nil {
					t.Fatalf("could not get v:statusmsg nvim variable: %v", err)
				}
				if gotWritelnErr != wantWritelnErr {
					t.Fatalf("b.WritelnErr(%q) = %q, want: %q", wantWritelnErr, gotWritelnErr, wantWritelnErr)
				}
			})

			t.Run("Notify", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				b := v.NewBatch()

				const wantNotifyMsg = `hello Notify`
				b.Notify(wantNotifyMsg, LogInfoLevel, make(map[string]interface{}))
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to Notify: %v", err)
				}
				var gotNotifyMsg string
				b.Exec(":messages", true, &gotNotifyMsg)
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to \":messages\" command: %v", err)
				}
				if wantNotifyMsg != gotNotifyMsg {
					t.Fatalf("Notify(%[1]q, %[2]q) = %[3]q, want: %[1]q", wantNotifyMsg, LogInfoLevel, gotNotifyMsg)
				}
			})

			t.Run("Notify/Error", func(t *testing.T) {
				defer func() {
					// cleanup v:statusmsg
					if err := v.SetVVar("statusmsg", ""); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
					// clear messages
					if _, err := v.Exec(":messages clear", false); err != nil {
						t.Fatalf("failed to SetVVar: %v", err)
					}
				}()

				b := v.NewBatch()

				const wantNotifyErr = `hello Notify Error`
				b.Notify(wantNotifyErr, LogErrorLevel, make(map[string]interface{}))
				if err := b.Execute(); err != nil {
					t.Fatalf("failed to Notify: %v", err)
				}
				var gotNotifyErr string
				b.VVar("errmsg", &gotNotifyErr)
				if err := b.Execute(); err != nil {
					t.Fatalf("could not get v:errmsg nvim variable: %v", err)
				}
				if wantNotifyErr != gotNotifyErr {
					t.Fatalf("Notify(%[1]q, %[2]q) = %[3]q, want: %[1]q", wantNotifyErr, LogErrorLevel, gotNotifyErr)
				}
			})
		})
	}
}

func testKey(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Run("FeedKeys", func(t *testing.T) {
				// cleanup current Buffer after tests.
				defer clearBuffer(t, v, Buffer(0))

				const (
					keys      = `iabc<ESC>`
					mode      = `n`
					escapeCSI = false
				)
				input, err := v.ReplaceTermcodes(keys, true, true, true)
				if err != nil {
					t.Fatal(err)
				}

				// clear current Buffer before run FeedKeys.
				clearBuffer(t, v, Buffer(0))

				if err := v.FeedKeys(input, mode, escapeCSI); err != nil {
					t.Fatal(err)
				}

				wantLines := []byte{'a', 'b', 'c'}

				gotLines, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(gotLines, wantLines) {
					t.Fatalf("FeedKeys(%s, %s, %t): got %v, want %v", input, mode, escapeCSI, gotLines, wantLines)
				}
			})

			t.Run("Input", func(t *testing.T) {
				// cleanup current Buffer after tests.
				defer clearBuffer(t, v, Buffer(0))

				const (
					keys      = `iabc<ESC>`
					mode      = `n`
					escapeCSI = false
				)
				input, err := v.ReplaceTermcodes(keys, true, true, true)
				if err != nil {
					t.Fatal(err)
				}

				// clear current Buffer before run FeedKeys.
				clearBuffer(t, v, Buffer(0))

				written, err := v.Input(input)
				if err != nil {
					t.Fatal(err)
				}
				if written != len(input) {
					t.Fatalf("Input(%s) = %d: want: %d", input, written, len(input))
				}

				wantLines := []byte{'a', 'b', 'c'}
				gotLines, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(gotLines, wantLines) {
					t.Fatalf("FeedKeys(%s, %s, %t): got %v, want %v", input, mode, escapeCSI, gotLines, wantLines)
				}
			})

			t.Run("InputMouse", func(t *testing.T) {
				defer func() {
					// cleanup current Buffer after tests.
					clearBuffer(t, v, Buffer(0))

					input, err := v.ReplaceTermcodes(`<ESC>`, true, true, true)
					if err != nil {
						t.Fatal(err)
					}
					if err := v.FeedKeys(input, `n`, true); err != nil {
						t.Fatal(err)
					}
				}()

				// clear current Buffer before run FeedKeys.
				clearBuffer(t, v, Buffer(0))

				lines := [][]byte{
					[]byte("foo bar baz"),
					[]byte("qux quux quuz"),
					[]byte("corge grault garply"),
					[]byte("waldo fred plugh"),
					[]byte("xyzzy thud"),
				}
				if err := v.SetBufferLines(Buffer(0), 0, -1, true, lines); err != nil {
					t.Fatal(err)
				}

				const (
					button       = `left`
					firestAction = `press`
					secondAction = `release`
					modifier     = ""
				)
				const (
					wantGrid = 20
					wantRow  = 2
					wantCol  = 5
				)
				if err := v.InputMouse(button, firestAction, modifier, wantGrid, wantRow, wantCol); err != nil {
					t.Fatal(err)
				}

				// TODO(zchee): assertion
			})

			t.Run("StringWidth", func(t *testing.T) {
				const str = "hello\t"
				got, err := v.StringWidth(str)
				if err != nil {
					t.Fatal(err)
				}
				if got != len(str) {
					t.Fatalf("StringWidth(%s) = %d, want: %d", str, got, len(str))
				}
			})

			t.Run("KeyMap", func(t *testing.T) {
				mode := "n"
				if err := v.SetKeyMap(mode, "y", "yy", make(map[string]bool)); err != nil {
					t.Fatal(err)
				}

				wantMaps := []*Mapping{
					{
						LHS:     "y",
						RHS:     "yy",
						Silent:  0,
						NoRemap: 0,
						Expr:    0,
						Buffer:  0,
						SID:     0,
						NoWait:  0,
					},
				}
				wantMapsLen := 0
				if nvimVersion.Minor >= 6 {
					lastMap := wantMaps[0]
					wantMaps = []*Mapping{
						{
							LHS:     "<C-L>",
							RHS:     "<Cmd>nohlsearch|diffupdate<CR><C-L>",
							Silent:  0,
							NoRemap: 1,
							Expr:    0,
							Buffer:  0,
							SID:     0,
							NoWait:  0,
						},
						{
							LHS:     "Y",
							RHS:     "y$",
							Silent:  0,
							NoRemap: 1,
							Expr:    0,
							Buffer:  0,
							SID:     0,
							NoWait:  0,
						},
						lastMap,
					}
					wantMapsLen = 2
				}
				got, err := v.KeyMap(mode)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, wantMaps) {
					for i, gotmap := range got {
						t.Logf(" got[%d]: %#v", i, gotmap)
					}
					for i, wantmap := range wantMaps {
						t.Logf("want[%d]: %#v", i, wantmap)
					}
					t.Fatalf("KeyMap(%s) = %#v, want: %#v", mode, got, wantMaps)
				}

				if err := v.DeleteKeyMap(mode, "y"); err != nil {
					t.Fatal(err)
				}

				got2, err := v.KeyMap(mode)
				if err != nil {
					t.Fatal(err)
				}
				if len(got2) != wantMapsLen {
					t.Fatalf("expected %d but got %#v", wantMapsLen, got2)
				}
			})

			t.Run("BufferKeyMap", func(t *testing.T) {
				mode := "n"
				buffer := Buffer(0)
				if err := v.SetBufferKeyMap(buffer, mode, "x", "xx", make(map[string]bool)); err != nil {
					t.Fatal(err)
				}

				wantMap := []*Mapping{
					{
						LHS:     "x",
						RHS:     "xx",
						Silent:  0,
						NoRemap: 0,
						Expr:    0,
						Buffer:  1,
						SID:     0,
						NoWait:  0,
					},
				}
				got, err := v.BufferKeyMap(buffer, mode)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(wantMap, got) {
					t.Fatalf("KeyMap(n) = %#v, want: %#v", got[0], wantMap[0])
				}

				if err := v.DeleteBufferKeyMap(buffer, mode, "x"); err != nil {
					t.Fatal(err)
				}

				got2, err := v.BufferKeyMap(buffer, mode)
				if err != nil {
					t.Fatal(err)
				}
				if wantLen := 0; len(got2) != wantLen {
					t.Fatalf("expected %d but got %#v", wantLen, got2)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("FeedKeys", func(t *testing.T) {
				// cleanup current Buffer after tests.
				defer clearBuffer(t, v, Buffer(0))

				b := v.NewBatch()

				const (
					keys      = `iabc<ESC>`
					mode      = `n`
					escapeCSI = false
				)
				var input string
				b.ReplaceTermcodes(keys, true, true, true, &input)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// clear current Buffer before run FeedKeys.
				clearBuffer(t, v, Buffer(0))

				b.FeedKeys(input, mode, escapeCSI)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantLines := []byte{'a', 'b', 'c'}
				gotLines, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(gotLines, wantLines) {
					t.Fatalf("FeedKeys(%s, %s, %t): got %v, want %v", keys, mode, escapeCSI, gotLines, wantLines)
				}
			})

			t.Run("Input", func(t *testing.T) {
				// cleanup current Buffer after tests.
				defer clearBuffer(t, v, Buffer(0))

				b := v.NewBatch()

				const (
					keys      = `iabc<ESC>`
					mode      = `n`
					escapeCSI = false
				)
				var input string
				b.ReplaceTermcodes(keys, true, true, true, &input)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// clear current Buffer before run FeedKeys.
				clearBuffer(t, v, Buffer(0))

				var written int
				b.Input(input, &written)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if written != len(input) {
					t.Fatalf("Input(%s) = %d: want: %d", input, written, len(input))
				}

				wantLines := []byte{'a', 'b', 'c'}
				var gotLines []byte
				b.CurrentLine(&gotLines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(gotLines, wantLines) {
					t.Fatalf("FeedKeys(%s, %s, %t): got %v, want %v", input, mode, escapeCSI, gotLines, wantLines)
				}
			})

			t.Run("InputMouse", func(t *testing.T) {
				defer func() {
					// cleanup current Buffer after tests.
					clearBuffer(t, v, Buffer(0))

					input, err := v.ReplaceTermcodes(`<ESC>`, true, true, true)
					if err != nil {
						t.Fatal(err)
					}
					if err := v.FeedKeys(input, `n`, true); err != nil {
						t.Fatal(err)
					}
				}()

				// clear current Buffer before run FeedKeys.
				clearBuffer(t, v, Buffer(0))

				lines := [][]byte{
					[]byte("foo bar baz"),
					[]byte("qux quux quuz"),
					[]byte("corge grault garply"),
					[]byte("waldo fred plugh"),
					[]byte("xyzzy thud"),
				}
				if err := v.SetBufferLines(Buffer(0), 0, -1, true, lines); err != nil {
					t.Fatal(err)
				}

				const (
					button       = `left`
					firestAction = `press`
					secondAction = `release`
					modifier     = ""
				)
				const (
					wantGrid = 20
					wantRow  = 2
					wantCol  = 5
				)
				b := v.NewBatch()
				b.InputMouse(button, firestAction, modifier, wantGrid, wantRow, wantCol)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				b.InputMouse(button, secondAction, modifier, wantGrid, wantRow, wantCol)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				// TODO(zchee): assertion
			})

			t.Run("StringWidth", func(t *testing.T) {
				b := v.NewBatch()

				const str = "hello\t"
				var got int
				b.StringWidth(str, &got)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if got != len(str) {
					t.Fatalf("StringWidth(%s) = %d, want: %d", str, got, len(str))
				}
			})

			t.Run("KeyMap", func(t *testing.T) {
				b := v.NewBatch()

				mode := "n"
				b.SetKeyMap(mode, "y", "yy", make(map[string]bool))
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantMaps := []*Mapping{
					{
						LHS:     "y",
						RHS:     "yy",
						Silent:  0,
						NoRemap: 0,
						Expr:    0,
						Buffer:  0,
						SID:     0,
						NoWait:  0,
					},
				}
				wantMapsLen := 0
				if nvimVersion.Minor >= 6 {
					lastMap := wantMaps[0]
					wantMaps = []*Mapping{
						{
							LHS:     "<C-L>",
							RHS:     "<Cmd>nohlsearch|diffupdate<CR><C-L>",
							Silent:  0,
							NoRemap: 1,
							Expr:    0,
							Buffer:  0,
							SID:     0,
							NoWait:  0,
						},
						{
							LHS:     "Y",
							RHS:     "y$",
							Silent:  0,
							NoRemap: 1,
							Expr:    0,
							Buffer:  0,
							SID:     0,
							NoWait:  0,
						},
						lastMap,
					}
					wantMapsLen = 2
				}
				var got []*Mapping
				b.KeyMap(mode, &got)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, wantMaps) {
					for i, gotmap := range got {
						t.Logf(" got[%d]: %#v", i, gotmap)
					}
					for i, wantmap := range wantMaps {
						t.Logf("want[%d]: %#v", i, wantmap)
					}
					t.Fatalf("KeyMap(%s) = %#v, want: %#v", mode, got, wantMaps)
				}

				b.DeleteKeyMap(mode, "y")
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var got2 []*Mapping
				b.KeyMap(mode, &got2)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if len(got2) != wantMapsLen {
					t.Fatalf("expected %d but got %#v", wantMapsLen, got2)
				}
			})

			t.Run("BufferKeyMap", func(t *testing.T) {
				mode := "n"
				b := v.NewBatch()

				buffer := Buffer(0)
				b.SetBufferKeyMap(buffer, mode, "x", "xx", make(map[string]bool))
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				wantMap := []*Mapping{
					{
						LHS:     "x",
						RHS:     "xx",
						Silent:  0,
						NoRemap: 0,
						Expr:    0,
						Buffer:  1,
						SID:     0,
						NoWait:  0,
					},
				}
				var got []*Mapping
				b.BufferKeyMap(buffer, mode, &got)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(wantMap, got) {
					t.Fatalf("KeyMap(n) = %#v, want: %#v", got[0], wantMap[0])
				}

				b.DeleteBufferKeyMap(buffer, mode, "x")
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var got2 []*Mapping
				b.BufferKeyMap(buffer, mode, &got2)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if len(got2) > 0 {
					t.Fatalf("expected 0 but got %#v", got2)
				}
			})
		})
	}
}

func testEval(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			var a, b string
			if err := v.Eval(`["hello", "world"]`, []*string{&a, &b}); err != nil {
				t.Fatal(err)
			}
			if a != "hello" || b != "world" {
				t.Fatalf("a=%q b=%q, want a=hello b=world", a, b)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var x, y string
			b.Eval(`["hello", "world"]`, []*string{&x, &y})
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if x != "hello" || y != "world" {
				t.Fatalf("a=%q b=%q, want a=hello b=world", x, y)
			}
		})
	}
}

func testBatch(v *Nvim) func(*testing.T) {
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
			t.Fatalf("unxpected error %T %v", e, e)
		}
		// Expect results proceeding error.
		for i := 0; i < errorIndex; i++ {
			if results[i] != i {
				t.Fatalf("result[i] = %d, want %d", results[i], i)
				break
			}
		}
		// No results after error.
		for i := errorIndex; i < len(results); i++ {
			if results[i] != -1 {
				t.Fatalf("result[i] = %d, want %d", results[i], -1)
				break
			}
		}

		// Execute should return marshal error for argument that cannot be marshaled.
		b.SetVar("batch0", make(chan bool))
		if err := b.Execute(); err == nil || !strings.Contains(err.Error(), "chan bool") {
			t.Fatalf("err = nil, expect error containing text 'chan bool'")
		}

		// Test call with empty argument list.
		var buf Buffer
		b.CurrentBuffer(&buf)
		if err = b.Execute(); err != nil {
			t.Fatalf("GetCurrentBuffer returns err %s: %#v", err, err)
		}
	}
}

func testMode(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			m, err := v.Mode()
			if err != nil {
				t.Fatal(err)
			}
			if m.Mode != "n" {
				t.Fatalf("Mode() returned %s, want n", m.Mode)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var m Mode
			b.Mode(&m)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if m.Mode != "n" {
				t.Fatalf("Mode() returned %s, want n", m.Mode)
			}
		})
	}
}

func testExecLua(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			var n int
			err := v.ExecLua("local a, b = ... return a + b", &n, 1, 2)
			if err != nil {
				t.Fatal(err)
			}
			if n != 3 {
				t.Fatalf("Mode() returned %v, want 3", n)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var n int
			b.ExecLua("local a, b = ... return a + b", &n, 1, 2)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if n != 3 {
				t.Fatalf("Mode() returned %v, want 3", n)
			}
		})
	}
}

func testHighlight(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			cm, err := v.ColorMap()
			if err != nil {
				t.Fatal(err)
			}

			const cmd = `highlight NewHighlight cterm=underline ctermbg=green guifg=red guibg=yellow guisp=blue gui=bold`
			if err := v.Command(cmd); err != nil {
				t.Fatal(err)
			}

			wantCTerm := &HLAttrs{
				Underline:       true,
				Foreground:      -1,
				Background:      10,
				Special:         -1,
				CtermForeground: -1,
				CtermBackground: -1,
			}
			wantGUI := &HLAttrs{
				Bold:            true,
				Foreground:      cm["Red"],
				Background:      cm["Yellow"],
				Special:         cm["Blue"],
				CtermForeground: -1,
				CtermBackground: -1,
			}

			var nsID int
			if err := v.Eval(`hlID('NewHighlight')`, &nsID); err != nil {
				t.Fatal(err)
			}

			const HLIDName = `Error`
			var wantErrorHLID = 64
			if nvimVersion.Minor >= 7 {
				wantErrorHLID = 66
			}

			goHLID, err := v.HLIDByName(HLIDName)
			if err != nil {
				t.Fatal(err)
			}
			if goHLID != wantErrorHLID {
				t.Fatalf("HLByID(%s)\n got %+v,\nwant %+v", HLIDName, goHLID, wantErrorHLID)
			}

			gotCTermHL, err := v.HLByID(nsID, false)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotCTermHL, wantCTerm) {
				t.Fatalf("HLByID(id, false)\n got %+v,\nwant %+v", gotCTermHL, wantCTerm)
			}

			gotGUIHL, err := v.HLByID(nsID, true)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotGUIHL, wantGUI) {
				t.Fatalf("HLByID(id, true)\n got %+v,\nwant %+v", gotGUIHL, wantGUI)
			}

			errorMsgHL, err := v.HLByName(`ErrorMsg`, true)
			if err != nil {
				t.Fatal(err)
			}
			errorMsgHL.Bold = true
			errorMsgHL.Underline = true
			errorMsgHL.Italic = true
			if err := v.SetHighlight(nsID, "ErrorMsg", errorMsgHL); err != nil {
				t.Fatal(err)
			}

			wantErrorMsgEHL := &HLAttrs{
				Bold:            true,
				Underline:       true,
				Italic:          true,
				Foreground:      16777215,
				Background:      16711680,
				Special:         -1,
				CtermForeground: -1,
				CtermBackground: -1,
			}
			if !reflect.DeepEqual(wantErrorMsgEHL, errorMsgHL) {
				t.Fatalf("SetHighlight:\nwant %#v\n got %#v", wantErrorMsgEHL, errorMsgHL)
			}

			const cmd2 = `hi NewHighlight2 guifg=yellow guibg=red gui=bold`
			if err := v.Command(cmd2); err != nil {
				t.Fatal(err)
			}
			var nsID2 int
			if err := v.Eval(`hlID('NewHighlight2')`, &nsID2); err != nil {
				t.Fatal(err)
			}
			if err := v.SetHighlightNameSpace(nsID2); err != nil {
				t.Fatal(err)
			}
			want := &HLAttrs{
				Bold:            true,
				Underline:       false,
				Undercurl:       false,
				Italic:          false,
				Reverse:         false,
				Foreground:      16776960,
				Background:      16711680,
				Special:         -1,
				CtermForeground: -1,
				CtermBackground: -1,
			}
			got, err := v.HLByID(nsID2, true)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("SetHighlight:\nwant %#v\n got %#v", want, got)
			}

			const wantRedColor = 16711680
			gotColor, err := v.ColorByName("red")
			if err != nil {
				t.Fatal(err)
			}
			if wantRedColor != gotColor {
				t.Fatalf("expected red color %d but got %d", wantRedColor, gotColor)
			}

			id, err := v.AddBufferHighlight(Buffer(0), 0, `NewHighlight2`, 0, 0, -1)
			if err != nil {
				t.Fatal(err)
			}
			if id < 0 {
				t.Fatalf("want id is not negative but got %d", id)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var cm map[string]int
			b.ColorMap(&cm)

			const cmd = `highlight NewHighlight cterm=underline ctermbg=green guifg=red guibg=yellow guisp=blue gui=bold`
			b.Command(cmd)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			wantCTerm := &HLAttrs{
				Underline:       true,
				Foreground:      -1,
				Background:      10,
				Special:         -1,
				CtermForeground: -1,
				CtermBackground: -1,
			}
			wantGUI := &HLAttrs{
				Bold:            true,
				Foreground:      cm[`Red`],
				Background:      cm[`Yellow`],
				Special:         cm[`Blue`],
				CtermForeground: -1,
				CtermBackground: -1,
			}

			var nsID int
			b.Eval("hlID('NewHighlight')", &nsID)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			const HLIDName = `Error`
			var wantErrorHLID = 64
			if nvimVersion.Minor >= 7 {
				wantErrorHLID = 66
			}

			var goHLID int
			b.HLIDByName(HLIDName, &goHLID)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if goHLID != wantErrorHLID {
				t.Fatalf("HLByID(%s)\n got %+v,\nwant %+v", HLIDName, goHLID, wantErrorHLID)
			}

			var gotCTermHL HLAttrs
			b.HLByID(nsID, false, &gotCTermHL)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(&gotCTermHL, wantCTerm) {
				t.Fatalf("HLByID(id, false)\n got %+v,\nwant %+v", &gotCTermHL, wantCTerm)
			}

			var gotGUIHL HLAttrs
			b.HLByID(nsID, true, &gotGUIHL)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(&gotGUIHL, wantGUI) {
				t.Fatalf("HLByID(id, true)\n got %+v,\nwant %+v", &gotGUIHL, wantGUI)
			}

			var errorMsgHL HLAttrs
			b.HLByName(`ErrorMsg`, true, &errorMsgHL)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			errorMsgHL.Bold = true
			errorMsgHL.Underline = true
			errorMsgHL.Italic = true
			b.SetHighlight(nsID, `ErrorMsg`, &errorMsgHL)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			wantErrorMsgEHL := &HLAttrs{
				Bold:            true,
				Underline:       true,
				Italic:          true,
				Foreground:      16777215,
				Background:      16711680,
				Special:         -1,
				CtermForeground: -1,
				CtermBackground: -1,
			}
			if !reflect.DeepEqual(&errorMsgHL, wantErrorMsgEHL) {
				t.Fatalf("SetHighlight:\ngot %#v\nwant %#v", &errorMsgHL, wantErrorMsgEHL)
			}

			const cmd2 = `hi NewHighlight2 guifg=yellow guibg=red gui=bold`
			b.Command(cmd2)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			var nsID2 int
			b.Eval("hlID('NewHighlight2')", &nsID2)
			b.SetHighlightNameSpace(nsID2)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			want := &HLAttrs{
				Bold:            true,
				Underline:       false,
				Undercurl:       false,
				Italic:          false,
				Reverse:         false,
				Foreground:      16776960,
				Background:      16711680,
				Special:         -1,
				CtermForeground: -1,
				CtermBackground: -1,
			}

			var got HLAttrs
			b.HLByID(nsID2, true, &got)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(&got, want) {
				t.Fatalf("SetHighlight:\n got %#v\nwant %#v", &got, want)
			}

			const wantRedColor = 16711680
			var gotColor int
			b.ColorByName("red", &gotColor)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if wantRedColor != gotColor {
				t.Fatalf("expected red color %d but got %d", wantRedColor, gotColor)
			}

			var id int
			b.AddBufferHighlight(Buffer(0), 0, `NewHighlight2`, 0, 0, -1, &id)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if id < 0 {
				t.Fatalf("want id is not negative but got %d", id)
			}
		})
	}
}

func testVirtualText(v *Nvim) func(*testing.T) {
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

		chunks := []TextChunk{
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

func testFloatingWindow(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			clearBuffer(t, v, 0) // clear curret buffer text
			curwin, err := v.CurrentWindow()
			if err != nil {
				t.Fatal(err)
			}

			wantWidth := 50
			wantHeight := 20

			cfg := &WindowConfig{
				Relative:  "cursor",
				Anchor:    "NW",
				Width:     40,
				Height:    10,
				Row:       1,
				Col:       0,
				Focusable: true,
				Style:     "minimal",
				ZIndex:    50,
			}
			w, err := v.OpenWindow(Buffer(0), true, cfg)
			if err != nil {
				t.Fatal(err)
			}
			if curwin == w {
				t.Fatal("same window number: floating window not focused")
			}

			if err := v.SetWindowWidth(w, wantWidth); err != nil {
				t.Fatal(err)
			}

			if err := v.SetWindowHeight(w, wantHeight); err != nil {
				t.Fatal(err)
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

			wantWinConfig := &WindowConfig{
				Relative:  "editor",
				Anchor:    "NW",
				Width:     40,
				Height:    10,
				Row:       1,
				Focusable: false,
				ZIndex:    50,
			}
			if err := v.SetWindowConfig(w, wantWinConfig); err != nil {
				t.Fatal(err)
			}
			gotWinConfig, err := v.WindowConfig(w)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotWinConfig, wantWinConfig) {
				t.Fatalf("want %#v but got %#v", wantWinConfig, gotWinConfig)
			}

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
			if err := v.WindowOption(w, "number", &numberOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "relativenumber", &relativenumberOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "cursorline", &cursorlineOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "cursorcolumn", &cursorcolumnOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "spell", &spellOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "list", &listOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "signcolumn", &signcolumnOpt); err != nil {
				t.Fatal(err)
			}
			if err := v.WindowOption(w, "colorcolumn", &colorcolumnOpt); err != nil {
				t.Fatal(err)
			}
			if numberOpt || relativenumberOpt || cursorlineOpt || cursorcolumnOpt || spellOpt || listOpt || signcolumnOpt != "auto" || colorcolumnOpt != "" {
				t.Fatal("expected minimal style")
			}

			if err := v.CloseWindow(w, true); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			clearBuffer(t, v, 0) // clear curret buffer text

			b := v.NewBatch()
			var curwin Window
			b.CurrentWindow(&curwin)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			wantWidth := 50
			wantHeight := 20

			cfg := &WindowConfig{
				Relative:  "cursor",
				Anchor:    "NW",
				Width:     40,
				Height:    10,
				Row:       1,
				Col:       0,
				Focusable: true,
				Style:     "minimal",
				ZIndex:    50,
			}
			var w Window
			b.OpenWindow(Buffer(0), true, cfg, &w)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if curwin == w {
				t.Fatal("same window number: floating window not focused")
			}

			b.SetWindowWidth(w, wantWidth)
			b.SetWindowHeight(w, wantHeight)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			var gotWidth int
			b.WindowWidth(w, &gotWidth)
			var gotHeight int
			b.WindowHeight(w, &gotHeight)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if gotWidth != wantWidth {
				t.Fatalf("got %d width but want %d", gotWidth, wantWidth)
			}
			if gotHeight != wantHeight {
				t.Fatalf("got %d height but want %d", gotHeight, wantHeight)
			}

			wantWinConfig := &WindowConfig{
				Relative:  "editor",
				Anchor:    "NW",
				Width:     40,
				Height:    10,
				Row:       1,
				Focusable: false,
				ZIndex:    50,
			}
			b.SetWindowConfig(w, wantWinConfig)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			gotWinConfig := new(WindowConfig)
			b.WindowConfig(w, gotWinConfig)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotWinConfig, wantWinConfig) {
				t.Fatalf("want %#v but got %#v", wantWinConfig, gotWinConfig)
			}

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
			b.WindowOption(w, "number", &numberOpt)
			b.WindowOption(w, "relativenumber", &relativenumberOpt)
			b.WindowOption(w, "cursorline", &cursorlineOpt)
			b.WindowOption(w, "cursorcolumn", &cursorcolumnOpt)
			b.WindowOption(w, "spell", &spellOpt)
			b.WindowOption(w, "list", &listOpt)
			b.WindowOption(w, "signcolumn", &signcolumnOpt)
			b.WindowOption(w, "colorcolumn", &colorcolumnOpt)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if numberOpt || relativenumberOpt || cursorlineOpt || cursorcolumnOpt || spellOpt || listOpt || signcolumnOpt != "auto" || colorcolumnOpt != "" {
				t.Fatal("expected minimal style")
			}

			b.CloseWindow(w, true)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testContext(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
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
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var ctxt map[string]interface{}
			b.Context(make(map[string][]string), &ctxt)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			var result interface{}
			b.LoadContext(ctxt, &result)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if result != nil {
				t.Fatal("expected result to nil")
			}
		})
	}
}

func testExtmarks(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			// setup buffer lines
			lines := [][]byte{
				[]byte("hello"),
				[]byte("world"),
			}
			if err := v.SetBufferLines(Buffer(0), 0, -1, true, lines); err != nil {
				t.Fatal(err)
			}

			// create namespace for test extmarks
			const extMarkName = "test_extmarks"
			nsID, err := v.CreateNamespace(extMarkName)
			if err != nil {
				t.Fatal(err)
			}

			const (
				wantExtMarkID = 1
				wantLine      = 1
				wantCol       = 3
			)
			gotExtMarkID, err := v.SetBufferExtmark(Buffer(0), nsID, wantLine, wantCol, make(map[string]interface{}))
			if err != nil {
				t.Fatal(err)
			}
			if gotExtMarkID != wantExtMarkID {
				t.Fatalf("got %d extMarkID but want %d", gotExtMarkID, wantExtMarkID)
			}

			extmarks, err := v.BufferExtmarks(Buffer(0), nsID, 0, -1, make(map[string]interface{}))
			if err != nil {
				t.Fatal(err)
			}
			if len(extmarks) > 1 {
				t.Fatalf("expected extmarks length to 1 but got %d", len(extmarks))
			}
			if extmarks[0].ID != gotExtMarkID {
				t.Fatalf("got %d extMarkID but want %d", extmarks[0].ID, wantExtMarkID)
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

			deleted, err := v.DeleteBufferExtmark(Buffer(0), nsID, gotExtMarkID)
			if err != nil {
				t.Fatal(err)
			}
			if !deleted {
				t.Fatalf("expected deleted but got %t", deleted)
			}

			if err := v.ClearBufferNamespace(Buffer(0), nsID, 0, -1); err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() {
				clearBuffer(t, v, Buffer(0)) // clear curret buffer text
			})
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			// setup buffer lines
			lines := [][]byte{
				[]byte("hello"),
				[]byte("world"),
			}
			b.SetBufferLines(Buffer(0), 0, -1, true, lines)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			// create namespace for test extmarks
			const extMarkName = "test_extmarks"
			var nsID int
			b.CreateNamespace(extMarkName, &nsID)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			const (
				wantExtMarkID = 2
				wantLine      = 1
				wantCol       = 3
			)
			var gotExtMarkID int
			b.SetBufferExtmark(Buffer(0), nsID, wantLine, wantCol, make(map[string]interface{}), &gotExtMarkID)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if gotExtMarkID != wantExtMarkID {
				t.Fatalf("got %d extMarkID but want %d", gotExtMarkID, wantExtMarkID)
			}

			var extmarks []ExtMark
			b.BufferExtmarks(Buffer(0), nsID, 0, -1, make(map[string]interface{}), &extmarks)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if len(extmarks) > 1 {
				t.Fatalf("expected extmarks length to 1 but got %d", len(extmarks))
			}
			if extmarks[0].ID != gotExtMarkID {
				t.Fatalf("got %d extMarkID but want %d", extmarks[0].ID, wantExtMarkID)
			}
			if extmarks[0].Row != wantLine {
				t.Fatalf("got %d extmarks Row but want %d", extmarks[0].Row, wantLine)
			}
			if extmarks[0].Col != wantCol {
				t.Fatalf("got %d extmarks Col but want %d", extmarks[0].Col, wantCol)
			}

			var pos []int
			b.BufferExtmarkByID(Buffer(0), nsID, gotExtMarkID, make(map[string]interface{}), &pos)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if pos[0] != wantLine {
				t.Fatalf("got %d extMark line but want %d", pos[0], wantLine)
			}
			if pos[1] != wantCol {
				t.Fatalf("got %d extMark col but want %d", pos[1], wantCol)
			}

			var deleted bool
			b.DeleteBufferExtmark(Buffer(0), nsID, gotExtMarkID, &deleted)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !deleted {
				t.Fatalf("expected deleted but got %t", deleted)
			}

			b.ClearBufferNamespace(Buffer(0), nsID, 0, -1)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() {
				clearBuffer(t, v, Buffer(0)) // clear curret buffer text
			})
		})
	}
}

func testRuntime(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		var runtimePath string
		if err := v.Eval("$VIMRUNTIME", &runtimePath); err != nil {
			t.Fatal(err)
		}
		viDiff := filepath.Join(runtimePath, "doc", "vi_diff.txt")
		vimDiff := filepath.Join(runtimePath, "doc", "vim_diff.txt")
		want := fmt.Sprintf("%s,%s", viDiff, vimDiff)

		binaryPath, err := exec.LookPath(BinaryName)
		if err != nil {
			t.Fatal(err)
		}
		nvimPrefix := filepath.Dir(filepath.Dir(binaryPath))

		wantPaths := []string{
			filepath.Join(nvimPrefix, "share", "nvim", "runtime"),
			filepath.Join(nvimPrefix, "lib", "nvim"),
		}
		switch runtime.GOOS {
		case "linux", "darwin":
			if nvimVersion.Minor <= 5 {
				oldRuntimePaths := []string{
					filepath.Join("/etc", "xdg", "nvim"),
					filepath.Join("/etc", "xdg", "nvim", "after"),
					filepath.Join("/usr", "local", "share", "nvim", "site"),
					filepath.Join("/usr", "local", "share", "nvim", "site", "after"),
					filepath.Join("/usr", "share", "nvim", "site"),
					filepath.Join("/usr", "share", "nvim", "site", "after"),
				}
				wantPaths = append(wantPaths, oldRuntimePaths...)
			}
		case "windows":
			if nvimVersion.Minor <= 5 {
				localAppDataDir := os.Getenv("LocalAppData")
				oldRuntimePaths := []string{
					filepath.Join(localAppDataDir, "nvim"),
					filepath.Join(localAppDataDir, "nvim", "after"),
					filepath.Join(localAppDataDir, "nvim-data", "site"),
					filepath.Join(localAppDataDir, "nvim-data", "site", "after"),
				}
				wantPaths = append(wantPaths, oldRuntimePaths...)
			}
		}
		sort.Strings(wantPaths)

		argName := filepath.Join("doc", "*_diff.txt")
		argAll := true

		t.Run("Nvim", func(t *testing.T) {
			t.Run("RuntimeFiles", func(t *testing.T) {
				files, err := v.RuntimeFiles(argName, argAll)
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(files)

				if len(files) != 2 {
					t.Fatalf("expected 2 length but got %d", len(files))
				}
				if got := strings.Join(files, ","); !strings.EqualFold(got, want) {
					t.Fatalf("RuntimeFiles(%s, %t): got %s but want %s", argName, argAll, got, want)
				}
			})

			t.Run("RuntimePaths", func(t *testing.T) {
				paths, err := v.RuntimePaths()
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(paths)

				if got, want := strings.Join(paths, ","), strings.Join(wantPaths, ","); !strings.EqualFold(got, want) {
					t.Fatalf("RuntimePaths():\n got %v\nwant %v", paths, wantPaths)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("RuntimeFiles", func(t *testing.T) {
				b := v.NewBatch()

				var files []string
				b.RuntimeFiles(argName, true, &files)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				sort.Strings(files)

				if len(files) != 2 {
					t.Fatalf("expected 2 length but got %d", len(files))
				}
				if got := strings.Join(files, ","); !strings.EqualFold(got, want) {
					t.Fatalf("RuntimeFiles(%s, %t): got %s but want %s", argName, argAll, got, want)
				}
			})

			t.Run("RuntimePaths", func(t *testing.T) {
				b := v.NewBatch()

				var paths []string
				b.RuntimePaths(&paths)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				sort.Strings(paths)

				if got, want := strings.Join(paths, ","), strings.Join(wantPaths, ","); !strings.EqualFold(got, want) {
					t.Fatalf("RuntimePaths():\n got %v\nwant %v", paths, wantPaths)
				}
			})
		})
	}
}

func testPutPaste(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Put", func(t *testing.T) {
			t.Run("Nvim", func(t *testing.T) {
				clearBuffer(t, v, Buffer(0)) // clear curret buffer text

				replacement := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
				if err := v.SetBufferText(Buffer(0), 0, 0, 0, 0, replacement); err != nil {
					t.Fatal(err)
				}

				const putText = "qux"
				putLines := []string{putText}
				if err := v.Put(putLines, "l", true, true); err != nil {
					t.Fatal(err)
				}

				want := append(replacement, []byte(putText))

				lines, err := v.BufferLines(Buffer(0), 0, -1, true)
				if err != nil {
					t.Fatal(err)
				}
				wantLines := bytes.Join(want, []byte("\n"))
				gotLines := bytes.Join(lines, []byte("\n"))
				if !bytes.Equal(wantLines, gotLines) {
					t.Fatalf("expected %s but got %s", string(wantLines), string(gotLines))
				}
			})

			t.Run("Batch", func(t *testing.T) {
				clearBuffer(t, v, 0) // clear curret buffer text

				b := v.NewBatch()

				replacement := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
				b.SetBufferText(Buffer(0), 0, 0, 0, 0, replacement)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				const putText = "qux"
				putLines := []string{putText}
				b.Put(putLines, "l", true, true)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				want := append(replacement, []byte(putText))

				var lines [][]byte
				b.BufferLines(Buffer(0), 0, -1, true, &lines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				wantLines := bytes.Join(want, []byte("\n"))
				gotLines := bytes.Join(lines, []byte("\n"))
				if !bytes.Equal(wantLines, gotLines) {
					t.Fatalf("expected %s but got %s", string(wantLines), string(gotLines))
				}
			})
		})

		t.Run("Paste", func(t *testing.T) {
			t.Run("Nvim", func(t *testing.T) {
				clearBuffer(t, v, 0) // clear curret buffer text

				state, err := v.Paste("!!", true, 1) // starts the paste
				if err != nil {
					t.Fatal(err)
				}
				if !state {
					t.Fatal("expect continue to pasting")
				}
				state, err = v.Paste("foo", true, 2) // continues the paste
				if err != nil {
					t.Fatal(err)
				}
				if !state {
					t.Fatal("expect continue to pasting")
				}
				state, err = v.Paste("bar", true, 2) // continues the paste
				if err != nil {
					t.Fatal(err)
				}
				if !state {
					t.Fatal("expect continue to pasting")
				}
				state, err = v.Paste("baz", true, 3) // ends the paste
				if err != nil {
					t.Fatal(err)
				}
				if !state {
					t.Fatal("expect not canceled")
				}

				lines, err := v.CurrentLine()
				if err != nil {
					t.Fatal(err)
				}
				const want = "!foobarbaz!"
				if want != string(lines) {
					t.Fatalf("got %s current lines but want %s", string(lines), want)
				}
			})

			t.Run("Batch", func(t *testing.T) {
				clearBuffer(t, v, 0) // clear curret buffer text

				b := v.NewBatch()

				var state, state2, state3, state4 bool
				b.Paste("!!", true, 1, &state)   // starts the paste
				b.Paste("foo", true, 2, &state2) // starts the paste
				b.Paste("bar", true, 2, &state3) // starts the paste
				b.Paste("baz", true, 3, &state4) // ends the paste
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !state || !state2 || !state3 || !state4 {
					t.Fatal("expect continue to pasting")
				}

				var lines []byte
				b.CurrentLine(&lines)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				const want = "!foobarbaz!"
				if want != string(lines) {
					t.Fatalf("got %s current lines but want %s", string(lines), want)
				}
			})
		})
	}
}

func testNamespace(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Namespace", func(t *testing.T) {
			t.Run("Nvim", func(t *testing.T) {
				const nsName = "test-nvim"
				nsID, err := v.CreateNamespace(nsName)
				if err != nil {
					t.Fatal(err)
				}

				nsIDs, err := v.Namespaces()
				if err != nil {
					t.Fatal(err)
				}

				gotID, ok := nsIDs[nsName]
				if !ok {
					t.Fatalf("not fount %s namespace ID", nsName)
				}

				if gotID != nsID {
					t.Fatalf("nsID mismatched: got: %d want: %d", gotID, nsID)
				}
			})

			t.Run("Batch", func(t *testing.T) {
				b := v.NewBatch()

				const nsName = "test-batch"
				var nsID int
				b.CreateNamespace(nsName, &nsID)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var nsIDs map[string]int
				b.Namespaces(&nsIDs)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				gotID, ok := nsIDs[nsName]
				if !ok {
					t.Fatalf("not fount %s namespace ID", nsName)
				}

				if gotID != nsID {
					t.Fatalf("nsID mismatched: got: %d want: %d", gotID, nsID)
				}
			})
		})
	}
}

func testOptions(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Option", func(t *testing.T) {
			tests := map[string]struct {
				name string
				want interface{}
			}{
				"background": {
					name: "background",
					want: "dark",
				},
			}

			for name, tt := range tests {
				t.Run("Nvim/"+name, func(t *testing.T) {
					var got interface{}
					if err := v.Option(tt.name, &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(tt.want, got) {
						t.Fatalf("got %#v but want %#v", got, tt.want)
					}
				})
			}

			for name, tt := range tests {
				t.Run("Batch/"+name, func(t *testing.T) {
					b := v.NewBatch()

					var got interface{}
					b.Option(tt.name, &got)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(tt.want, got) {
						t.Fatalf("got %#v but want %#v", got, tt.want)
					}
				})
			}
		})

		t.Run("SetOption", func(t *testing.T) {
			tests := map[string]struct {
				name  string
				value interface{}
				want  interface{}
			}{
				"background": {
					name: "background",
					want: "light",
				},
			}

			for name, tt := range tests {
				t.Run("Nvim/"+name, func(t *testing.T) {
					if err := v.SetOption(tt.name, tt.want); err != nil {
						t.Fatal(err)
					}

					var got interface{}
					if err := v.Option(tt.name, &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(tt.want, got) {
						t.Fatalf("got %#v but want %#v", got, tt.want)
					}
				})
			}

			for name, tt := range tests {
				t.Run("Batch/"+name, func(t *testing.T) {
					b := v.NewBatch()
					b.SetOption(tt.name, tt.want)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}

					var got interface{}
					b.Option(tt.name, &got)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(tt.want, got) {
						t.Fatalf("got %#v but want %#v", got, tt.want)
					}
				})
			}
		})

		t.Run("OptionInfo", func(t *testing.T) {
			tests := map[string]struct {
				name string
				want *OptionInfo
			}{
				"filetype": {
					name: "filetype",
					want: &OptionInfo{
						Name:          "filetype",
						ShortName:     "ft",
						Type:          "string",
						Default:       "",
						WasSet:        false,
						LastSetSid:    0,
						LastSetLinenr: 0,
						LastSetChan:   0,
						Scope:         "buf",
						GlobalLocal:   false,
						CommaList:     false,
						FlagList:      false,
					},
				},
				"cmdheight": {
					name: "cmdheight",
					want: &OptionInfo{
						Name:          "cmdheight",
						ShortName:     "ch",
						Type:          "number",
						Default:       int64(1),
						WasSet:        false,
						LastSetSid:    0,
						LastSetLinenr: 0,
						LastSetChan:   0,
						Scope:         "global",
						GlobalLocal:   false,
						CommaList:     false,
						FlagList:      false,
					},
				},
				"hidden": {
					name: "hidden",
					want: &OptionInfo{
						Name:          "hidden",
						ShortName:     "hid",
						Type:          "boolean",
						Default:       true,
						WasSet:        false,
						LastSetSid:    0,
						LastSetLinenr: 0,
						LastSetChan:   0,
						Scope:         "global",
						GlobalLocal:   false,
						CommaList:     false,
						FlagList:      false,
					},
				},
			}

			for name, tt := range tests {
				t.Run("Nvim/"+name, func(t *testing.T) {
					if name == "hidden" {
						skipVersion(t, "v0.6.0")
					}

					got, err := v.OptionInfo(tt.name)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(tt.want, got) {
						t.Fatalf("got %#v but want %#v", got, tt.want)
					}
				})
			}

			for name, tt := range tests {
				t.Run("Batch/"+name, func(t *testing.T) {
					if name == "hidden" {
						skipVersion(t, "v0.6.0")
					}

					b := v.NewBatch()

					var got OptionInfo
					b.OptionInfo(tt.name, &got)
					if err := b.Execute(); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(tt.want, &got) {
						t.Fatalf("got %#v but want %#v", &got, tt.want)
					}
				})
			}
		})
	}
}

func testAllOptionsInfo(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		want := &OptionInfo{
			Name:          "",
			ShortName:     "",
			Type:          "",
			Default:       nil,
			WasSet:        false,
			LastSetSid:    0,
			LastSetLinenr: 0,
			LastSetChan:   0,
			Scope:         "",
			GlobalLocal:   false,
			CommaList:     false,
			FlagList:      false,
		}

		got, err := v.AllOptionsInfo()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("got %v but want %v", got, want)
		}

		b := v.NewBatch()
		var got2 OptionInfo
		b.AllOptionsInfo(&got2)
		if err := b.Execute(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, &got2) {
			t.Fatalf("got %v but want %v", got2, want)
		}
	}
}

func testOptionsInfo(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		tests := map[string]struct {
			name string
			want *OptionInfo
		}{
			"filetype": {
				name: "filetype",
				want: &OptionInfo{
					Name:          "filetype",
					ShortName:     "ft",
					Type:          "string",
					Default:       "",
					WasSet:        false,
					LastSetSid:    0,
					LastSetLinenr: 0,
					LastSetChan:   0,
					Scope:         "buf",
					GlobalLocal:   false,
					CommaList:     false,
					FlagList:      false,
				},
			},
			"cmdheight": {
				name: "cmdheight",
				want: &OptionInfo{
					Name:          "cmdheight",
					ShortName:     "ch",
					Type:          "number",
					Default:       int64(1),
					WasSet:        false,
					LastSetSid:    0,
					LastSetLinenr: 0,
					LastSetChan:   0,
					Scope:         "global",
					GlobalLocal:   false,
					CommaList:     false,
					FlagList:      false,
				},
			},
			"hidden": {
				name: "hidden",
				want: &OptionInfo{
					Name:          "hidden",
					ShortName:     "hid",
					Type:          "boolean",
					Default:       true,
					WasSet:        false,
					LastSetSid:    0,
					LastSetLinenr: 0,
					LastSetChan:   0,
					Scope:         "global",
					GlobalLocal:   false,
					CommaList:     false,
					FlagList:      false,
				},
			},
		}

		for name, tt := range tests {
			t.Run("Nvim/"+name, func(t *testing.T) {
				if name == "hidden" {
					skipVersion(t, "v0.6.0")
				}

				got, err := v.OptionInfo(tt.name)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tt.want, got) {
					t.Fatalf("got %#v but want %#v", got, tt.want)
				}
			})
		}

		for name, tt := range tests {
			t.Run("Batch/"+name, func(t *testing.T) {
				if name == "hidden" {
					skipVersion(t, "v0.6.0")
				}

				b := v.NewBatch()

				var got OptionInfo
				b.OptionInfo(tt.name, &got)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(tt.want, &got) {
					t.Fatalf("got %#v but want %#v", &got, tt.want)
				}
			})
		}
	}
}

func testOptionsValue(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		tests := map[string]struct {
			name  string
			opts  map[string]OptionValueScope
			want  interface{}
			value interface{}
		}{
			"equalalways": {
				name: "equalalways",
				opts: map[string]OptionValueScope{
					"scope": GlobalScope,
				},
				want:  true,
				value: false,
			},
			"lazyredraw": {
				name: "lazyredraw",
				opts: map[string]OptionValueScope{
					"scope": LocalScope,
				},
				want:  false,
				value: true,
			},
		}
		for name, tt := range tests {
			t.Run(path.Join(name, "Nvim"), func(t *testing.T) {
				skipVersion(t, "v0.7.0")

				var result interface{}
				if err := v.OptionValue(tt.name, tt.opts, &result); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(result, tt.want) {
					t.Fatalf("got %#v but want %#v", result, tt.want)
				}

				if err := v.SetOptionValue(tt.name, tt.value, tt.opts); err != nil {
					t.Fatal(err)
				}

				var result2 interface{}
				if err := v.OptionValue(tt.name, tt.opts, &result2); err != nil {
					t.Fatal(err)
				}
				if reflect.DeepEqual(result, result2) {
					t.Fatalf("got %#v but want %#v", result, result2)
				}

				if err := v.SetOptionValue(tt.name, result, tt.opts); err != nil {
					t.Fatal(err)
				}
			})
		}

		for name, tt := range tests {
			t.Run(path.Join(name, "Batch"), func(t *testing.T) {
				skipVersion(t, "v0.7.0")

				b := v.NewBatch()

				var result interface{}
				b.OptionValue(tt.name, tt.opts, &result)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(result, tt.want) {
					t.Fatalf("got %#v but want %#v", result, tt.want)
				}

				b.SetOptionValue(tt.name, tt.value, tt.opts)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				var result2 interface{}
				b.OptionValue(tt.name, tt.opts, &result2)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if reflect.DeepEqual(result, result2) {
					t.Fatalf("got %#v but want %#v", result, result2)
				}

				b.SetOptionValue(tt.name, result, tt.opts)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

// TODO(zchee): correct testcase
func testTerm(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			buf, err := v.CreateBuffer(true, true)
			if err != nil {
				t.Fatal(err)
			}

			cfg := &WindowConfig{
				Relative: "editor",
				Width:    79,
				Height:   31,
				Row:      1,
				Col:      1,
				ZIndex:   50,
			}
			if _, err := v.OpenWindow(buf, false, cfg); err != nil {
				t.Fatal(err)
			}

			termID, err := v.OpenTerm(buf, make(map[string]interface{}))
			if err != nil {
				t.Fatal(err)
			}

			data := "\x1b[38;2;00;00;255mTRUECOLOR\x1b[0m"
			if err := v.Call("chansend", nil, termID, data); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var buf Buffer
			b.CreateBuffer(true, true, &buf)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			cfg := &WindowConfig{
				Relative: "editor",
				Width:    79,
				Height:   31,
				Row:      1,
				Col:      1,
				ZIndex:   50,
			}
			var win Window
			b.OpenWindow(buf, false, cfg, &win)

			var termID int
			b.OpenTerm(buf, make(map[string]interface{}), &termID)

			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			data := "\x1b[38;2;00;00;255mTRUECOLOR\x1b[0m"
			b.Call("chansend", nil, termID, data)

			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testChannelClientInfo(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		const clientNamePrefix = "testClient"

		var (
			clientVersion = ClientVersion{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: "-dev",
				Commit:     "e07b9dde387bc817d36176bbe1ce58acd3c81921",
			}
			clientType    = RemoteClientType
			clientMethods = map[string]*ClientMethod{
				"foo": {
					Async: true,
					NArgs: ClientMethodNArgs{
						Min: 0,
						Max: 1,
					},
				},
				"bar": {
					Async: false,
					NArgs: ClientMethodNArgs{
						Min: 0,
						Max: 0,
					},
				},
			}
			clientAttributes = ClientAttributes{
				ClientAttributeKeyLicense: "Apache-2.0",
			}
		)

		t.Run("Nvim", func(t *testing.T) {
			t.Run("Channels", func(t *testing.T) {
				chans, err := v.Channels()
				if err != nil {
					t.Fatal(err)
				}

				if len(chans) == 0 {
					t.Fatalf("want any channels but got %d", len(chans))
				}
			})

			clientName := clientNamePrefix + "Nvim"

			t.Run("SetClientInfo", func(t *testing.T) {
				if err := v.SetClientInfo(clientName, clientVersion, clientType, clientMethods, clientAttributes); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("ChannelInfo", func(t *testing.T) {
				wantClient := &Client{
					Name:       clientName,
					Version:    clientVersion,
					Type:       clientType,
					Methods:    clientMethods,
					Attributes: clientAttributes,
				}
				wantChannel := &Channel{
					Stream: "stdio",
					Mode:   "rpc",
					Client: wantClient,
				}

				gotChannel, err := v.ChannelInfo(int(channelID))
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(gotChannel, wantChannel) {
					t.Fatalf("got %#v channel but want %#v channel", gotChannel, wantChannel)
				}
			})
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			t.Run("Channels", func(t *testing.T) {
				var chans []*Channel
				b.Channels(&chans)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if len(chans) == 0 {
					t.Fatalf("want any channels but got %d", len(chans))
				}
			})

			clientName := clientNamePrefix + "Batch"

			t.Run("SetClientInfo", func(t *testing.T) {
				b.SetClientInfo(clientName, clientVersion, clientType, clientMethods, clientAttributes)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("ChannelInfo", func(t *testing.T) {
				wantClient := &Client{
					Name:       clientName,
					Version:    clientVersion,
					Type:       clientType,
					Methods:    clientMethods,
					Attributes: clientAttributes,
				}
				wantChannel := &Channel{
					Stream: "stdio",
					Mode:   "rpc",
					Client: wantClient,
				}

				var gotChannel Channel
				b.ChannelInfo(int(channelID), &gotChannel)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(&gotChannel, wantChannel) {
					t.Fatalf("got %#v channel but want %#v channel", &gotChannel, wantChannel)
				}
			})
		})
	}
}

func testUI(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			gotUIs, err := v.UIs()
			if err != nil {
				t.Fatal(err)
			}

			if len(gotUIs) > 0 || gotUIs != nil {
				t.Fatalf("expected ui empty but non-zero: %#v", gotUIs)
			}

			v.RegisterHandler("redraw", func(updates ...[]interface{}) {})
			if err := v.AttachUI(500, 400, make(map[string]interface{})); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := v.DetachUI(); err != nil {
					t.Fatal(err)
				}
			})

			if err := v.TryResizeUI(50, 40); err != nil {
				t.Fatal(err)
			}

			if err := v.SetUIOption("rgb", true); err != nil {
				t.Fatal(err)
			}

			if err := v.TryResizeUIGrid(1, 50, 40); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			var gotUIs []*UI
			b.UIs(&gotUIs)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if len(gotUIs) > 0 || gotUIs != nil {
				t.Fatalf("expected ui empty but non-zero: %#v", gotUIs)
			}

			v.RegisterHandler("redraw", func(updates ...[]interface{}) {})
			b.AttachUI(500, 400, make(map[string]interface{}))
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				b.DetachUI()
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
			})

			b.TryResizeUI(50, 40)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			b.SetUIOption("rgb", true)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			b.TryResizeUIGrid(1, 50, 40)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testProc(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Run("Proc", func(t *testing.T) {
				pid := os.Getpid()
				ppid := os.Getppid()

				wantProcess := &Process{
					Name: "nvim.test",
					PID:  pid,
					PPID: ppid,
				}
				if runtime.GOOS == "windows" {
					wantProcess.Name = "nvim.test.exe"
				}

				gotProc, err := v.Proc(pid)
				if err != nil {
					t.Fatal(err)
				}

				if gotProc.Name != wantProcess.Name {
					t.Fatalf("got %s Process.Name but want %s", gotProc.Name, wantProcess.Name)
				}

				if gotProc.PID != wantProcess.PID {
					t.Fatalf("got %d Process.PID but want %d", gotProc.PID, wantProcess.PID)
				}

				if gotProc.PPID != wantProcess.PPID {
					t.Fatalf("got %d Process.PPID but want %d", gotProc.PPID, wantProcess.PPID)
				}
			})

			t.Run("ProcChildren", func(t *testing.T) {
				_, err := v.ProcChildren(os.Getpid())
				if err != nil {
					t.Fatal(err)
				}
				// TODO(zchee): assert processes
			})
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("Proc", func(t *testing.T) {
				b := v.NewBatch()

				pid := os.Getpid()
				ppid := os.Getppid()

				wantProcess := &Process{
					Name: "nvim.test",
					PID:  pid,
					PPID: ppid,
				}
				if runtime.GOOS == "windows" {
					wantProcess.Name = "nvim.test.exe"
				}

				var gotProc Process
				b.Proc(pid, &gotProc)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}

				if gotProc.Name != wantProcess.Name {
					t.Fatalf("got %s Process.Name but want %s", gotProc.Name, wantProcess.Name)
				}

				if gotProc.PID != wantProcess.PID {
					t.Fatalf("got %d Process.PID but want %d", gotProc.PID, wantProcess.PID)
				}

				if gotProc.PPID != wantProcess.PPID {
					t.Fatalf("got %d Process.PPID but want %d", gotProc.PPID, wantProcess.PPID)
				}
			})

			t.Run("ProcChildren", func(t *testing.T) {
				b := v.NewBatch()

				var processes []uint
				b.ProcChildren(os.Getpid(), &processes)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				// TODO(zchee): assert processes
			})
		})
	}
}

func testMark(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			// set dummy lines
			lines := [][]byte{
				[]byte("a"),
				[]byte("bit of"),
				[]byte("text"),
			}
			if err := v.SetBufferLines(Buffer(0), 0, -1, true, lines); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				clearBuffer(t, v, Buffer(0))
			})

			// set cursor position
			if err := v.SetWindowCursor(Window(0), [2]int{3, 0}); err != nil {
				t.Fatal(err)
			}

			// set buffer name
			cwd, err := os.Getwd() // buffer name is full path
			if err != nil {
				t.Fatal(err)
			}
			// shrink home path
			cwd = strings.ReplaceAll(cwd, os.Getenv("HOME"), "~")
			bufName := filepath.Join(cwd, "test_mark_buffer")
			if err := v.SetBufferName(Buffer(0), bufName); err != nil {
				t.Fatal(err)
			}

			// set mark
			const mark = "X"
			if err := v.Command(fmt.Sprintf("mark %s", mark)); err != nil {
				t.Fatal(err)
			}

			gotMark, err := v.Mark(mark, make(map[string]interface{}))
			if err != nil {
				t.Fatal(err)
			}
			wantMark := &Mark{
				Row:        3,
				Col:        0,
				Buffer:     0,
				BufferName: bufName,
			}

			if !reflect.DeepEqual(gotMark, wantMark) {
				t.Fatalf("got %#v mark but want %#v", gotMark, wantMark)
			}

			deleted, err := v.DeleteMark(mark)
			if err != nil {
				t.Fatal(err)
			}
			if !deleted {
				t.Fatalf("could not delete %s mark", mark)
			}

			gotMark2, err := v.Mark(mark, make(map[string]interface{}))
			if err != nil {
				t.Fatal(err)
			}
			wantMark2 := &Mark{
				Row:        0,
				Col:        0,
				Buffer:     0,
				BufferName: "",
			}

			if !reflect.DeepEqual(gotMark2, wantMark2) {
				t.Fatalf("got %#v mark but want %#v", gotMark2, wantMark2)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			// set dummy lines
			lines := [][]byte{
				[]byte("a"),
				[]byte("bit of"),
				[]byte("text"),
			}
			b.SetBufferLines(Buffer(0), 0, -1, true, lines)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				clearBuffer(t, v, Buffer(0))
			})

			// set cursor position
			b.SetWindowCursor(Window(0), [2]int{3, 0})
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			// set buffer name
			cwd, err := os.Getwd() // buffer name is full path
			if err != nil {
				t.Fatal(err)
			}
			// shrink home path
			cwd = strings.ReplaceAll(cwd, os.Getenv("HOME"), "~")
			bufName := filepath.Join(cwd, "test_mark_buffer")
			b.SetBufferName(Buffer(0), bufName)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			// set mark
			const mark = "X"
			b.Command(fmt.Sprintf("mark %s", mark))
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			gotMark := new(Mark)
			b.Mark(mark, make(map[string]interface{}), gotMark)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			wantMark := &Mark{
				Row:        3,
				Col:        0,
				Buffer:     0,
				BufferName: bufName,
			}

			if !reflect.DeepEqual(gotMark, wantMark) {
				t.Fatalf("got %#v mark but want %#v", gotMark, wantMark)
			}

			var deleted bool
			b.DeleteMark(mark, &deleted)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !deleted {
				t.Fatalf("could not delete %s mark", mark)
			}

			gotMark2 := new(Mark)
			b.Mark(mark, make(map[string]interface{}), gotMark2)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			wantMark2 := &Mark{
				Row:        0,
				Col:        0,
				Buffer:     0,
				BufferName: "",
			}

			if !reflect.DeepEqual(gotMark2, wantMark2) {
				t.Fatalf("got %#v mark but want %#v", gotMark2, wantMark2)
			}
		})
	}
}

func testStatusLine(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			opts := map[string]interface{}{
				"highlights": true,
			}
			gotStatusLine, err := v.EvalStatusLine("TextWithNoHighlight", opts)
			if err != nil {
				t.Fatal(err)
			}

			wantStatusLine := map[string]interface{}{
				"highlights": []interface{}{
					map[string]interface{}{
						"group": "StatusLine",
						"start": int64(0),
					},
				},
				"str":   "TextWithNoHighlight",
				"width": 19,
			}

			gotHighlight := gotStatusLine["highlights"].([]interface{})[0].(map[string]interface{})
			wantHighlight := wantStatusLine["highlights"].([]interface{})[0].(map[string]interface{})
			if !reflect.DeepEqual(gotHighlight["group"], wantHighlight["group"]) {
				t.Fatalf("got %#v highlight group but want %#v", gotHighlight["group"], wantHighlight["group"])
			}
			if gotHighlight["start"] != wantHighlight["start"] {
				t.Fatalf("got %#v highlight start but want %#v", gotHighlight["start"], wantHighlight["start"])
			}

			gotStr := gotStatusLine["str"]
			wantStr := gotStatusLine["str"]
			if gotStr != wantStr {
				t.Fatalf("got %#v str but want %#v", gotStr, wantStr)
			}

			gotWidth := gotStatusLine["width"]
			wantWidth := gotStatusLine["width"]
			if gotWidth != wantWidth {
				t.Fatalf("got %#v width but want %#v", gotWidth, wantWidth)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			opts := map[string]interface{}{
				"highlights": true,
			}
			var gotStatusLine map[string]interface{}
			b.EvalStatusLine("TextWithNoHighlight", opts, &gotStatusLine)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			wantStatusLine := map[string]interface{}{
				"highlights": []interface{}{
					map[string]interface{}{
						"group": "StatusLine",
						"start": int64(0),
					},
				},
				"str":   "TextWithNoHighlight",
				"width": 19,
			}

			gotHighlight := gotStatusLine["highlights"].([]interface{})[0].(map[string]interface{})
			wantHighlight := wantStatusLine["highlights"].([]interface{})[0].(map[string]interface{})
			if !reflect.DeepEqual(gotHighlight["group"], wantHighlight["group"]) {
				t.Fatalf("got %#v highlight group but want %#v", gotHighlight["group"], wantHighlight["group"])
			}
			if gotHighlight["start"] != wantHighlight["start"] {
				t.Fatalf("got %#v highlight start but want %#v", gotHighlight["start"], wantHighlight["start"])
			}

			gotStr := gotStatusLine["str"]
			wantStr := gotStatusLine["str"]
			if gotStr != wantStr {
				t.Fatalf("got %#v str but want %#v", gotStr, wantStr)
			}

			gotWidth := gotStatusLine["width"]
			wantWidth := gotStatusLine["width"]
			if gotWidth != wantWidth {
				t.Fatalf("got %#v width but want %#v", gotWidth, wantWidth)
			}
		})
	}
}
