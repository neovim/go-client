package nvimtest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/neovim/go-client/nvim"
)

// NewChildProcess returns the new *Nvim, and registers cleanup to tb.Cleanup.
func NewChildProcess(tb testing.TB) *nvim.Nvim {
	tb.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	tb.Cleanup(func() {
		cancel()
	})

	tmpdir := tb.TempDir()
	envs := os.Environ()
	envs = append(envs, []string{
		fmt.Sprintf("XDG_CONFIG_HOME=%s", tmpdir),
		fmt.Sprintf("XDG_DATA_HOME=%s", tmpdir),
		fmt.Sprintf("NVIM_LOG_FILE=%s", os.DevNull),
	}...)
	opts := []nvim.ChildProcessOption{
		nvim.ChildProcessCommand(nvim.BinaryName),
		nvim.ChildProcessArgs(
			"--clean",             // Mimics a fresh install of Nvim. See :help --clean
			"--embed",             // Use stdin/stdout as a msgpack-RPC channel, so applications can embed and control Nvim via the RPC API.
			"--headless",          // Start without UI, and do not wait for nvim_ui_attach
			"-c", "set packpath=", // Clean packpath
		),
		nvim.ChildProcessContext(ctx),
		nvim.ChildProcessEnv(envs),
		nvim.ChildProcessServe(true),
		nvim.ChildProcessLogf(tb.Logf),
	}
	n, err := nvim.NewChildProcess(opts...)
	if err != nil {
		tb.Fatal(err)
	}

	tb.Cleanup(func() {
		if err := n.Close(); err != nil {
			tb.Fatal(err)
		}
	})

	return n
}
