package nvim

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newChildProcess(tb testing.TB) (v *Nvim, cleanup func()) {
	tb.Helper()

	ctx := context.Background()
	opts := []ChildProcessOption{
		ChildProcessArgs("-u", "NONE", "-n", "--embed", "--headless", "--noplugin"),
		ChildProcessContext(ctx),
		ChildProcessLogf(tb.Logf),
	}
	if runtime.GOOS == "windows" {
		opts = append(opts, ChildProcessCommand("nvim.exe"))
	}
	n, err := NewChildProcess(opts...)
	if err != nil {
		tb.Fatal(err)
	}
	v = n

	done := make(chan error, 1)
	go func() {
		done <- v.Serve()
	}()

	cleanup = func() {
		if err := v.Close(); err != nil {
			tb.Fatal(err)
		}

		err := <-done
		if err != nil {
			tb.Fatal(err)
		}

		const nvimlogFile = ".nvimlog"
		wd, err := os.Getwd()
		if err != nil {
			tb.Fatal(err)
		}
		if walkErr := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if fname := info.Name(); fname == nvimlogFile {
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to remove %s file: %w", path, err)
				}
			}

			return nil
		}); walkErr != nil {
			tb.Fatal(fmt.Errorf("walkErr: %w", errors.Unwrap(walkErr)))
		}
	}

	return v, cleanup
}

type version struct {
	Major int
	Minor int
	Patch int
}

var nvimVersion version

func parseVersion(tb testing.TB, version string) (major, minor, patch int) {
	tb.Helper()

	version = strings.TrimPrefix(version, "v")
	vpair := strings.Split(version, ".")
	if len(vpair) != 3 {
		tb.Fatal("could not parse neovim version")
	}

	var err error
	major, err = strconv.Atoi(vpair[0])
	if err != nil {
		tb.Fatal(err)
	}
	minor, err = strconv.Atoi(vpair[1])
	if err != nil {
		tb.Fatal(err)
	}
	patch, err = strconv.Atoi(vpair[2])
	if err != nil {
		tb.Fatal(err)
	}

	return major, minor, patch
}

func skipVersion(tb testing.TB, version string) {
	major, minor, patch := parseVersion(tb, version)

	const skipFmt = "skip test: current neovim version v%d.%d.%d but expected version %s"
	if nvimVersion.Major < major || nvimVersion.Minor < minor || nvimVersion.Patch < patch {
		tb.Skipf(skipFmt, nvimVersion.Major, nvimVersion.Minor, nvimVersion.Patch, version)
	}
}

func TestAPI(t *testing.T) {
	t.Parallel()

	v, cleanup := newChildProcess(t)
	defer cleanup()

	apiInfo, err := v.APIInfo()
	if err != nil {
		t.Fatal(err)
	}
	if len(apiInfo) != 2 {
		t.Fatalf("unknown APIInfo: %#v", apiInfo)
	}
	info, ok := apiInfo[1].(map[string]interface{})
	if !ok {
		t.Fatalf("apiInfo[1] is not map[string]interface{} type: %T", apiInfo[1])
	}
	infoV := info["version"].(map[string]interface{})
	nvimVersion.Major = int(infoV["major"].(int64))
	nvimVersion.Minor = int(infoV["minor"].(int64))
	nvimVersion.Patch = int(infoV["patch"].(int64))

	t.Run("BufAttach", testBufAttach(v))
	t.Run("SimpleHandler", testSimpleHandler(v))
	t.Run("Buffer", testBuffer(v))
	t.Run("Window", testWindow(v))
	t.Run("Tabpage", testTabpage(v))
	t.Run("Lines", testLines(v))
	t.Run("Var", testVar(v))
	t.Run("Message", testMessage(v))
	t.Run("StructValue", testStructValue(v))
	t.Run("Eval", testEval(v))
	t.Run("Batch", testBatch(v))
	t.Run("CallWithNoArgs", testCallWithNoArgs(v))
	t.Run("Mode", testMode(v))
	t.Run("ExecLua", testExecLua(v))
	t.Run("Highlight", testHighlight(v))
	t.Run("VirtualText", testVirtualText(v))
	t.Run("FloatingWindow", testFloatingWindow(v))
	t.Run("Context", testContext(v))
	t.Run("Extmarks", testExtmarks(v))
	t.Run("RuntimeFiles", testRuntimeFiles(v))
	t.Run("AllOptionsInfo", testAllOptionsInfo(v))
	t.Run("OptionsInfo", testOptionsInfo(v))
	t.Run("OpenTerm", testTerm(v))
}

