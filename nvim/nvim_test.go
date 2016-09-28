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

	// Simple handler.
	{
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
	}

	// Buffers
	{
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

	// Windows
	{
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

	// Tabpage
	{
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

	// Lines
	{
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

	// Vars
	{
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

	// Set variable to struct value
	{
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

	// Eval
	{
		var a, b string
		if err := v.Eval(`["hello", "world"]`, []*string{&a, &b}); err != nil {
			t.Error(err)
		}
		if a != "hello" || b != "world" {
			t.Errorf("a=%q b=%q, want a=hello b=world", a, b)
		}
	}

	// Pipeline
	{
		p := v.NewPipeline()
		results := make([]int, 128)

		for i := range results {
			p.SetVar(fmt.Sprintf("v%d", i), i)
		}

		for i := range results {
			p.Var(fmt.Sprintf("v%d", i), &results[i])
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
		p.Var("v3", &i)
		if err := p.Wait(); err != nil {
			log.Fatal(err)
		}
		if i != 3 {
			t.Fatalf("i = %d, want %d", i, 3)
		}
	}

	// Call with no args.
	{
		var wd string
		err := v.Call("getcwd", &wd)
		if err != nil {
			t.Fatal(err)
		}
	}

}
