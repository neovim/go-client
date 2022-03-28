package plugin_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/nvimtest"
	"github.com/neovim/go-client/nvim/plugin"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	type testEval struct {
		BaseDir string `eval:"fnamemodify(getcwd(), ':t')"`
	}

	p := plugin.New(nvimtest.NewChildProcess(t))

	// SimpleHandler
	p.Handle("hello",
		func(s string) (string, error) {
			return "Hello, " + s, nil
		},
	)

	// FunctionHandler
	p.HandleFunction(
		&plugin.FunctionOptions{Name: "Hello"},
		func(args []string) (string, error) {
			return "Hello, " + strings.Join(args, " "), nil
		},
	)

	// FunctionEvalHandler
	p.HandleFunction(
		&plugin.FunctionOptions{Name: "TestEval", Eval: "*"},
		func(_ []string, eval *testEval) (string, error) {
			return eval.BaseDir, nil
		},
	)

	// CommandHandler
	p.HandleCommand(
		&plugin.CommandOptions{
			Name:     "Hello",
			NArgs:    "*",
			Range:    "%",
			Addr:     "buffers",
			Complete: "buffer",
			Bang:     true,
			Register: true,
			Bar:      true,
		},
		func(n *nvim.Nvim, args []string) error {
			chunks := []nvim.TextChunk{
				{Text: `Hello`},
			}
			for _, arg := range args {
				chunks = append(chunks, nvim.TextChunk{Text: arg})
			}

			return n.Echo(chunks, true, make(map[string]interface{}))
		},
	)

	// CommandRangeDotHandler
	p.HandleCommand(
		&plugin.CommandOptions{
			Name:  "HelloRangeDot",
			Range: ".",
		},
		func(n *nvim.Nvim) error {
			chunks := []nvim.TextChunk{
				{Text: `Hello`},
			}

			return n.Echo(chunks, true, make(map[string]interface{}))
		},
	)

	// CommandCountHandler
	p.HandleCommand(
		&plugin.CommandOptions{
			Name:  "HelloCount",
			Count: "0",
		},
		func(n *nvim.Nvim) error {
			chunks := []nvim.TextChunk{
				{Text: `Hello`},
			}

			return n.Echo(chunks, true, make(map[string]interface{}))
		},
	)

	// CommandEvalHandler
	p.HandleCommand(
		&plugin.CommandOptions{
			Name: "HelloEval",
			Eval: "*",
		},
		func(n *nvim.Nvim, eval *testEval) error {
			chunks := []nvim.TextChunk{
				{
					Text: eval.BaseDir,
				},
			}

			return n.Echo(chunks, true, make(map[string]interface{}))
		},
	)

	if err := p.RegisterForTests(); err != nil {
		t.Fatalf("register handlers for test: %v", err)
	}

	t.Run("SimpleHandler", func(t *testing.T) {
		result, err := p.Nvim.Exec(`:echo Hello('John', 'Doe')`, true)
		if err != nil {
			t.Fatalf("exec 'echo Hello' function: %v", err)
		}

		expected := `Hello, John Doe`
		if result != expected {
			t.Fatalf("Hello returned %q, want %q", result, expected)
		}
	})

	t.Run("FunctionHandler", func(t *testing.T) {
		cid := p.Nvim.ChannelID()
		var result string
		if err := p.Nvim.Call(`rpcrequest`, &result, cid, `hello`, `world`); err != nil {
			t.Fatalf("call rpcrequest(%v, %v, %v, %v): %v", &result, cid, "hello", "world", err)
		}

		expected2 := `Hello, world`
		if result != expected2 {
			t.Fatalf("hello returned %q, want %q", result, expected2)
		}
	})

	t.Run("FunctionEvalHandler", func(t *testing.T) {
		var result string
		if err := p.Nvim.Eval(`TestEval()`, &result); err != nil {
			t.Fatalf("eval 'TestEval()' function: %v", err)
		}

		expected3 := `plugin`
		if result != expected3 {
			t.Fatalf("EvalTest returned %q, want %q", result, expected3)
		}
	})

	t.Run("CommandHandler", func(t *testing.T) {
		result, err := p.Nvim.Exec(`Hello World`, true)
		if err != nil {
			t.Fatalf("exec 'Hello' command: %v", err)
		}

		expected := `Helloorld`
		if result != expected {
			t.Fatalf("Hello returned %q, want %q", result, expected)
		}
	})

	t.Run("CommandRangeDotHandler", func(t *testing.T) {
		result, err := p.Nvim.Exec(`HelloRangeDot`, true)
		if err != nil {
			t.Fatalf("exec 'Hello' command: %v", err)
		}

		expected := `Hello`
		if result != expected {
			t.Fatalf("Hello returned %q, want %q", result, expected)
		}
	})

	t.Run("CommandCountHandler", func(t *testing.T) {
		result, err := p.Nvim.Exec(`HelloCount`, true)
		if err != nil {
			t.Fatalf("exec 'Hello' command: %v", err)
		}

		expected := `Hello`
		if result != expected {
			t.Fatalf("Hello returned %q, want %q", result, expected)
		}
	})

	t.Run("CommandEvalHandler", func(t *testing.T) {
		result, err := p.Nvim.Exec(`HelloEval`, true)
		if err != nil {
			t.Fatalf("exec 'Hello' command: %v", err)
		}

		expected := `plugin`
		if result != expected {
			t.Fatalf("Hello returned %q, want %q", result, expected)
		}
	})
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	p := plugin.New(nvimtest.NewChildProcess(t))

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
	p.Handle(event1, func() error { return eventFn1(t, p.Nvim) })

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
	p.Handle(event2, func() error { return eventFn2(t, p.Nvim) })

	if err := p.RegisterForTests(); err != nil {
		t.Fatalf("register for test: %v", err)
	}

	if err := p.Nvim.Subscribe(event1); err != nil {
		t.Fatalf("subscribe(%v): %v", event1, err)
	}

	b := p.Nvim.NewBatch()
	b.Subscribe(event2)
	if err := b.Execute(); err != nil {
		t.Fatalf("batch execute: %v", err)
	}

	// warm-up
	var result int
	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q)`, event1), &result); err != nil {
		t.Fatalf("eval rpcnotify for warm-up of event1: %v", err)
	}
	if result != 1 {
		t.Fatalf("expect 1 but got %d", result)
	}

	var result2 int
	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 1, 2, 3)`, event1), &result2); err != nil {
		t.Fatalf("eval rpcnotify for event1: %v", err)
	}
	if result2 != 1 {
		t.Fatalf("expect 1 but got %d", result2)
	}

	var result3 int
	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 4, 5, 6)`, event2), &result3); err != nil {
		t.Fatalf("eval rpcnotify for event2: %v", err)
	}
	if result3 != 1 {
		t.Fatalf("expect 1 but got %d", result3)
	}

	if err := p.Nvim.Unsubscribe(event1); err != nil {
		t.Fatalf("unsubscribe event1: %v", err)
	}

	b.Unsubscribe(event2)
	if err := b.Execute(); err != nil {
		t.Fatalf("unsubscribe event2: %v", err)
	}

	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 7, 8, 9)`, event1), nil); err != nil {
		t.Fatalf("ensure rpcnotify to event1 is no-op: %v", err)
	}

	if err := p.Nvim.Eval(fmt.Sprintf(`rpcnotify(0, %q, 10, 11, 12)`, event2), nil); err != nil {
		t.Fatalf("ensure rpcnotify to event2 is no-op: %v", err)
	}
}