func testBufAttach(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
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
			t.Fatal(errors.New("could not attach buffer"))
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
			t.Run("Buffers", func(t *testing.T) {
				bufs, err := v.Buffers()
				if err != nil {
					t.Fatal(err)
				}
				if len(bufs) != 1 {
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
		})

		t.Run("Batch", func(t *testing.T) {
			t.Run("Buffers", func(t *testing.T) {
				b := v.NewBatch()

				var bufs []Buffer
				b.Buffers(&bufs)
				if err := b.Execute(); err != nil {
					t.Fatal(err)
				}
				if len(bufs) != 1 {
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
		})
	}
}

func testTabpage(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Parallel()

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
		})

		t.Run("Batch", func(t *testing.T) {
			t.Parallel()

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
		})
	}
}

func testLines(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
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

func testVar(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			if err := v.SetVar("gvar", "gval"); err != nil {
				t.Fatal(err)
			}

			var value interface{}
			if err := v.Var("gvar", &value); err != nil {
				t.Fatal(err)
			}
			if value != "gval" {
				t.Fatalf("got %v, want %q", value, "gval")
			}

			if err := v.SetVar("gvar", ""); err != nil {
				t.Fatal(err)
			}
			value = nil
			if err := v.Var("gvar", &value); err != nil {
				t.Fatal(err)
			}
			if value != "" {
				t.Fatalf("got %v, want %q", value, "")
			}
		})

		t.Run("Batch", func(t *testing.T) {
			b := v.NewBatch()

			b.SetVar("gvar", "gval")
			var value interface{}
			b.Var("gvar", &value)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if value != "gval" {
				t.Fatalf("got %v, want %q", value, "gval")
			}

			b.SetVar("gvar", "")
			value = nil
			b.Var("gvar", &value)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if value != "" {
				t.Fatalf("got %v, want %q", value, "")
			}
		})
	}
}

func testMessage(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
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

			// cleanup v:statusmsg
			if err := v.SetVVar("statusmsg", ""); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

			// clear messages
			if _, err := v.Exec(":messages clear", false); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

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

			// cleanup v:statusmsg
			if err := v.SetVVar("statusmsg", ""); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

			// clear messages
			if _, err := v.Exec(":messages clear", false); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

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

			// clear messages
			if _, err := v.Exec(":messages clear", false); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

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

			// clear messages
			if _, err := v.Exec(":messages clear", false); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

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

			// clear messages
			if _, err := v.Exec(":messages clear", false); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}
		})

		t.Run("Batch", func(t *testing.T) {
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

			// cleanup v:statusmsg
			if err := v.SetVVar("statusmsg", ""); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

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

			// cleanup v:statusmsg
			if err := v.SetVVar("statusmsg", ""); err != nil {
				t.Fatalf("failed to SetVVar: %v", err)
			}

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

			// clear messages
			b.Exec(":messages clear", false, new(string))
			if err := b.Execute(); err != nil {
				t.Fatalf("failed to \":messages clear\" command: %v", err)
			}

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

			// clear messages
			b.Exec(":messages clear", false, new(string))
			b.SetVVar("statusmsg", "")
			if err := b.Execute(); err != nil {
				t.Fatalf("failed to \":messages clear\" command: %v", err)
			}

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

			// clear messages
			b.Exec(":messages clear", false, new(string))
			if err := b.Execute(); err != nil {
				t.Fatalf("failed to \":messages clear\" command: %v", err)
			}
		})
	}
}

func testStructValue(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Parallel()

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
				t.Fatalf("got %+v, want %+v", &actual, &expected)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			t.Parallel()

			b := v.NewBatch()

			var expected, actual struct {
				Str string
				Num int
			}
			expected.Str = "Hello"
			expected.Num = 42
			b.SetVar("structvar", &expected)
			b.Var("structvar", &actual)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(&actual, &expected) {
				t.Fatalf("got %+v, want %+v", &actual, &expected)
			}
		})
	}
}

