// Copyright 2016 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package plugin is a Nvim remote plugin host.
package plugin

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/neovim/go-client/nvim"
)

// Main implements the main function for a Nvim remote plugin.
//
// Plugin applications call the Main function to run the plugin. The Main
// function creates a Nvim client, calls the supplied function to register
// handlers with the plugin and then runs the server loop to handle requests
// from Nvim.
//
// Plugin applications should use the default logger in the standard log
// package for logging. If the environment variable NVIM_GO_LOG_FILE is set,
// then the default logger is configured to append to the file specified by the
// environment variable.
//
// Run the plugin application with the command line option --manifest=hostName
// to print the plugin manifest to stdout. Add the manifest manually to a
// Vimscript file. The :UpdateRemotePlugins command is not supported at this
// time.
//
// If the --manifest=host command line flag is specified, then Main prints the
// plugin manifest to stdout insead of running the application as a plugin.
func Main(registerHandlers func(p *Plugin) error) {
	pluginHost := flag.String("manifest", "", "Write plugin manifest for `host` to stdout")
	flag.Parse()

	if *pluginHost != "" {
		log.SetFlags(0)
		p := New(nil)
		if err := registerHandlers(p); err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(p.Manifest(*pluginHost))
		return
	}

	stdout := os.Stdout
	if fname := os.Getenv("NVIM_GO_LOG_FILE"); fname != "" {
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		os.Stdout = f
		os.Stderr = f
		log.SetOutput(f)
		log.SetPrefix(fmt.Sprintf("%8d ", os.Getpid()))
		log.Print("Plugin Start")
		defer log.Print("Plugin Exit")
	} else {
		log.SetFlags(0)
		os.Stdout = os.Stderr
	}
	v, err := nvim.New(os.Stdin, stdout, stdout, log.Printf)
	if err != nil {
		log.Fatal(err)
	}
	p := New(v)
	if err := registerHandlers(p); err != nil {
		log.Fatal(err)
	}
	if err := v.Serve(); err != nil {
		log.Fatal(err)
	}
}
