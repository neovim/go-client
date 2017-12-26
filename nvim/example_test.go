package nvim_test

import (
	"fmt"
	"log"
	"os"

	"github.com/neovim/go-client/nvim"
)

// This program lists the names of the Nvim buffers when run from an Nvim
// terminal. It dials to Nvim using the $NVIM_LISTEN_ADDRESS and fetches all of
// the buffer names in one call using a batch.
func Example() {
	// Get address from environment variable set by Nvim.
	addr := os.Getenv("NVIM_LISTEN_ADDRESS")
	if addr == "" {
		log.Fatal("NVIM_LISTEN_ADDRESS not set")
	}

	// Dial with default options.
	v, err := nvim.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup on return.
	defer v.Close()

	bufs, err := v.Buffers()
	if err != nil {
		log.Fatal(err)
	}

	// Get the names using a single atomic call to Nvim.
	names := make([]string, len(bufs))
	b := v.NewBatch()
	for i, buf := range bufs {
		b.BufferName(buf, &names[i])
	}
	if err := b.Execute(); err != nil {
		log.Fatal(err)
	}

	// Print the names.
	for _, name := range names {
		fmt.Println(name)
	}
}
