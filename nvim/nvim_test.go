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
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

func newEmbeddedNvim(t *testing.T) (*Nvim, func()) {
	v, err := NewEmbedded(&EmbedOptions{
		Args: []string{"-u", "NONE", "-n"},
		Env:  []string{},
		Logf: t.Logf,
	})
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- v.Serve()
	}()

	return v, func() {
		e1 := v.Close()
		e2 := <-done
		if e1 != nil {
			t.Fatal(e1)
		}
		if e2 != nil {
			t.Fatal(e2)
		}
	}
}

func helloHandler(s string) (string, error) {
	return "Hello, " + s, nil
}

func TestAPI(t *testing.T) {
	v, cleanup := newEmbeddedNvim(t)
	defer cleanup()
	cid := v.ChannelID()
	if cid <= 0 {
		t.Fatal("could not get channel id")
	}

	t.Run("simpleHandler", func(t *testing.T) {
		if err := v.RegisterHandler("hello", helloHandler); err != nil {
			t.Fatal(err)
		}
		var result string
		if err := v.Call("rpcrequest", &result, cid, "hello", "world"); err != nil {
			t.Fatal(err)
		}
		expected := "Hello, world"
		if result != expected {
			t.Errorf("hello returned %q, want %q", result, expected)
		}
	})

	t.Run("buffer", func(t *testing.T) {
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
	})

	t.Run("window", func(t *testing.T) {
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
	})

	t.Run("tabpage", func(t *testing.T) {
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
	})

	t.Run("lines", func(t *testing.T) {
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
	})

	t.Run("var", func(t *testing.T) {
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
	})

	t.Run("structValue", func(t *testing.T) {
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
	})

	t.Run("eval", func(t *testing.T) {
		var a, b string
		if err := v.Eval(`["hello", "world"]`, []*string{&a, &b}); err != nil {
			t.Error(err)
		}
		if a != "hello" || b != "world" {
			t.Errorf("a=%q b=%q, want a=hello b=world", a, b)
		}
	})

	t.Run("batch", func(t *testing.T) {
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

	})

	t.Run("pipeline", func(t *testing.T) {
		p := v.NewPipeline()
		results := make([]int, 128)

		for i := range results {
			p.SetVar(fmt.Sprintf("batch%d", i), i)
		}

		for i := range results {
			p.Var(fmt.Sprintf("batch%d", i), &results[i])
		}

		if err := p.Wait(); err != nil {
			t.Fatal(err)
		}

		for i := range results {
			if results[i] != i {
				t.Fatalf("result = %d, want %d", results[i], i)
			}
		}

		// Reuse pipeline

		var i int
		p.Var("batch3", &i)
		if err := p.Wait(); err != nil {
			log.Fatal(err)
		}
		if i != 3 {
			t.Fatalf("i = %d, want %d", i, 3)
		}
	})

	t.Run("callWithNoArgs", func(t *testing.T) {
		var wd string
		err := v.Call("getcwd", &wd)
		if err != nil {
			t.Fatal(err)
		}
	})

}
