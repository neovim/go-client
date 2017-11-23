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
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
)

// Main implements the main function for a Nvim remote plugin.
//
// Plugin applications call the Main function to run the plugin. The Main
// function creates a Nvim client, calls the supplied function to register
// handlers with the plugin and then runs the server loop to handle requests
// from Nvim.
//
// Applications should use the default logger in the standard log package to
// write to Nvim's log.
//
// Run the plugin application with the command line option --manifest=hostName
// to print the plugin manifest to stdout. Add the manifest manually to a
// Vimscript file. The :UpdateRemotePlugins command is not supported at this
// time.
//
// If the --manifest=host command line flag is specified, then Main prints the
// plugin manifest to stdout insead of running the application as a plugin.
// If the --location=vimfile command line flag is specified, then plugin
// manifest will be automatically written to .vim file.
func Main(registerHandlers func(p *Plugin) error) {
	pluginHost := flag.String("manifest", "", "Write plugin manifest for `host` to stdout")
	vimFilePath := flag.String("location", "", "if this option is enable, manifest is automatically written `.vim file`")
	flag.Parse()

	if *pluginHost != "" {
		log.SetFlags(0)
		p := New(nil)
		if err := registerHandlers(p); err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write(p.Manifest(*pluginHost))
	}

	if *vimFilePath != "" && *pluginHost != "" {
		log.SetFlags(0)
		p := New(nil)
		if err := registerHandlers(p); err != nil {
			log.Fatal(err)
		}
		if err := replaceManifest(*vimFilePath, p.Manifest(*pluginHost)); err != nil {
			log.Fatal(err)
		}
		return
	}

	stdout := os.Stdout
	os.Stdout = os.Stderr
	log.SetFlags(0)

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

func replaceManifest(path string, newManifest []byte) error {
	lines := strings.Split(string(newManifest), "\n")
	if len(lines) == 0 {
		return errors.New("no manifest")
	}

	head := lines[0]

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("no such file: %s", path)
	}

	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	flg := true
	temporaryScript := make([]string, 0)
	for scanner.Scan() {
		if scanner.Text() == head {
			flg = false
		}
		if scanner.Text() == "\\ ])" {
			flg = true
		}
		if flg {
			if scanner.Text() != "\\ ])" {
				temporaryScript = append(temporaryScript, scanner.Text())
			}
		}
	}

	newManifestLines := strings.Split(string(newManifest), "\n")
	temporaryScript = append(temporaryScript, newManifestLines...)
	var script []byte
	for i, t := range temporaryScript {
		script = append(script, []byte(t)...)
		if i != len(temporaryScript)-1 {
			script = append(script, []byte("\n")...)
		}
	}

	if err := ioutil.WriteFile(path, []byte(script), 0666); err != nil {
		return err
	}

	return nil
}
