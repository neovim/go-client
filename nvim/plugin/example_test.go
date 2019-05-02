package plugin_test

import (
	"strings"

	"github.com/neovim/go-client/nvim/plugin"
)

// This plugin adds the Hello function to Nvim.
func Example() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleFunction(&plugin.FunctionOptions{Name: "Hello"}, func(args []string) (string, error) {
			return "Hello, " + strings.Join(args, " "), nil
		})
		return nil
	})
}
