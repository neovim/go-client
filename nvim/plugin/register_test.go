package plugin_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

func newChildProcess(tb testing.TB) (p *plugin.Plugin, cleanup func()) {
	tb.Helper()

	env := []string{}
	if v := os.Getenv("VIM"); v != "" {
		env = append(env, "VIM="+v)
	}

	ctx := context.Background()
	opts := []nvim.ChildProcessOption{
		nvim.ChildProcessCommand(nvim.BinaryName),
		nvim.ChildProcessArgs("-u", "NONE", "-n", "-i", "NONE", "--embed", "--headless"),
		nvim.ChildProcessContext(ctx),
		nvim.ChildProcessLogf(tb.Logf),
		nvim.ChildProcessEnv(env),
	}
	v, err := nvim.NewChildProcess(opts...)
	if err != nil {
		tb.Fatal(err)
	}

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
		}); walkErr != nil && !os.IsNotExist(err) {
			tb.Fatal(fmt.Errorf("walkErr: %w", errors.Unwrap(walkErr)))
		}
	}

	return plugin.New(v), cleanup
}

func TestRegister(t *testing.T) {
	p, cleanup := newChildProcess(t)
	t.Cleanup(func() {
		cleanup()
	})

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
		expected := "Hello, John Doe"
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