func testEval(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Parallel()

			var a, b string
			if err := v.Eval(`["hello", "world"]`, []*string{&a, &b}); err != nil {
				t.Fatal(err)
			}
			if a != "hello" || b != "world" {
				t.Fatalf("a=%q b=%q, want a=hello b=world", a, b)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			t.Parallel()

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

func testCallWithNoArgs(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		var wd string
		err := v.Call("getcwd", &wd)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testMode(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		m, err := v.Mode()
		if err != nil {
			t.Fatal(err)
		}
		if m.Mode != "n" {
			t.Fatalf("Mode() returned %s, want n", m.Mode)
		}
	}
}

func testExecLua(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Parallel()

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
			t.Parallel()

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
			t.Parallel()

			cm, err := v.ColorMap()
			if err != nil {
				t.Fatal(err)
			}

			const cmd = `highlight NewHighlight cterm=underline ctermbg=green guifg=red guibg=yellow guisp=blue gui=bold`
			if err := v.Command(cmd); err != nil {
				t.Fatal(err)
			}

			wantCTerm := &HLAttrs{
				Underline:  true,
				Foreground: -1,
				Background: 10,
				Special:    -1,
			}
			wantGUI := &HLAttrs{
				Bold:       true,
				Foreground: cm["Red"],
				Background: cm["Yellow"],
				Special:    cm["Blue"],
			}

			var nsID int
			if err := v.Eval(`hlID('NewHighlight')`, &nsID); err != nil {
				t.Fatal(err)
			}

			const (
				HLIDName      = "Error"
				wantErrorHLID = 137
			)
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
				Bold:       true,
				Underline:  true,
				Italic:     true,
				Foreground: 16777215,
				Background: 16711680,
				Special:    -1,
			}
			if !reflect.DeepEqual(wantErrorMsgEHL, errorMsgHL) {
				t.Fatalf("SetHighlight:\nwant %#v\n got %#v", wantErrorMsgEHL, errorMsgHL)
			}

			const cmd2 = "hi NewHighlight2 guifg=yellow guibg=red gui=bold"
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
				Bold:       true,
				Underline:  false,
				Undercurl:  false,
				Italic:     false,
				Reverse:    false,
				Foreground: 16776960,
				Background: 16711680,
				Special:    -1,
			}
			got, err := v.HLByID(nsID2, true)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("SetHighlight:\nwant %#v\n got %#v", want, got)
			}
		})

		t.Run("Batch", func(t *testing.T) {
			t.Parallel()

			b := v.NewBatch()

			var cm map[string]int
			b.ColorMap(&cm)

			const cmd = `highlight NewHighlight cterm=underline ctermbg=green guifg=red guibg=yellow guisp=blue gui=bold`
			b.Command(cmd)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			wantCTerm := &HLAttrs{
				Underline:  true,
				Foreground: -1,
				Background: 10,
				Special:    -1,
			}
			wantGUI := &HLAttrs{
				Bold:       true,
				Foreground: cm[`Red`],
				Background: cm[`Yellow`],
				Special:    cm[`Blue`],
			}

			var nsID int
			b.Eval("hlID('NewHighlight')", &nsID)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}

			const (
				HLIDName      = `Error`
				wantErrorHLID = 137
			)
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
				Bold:       true,
				Underline:  true,
				Italic:     true,
				Foreground: 16777215,
				Background: 16711680,
				Special:    -1,
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
				Bold:       true,
				Underline:  false,
				Undercurl:  false,
				Italic:     false,
				Reverse:    false,
				Foreground: 16776960,
				Background: 16711680,
				Special:    -1,
			}

			var got HLAttrs
			b.HLByID(nsID2, true, &got)
			if err := b.Execute(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(&got, want) {
				t.Fatalf("SetHighlight:\n got %#v\nwant %#v", &got, want)
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

func testContext(v *Nvim) func(*testing.T) {
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

func testExtmarks(v *Nvim) func(*testing.T) {
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
		if extmarks[0].ID != gotExtMarkID {
			t.Fatalf("got %d extMarkID but want %d", extmarks[0].ID, extMarkID)
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

func testRuntimeFiles(v *Nvim) func(*testing.T) {
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
			name := name
			tt := tt
			t.Run("Nvim/"+name, func(t *testing.T) {
				t.Parallel()

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
			name := name
			tt := tt
			t.Run("Batch/"+name, func(t *testing.T) {
				t.Parallel()

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

// TODO(zchee): correct testcase
func testTerm(v *Nvim) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Nvim", func(t *testing.T) {
			t.Parallel()

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
			t.Parallel()

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

func TestDial(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported dial unix socket on windows GOOS")
	}

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

// clearBuffer clears the buffer lines.
func clearBuffer(tb testing.TB, v *Nvim, buffer Buffer) {
	tb.Helper()

	if err := v.SetBufferLines(buffer, 0, -1, true, bytes.Fields(nil)); err != nil {
		tb.Fatal(err)
	}
}

func TestLogLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{
			name:  "Trace",
			level: LogTraceLevel,
			want:  "TraceLevel",
		},
		{
			name:  "Debug",
			level: LogDebugLevel,
			want:  "DebugLevel",
		},
		{
			name:  "Info",
			level: LogInfoLevel,
			want:  "InfoLevel",
		},
		{
			name:  "Warn",
			level: LogWarnLevel,
			want:  "WarnLevel",
		},
		{
			name:  "Error",
			level: LogErrorLevel,
			want:  "ErrorLevel",
		},
		{
			name:  "unkonwn",
			level: LogLevel(-1),
			want:  "unkonwn Level",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.level.String(); got != tt.want {
				t.Errorf("LogLevel.String() = %v, want %v", tt.want, got)
			}
		})
	}
}
