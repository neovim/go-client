package nvim

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// newChildProcess returns the new *Nvim, and registers cleanup to tb.Cleanup.
func newChildProcess(tb testing.TB, opts ...ChildProcessOption) (v *Nvim) {
	tb.Helper()

	envs := os.Environ()
	envs = append(envs, []string{
		"XDG_CONFIG_HOME=",
		"XDG_DATA_HOME=",
	}...)

	ctx := context.Background()
	copts := []ChildProcessOption{
		ChildProcessCommand(BinaryName),
		ChildProcessArgs(
			"-u", "NONE",
			"-n",
			"-i", "NONE",
			"--headless",
			"--embed",
		),
		ChildProcessContext(ctx),
		ChildProcessLogf(tb.Logf),
		ChildProcessEnv(envs),
	}
	copts = append(copts, opts...)
	n, err := NewChildProcess(copts...)
	if err != nil {
		tb.Fatal(err)
	}
	v = n

	done := make(chan error, 1)
	go func() {
		done <- v.Serve()
	}()

	tb.Cleanup(func() {
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
	})

	if err := v.Command("set packpath="); err != nil {
		tb.Fatal(err)
	}

	return v
}

func TestDial(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported dial unix socket on windows GOOS")
	}

	t.Parallel()

	v1 := newChildProcess(t)
	var addr string
	if err := v1.VVar("servername", &addr); err != nil {
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
		Path: BinaryName,
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

func TestCallWithNoArgs(t *testing.T) {
	t.Parallel()

	v := newChildProcess(t)

	var wd string
	err := v.Call("getcwd", &wd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCallWithNoArgsWithDisabledEmbed(t *testing.T) {
	t.Parallel()

	v := newChildProcess(t, ChildProcessArgs("--embed"), ChildProcessDisableEmbed())

	var wd string
	err := v.Call("getcwd", &wd)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStructValue(t *testing.T) {
	t.Parallel()

	v := newChildProcess(t)

	t.Run("Nvim", func(t *testing.T) {
		var expected, actual struct {
			Str string
			Num int
		}
		expected.Str = `Hello`
		expected.Num = 42
		if err := v.SetVar(`structvar`, &expected); err != nil {
			t.Fatal(err)
		}
		if err := v.Var(`structvar`, &actual); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(&actual, &expected) {
			t.Fatalf("SetVar: got %+v, want %+v", &actual, &expected)
		}
	})

	t.Run("Batch", func(t *testing.T) {
		b := v.NewBatch()

		var expected, actual struct {
			Str string
			Num int
		}
		expected.Str = `Hello`
		expected.Num = 42
		b.SetVar(`structvar`, &expected)
		b.Var(`structvar`, &actual)
		if err := b.Execute(); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(&actual, &expected) {
			t.Fatalf("SetVar: got %+v, want %+v", &actual, &expected)
		}
	})
}
