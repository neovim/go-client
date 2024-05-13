This this example Neovim plugin shows how to invoke a [Go](https://go.dev/)
function from a plugin.

The plugin starts a Go program containing the function as a child process. The
plugin invokes functions in the child process using
[RPC](https://neovim.io/doc/user/api.html#RPC).

Use the following steps to run the plugin:

1. Build the program with the [go tool](https://golang.org/cmd/go/) to an
   executable named `helloremote`. Ensure that the executable is in a directory in
   the `PATH` environment variable.
   ```
   $ cd helloremote
   $ go build
   ```
1. Install the plugin in this directory using a plugin manager or by adding
   this directory to the
   [runtimepath](https://neovim.io/doc/user/options.html#'runtimepath').
1. Start Nvim and run the following command:
   ```vim
   :Hello world!
   ```
