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

func TestSubscribe(t *testing.T) {
	p, cleanup := newChildProcess(t)
	t.Cleanup(func() {
		cleanup()
	})

	const event1 = "event1"
	eventFn1 := func(t *testing.T, v *nvim.Nvim) error {
		return v.RegisterHandler(event1, func(event ...interface{}) {
			if event[0] != int64(1) {
				t.Fatalf("expected event[0] is 1 but got %d", event[0])
			}
			if event[1] != int64(2) {
				t.Fatalf("expected event[1] is 2 but got %d", event[1])
			}
			if event[2] != int64(3) {
				t.Fatalf("expected event[2] is 3 but got %d", event[2])
			}
		})
	}
	const event2 = "event2"
	eventFn2 := func(t *testing.T, v *nvim.Nvim) error {
		return v.RegisterHandler(event1, func(event ...interface{}) {
			if event[0] != int64(4) {
				t.Fatalf("expected event[0] is 4 but got %d", event[0])
			}
			if event[1] != int64(5) {
				t.Fatalf("expected event[1] is 5 but got %d", event[1])
			}
			if event[2] != int64(6) {
				t.Fatalf("expected event[2] is 6 but got %d", event[2])
			}
		})
	}
	p.Handle(event1, func() error { return eventFn1(t, p.Nvim) })
	p.Handle(event2, func() error { return eventFn2(t, p.Nvim) })

	if err := p.RegisterForTests(); err != nil {
		t.Fatal(err)
	}

	if err := p.Nvim.Subscribe(event1); err != nil {
		t.Fatal(err)
	}
	b := p.Nvim.NewBatch()
	b.Subscribe(event2)
	if err := b.Execute(); err != nil {
		t.Fatal(err)
	}

	// warm-up
	var result int
	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q)`, event1), &result); err != nil {
		t.Fatal(err)
	}
	if result != 1 {
		t.Fatalf("expect 1 but got %d", result)
	}

	var result2 int
	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 1, 2, 3)`, event1), &result2); err != nil {
		t.Fatal(err)
	}
	if result2 != 1 {
		t.Fatalf("expect 1 but got %d", result2)
	}

	var result3 int
	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 4, 5, 6)`, event2), &result3); err != nil {
		t.Fatal(err)
	}
	if result3 != 1 {
		t.Fatalf("expect 1 but got %d", result3)
	}

	if err := p.Nvim.Unsubscribe(event1); err != nil {
		t.Fatal(err)
	}
	b.Unsubscribe(event2)
	if err := b.Execute(); err != nil {
		t.Fatal(err)
	}

	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 7, 8, 9)`, event1), nil); err != nil {
		t.Fatal(err)
	}

	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 10, 11, 12)`, event2), nil); err != nil {
		t.Fatal(err)
	}
}
