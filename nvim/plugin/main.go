// Package plugin is a Nvim remote plugin host.
package plugin

import (
	"bytes"
	"flag"
	"log"
	"os"
	"regexp"

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
	vimFilePath := flag.String("location", "", "Manifest is automatically written to `.vim file`")
	flag.Parse()

	if *pluginHost != "" {
		log.SetFlags(0)
		p := New(nil)
		if err := registerHandlers(p); err != nil {
			log.Fatal(err)
		}
		manifest := p.Manifest(*pluginHost)
		if *vimFilePath != "" {
			if err := overwriteManifest(*vimFilePath, *pluginHost, manifest); err != nil {
				log.Fatal(err)
			}
		} else {
			os.Stdout.Write(manifest)
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

func overwriteManifest(path, host string, manifest []byte) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	output := replaceManifest(host, input, manifest)
	return os.WriteFile(path, output, 0666)
}

func replaceManifest(host string, input, manifest []byte) []byte {
	p := regexp.MustCompile(`(?ms)^call remote#host#RegisterPlugin\('` + regexp.QuoteMeta(host) + `'.*?^\\ ]\)$`)
	match := p.FindIndex(input)
	var output []byte
	if match == nil {
		if len(input) > 0 && input[len(input)-1] != '\n' {
			input = append(input, '\n')
		}
		output = append(input, manifest...)
	} else {
		if match[1] != len(input) {
			// No need for trailing \n if in middle of file.
			manifest = bytes.TrimSuffix(manifest, []byte{'\n'})
		}
		output = append([]byte{}, input[:match[0]]...)
		output = append(output, manifest...)
		output = append(output, input[match[1]:]...)
	}
	return output
}
