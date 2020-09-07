package plugin_test

import (
	"fmt"
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
	v, err := nvim.NewChildProcess(
		nvim.ChildProcessArgs("-u", "NONE", "-n", "--embed"),
		nvim.ChildProcessEnv(env),
		nvim.ChildProcessLogf(t.Logf))
	if err != nil {
		t.Fatal(err)
	}

	return plugin.New(v), func() {
		if err := v.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

type testEval struct {
	Buffer int64 `eval:"nvim_get_current_buf()"`
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
	p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufEnter", Pattern: "*", Eval: "*"},
		func(eval *testEval) error {
			fmt.Printf("eval: %#v\n", eval)
			// if err := p.Nvim.SetVar("fname", eval.Buffer); err != nil {
			// 	return err
			// }
			return nil
		})

	err := p.RegisterForTests()
	if err != nil {
		t.Fatal(err)
	}

	result, err := p.Nvim.CommandOutput(":echo Hello('John', 'Doe')")
	if err != nil {
		t.Error(err)
	}
	expected := "Hello, John Doe"
	if result != expected {
		t.Errorf("Hello returned %q, want %q", result, expected)
	}

	cid := p.Nvim.ChannelID()

	var result2 string
	if err := p.Nvim.Call("rpcrequest", &result2, cid, "hello", "world"); err != nil {
		t.Fatal(err)
	}

	expected2 := "Hello, world"
	if result2 != expected2 {
		t.Errorf("hello returned %q, want %q", result2, expected2)
	}

	// result3, err := p.Nvim.CommandOutput(":echomsg g:fname")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// var result3 string
	// if err := p.Nvim.Var("fname", &result3); err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("result2: %s", result3)
}
