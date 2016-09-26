// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin_test

import (
	"os"
	"strings"
	"testing"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func newEmbeddedPlugin(t *testing.T) (*plugin.Plugin, func()) {
	env := []string{}
	if v := os.Getenv("VIM"); v != "" {
		env = append(env, "VIM="+v)
	}
	v, err := nvim.NewEmbedded(&nvim.EmbedOptions{
		Args: []string{"-u", "NONE", "-n"},
		Env:  env,
		Logf: t.Logf,
	})
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- v.Serve()
	}()

	return plugin.New(v), func() {
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

func TestRegister(t *testing.T) {
	p, cleanup := newEmbeddedPlugin(t)
	defer cleanup()

	p.Handle("hello", func(s string) (string, error) {
		return "Hello, " + s, nil
	})
	p.HandleFunction(&plugin.FunctionOptions{Name: "Hello"}, func(args []string) (string, error) {
		return "Hello, " + strings.Join(args, " "), nil
	})

	if err := p.RegisterForTests(); err != nil {
		t.Fatal(err)
	}

	{
		result, err := p.Nvim.CommandOutput(":echo Hello('John', 'Doe')")
		if err != nil {
			t.Error(err)
		}
		expected := "\nHello, John Doe"
		if result != expected {
			t.Errorf("Hello returned %q, want %q", result, expected)
		}
	}

	{
		cid := p.Nvim.ChannelID()

		var result string
		if err := p.Nvim.Call("rpcrequest", &result, cid, "hello", "world"); err != nil {
			t.Fatal(err)
		}

		expected := "Hello, world"
		if result != expected {
			t.Errorf("hello returned %q, want %q", result, expected)
		}
	}
}
