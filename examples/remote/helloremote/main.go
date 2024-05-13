package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func hello(v *nvim.Nvim, args []string) error {
	return v.WriteOut(fmt.Sprintf("Hello %s\n", strings.Join(args, " ")))
}

func main() {
	// Turn off timestamps in output.
	log.SetFlags(0)

	// Direct writes by the application to stdout garble the RPC stream.
	// Redirect the application's direct use of stdout to stderr.
	stdout := os.Stdout
	os.Stdout = os.Stderr

	// Create a client connected to stdio. Configure the client to use the
	// standard log package for logging.
	v, err := nvim.New(os.Stdin, stdout, stdout, log.Printf)
	if err != nil {
		log.Fatal(err)
	}

	// Register function with the client.
	v.RegisterHandler("hello", hello)

	// Run the RPC message loop. The Serve function returns when
	// nvim closes.
	if err := v.Serve(); err != nil {
		log.Fatal(err)
	}
}
