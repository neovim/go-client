// Command gen generates apiimp.go from api.mpack.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/neovim/go-client/msgpack"
)

var (
	flagGenerate   string
	flagCompare    bool
	flagDump       bool
	flagNvimCommit string
)

func init() {
	spew.Config = spew.ConfigState{
		Indent:           " ",
		SortKeys:         false,
		ContinueOnMethod: true,
	}

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage of gen [api.mpack]:")
		flag.PrintDefaults()
	}

	flag.StringVar(&flagGenerate, "generate", "", "Generate implementation from apidef.go and write to `file`")
	flag.BoolVar(&flagCompare, "compare", false, "Compare apidef.go to the output of api.mpack")
	flag.BoolVar(&flagDump, "dump", false, "Print api.mpack as a JSON")
	flag.StringVar(&flagNvimCommit, "commit", "master", "Full commit hash of neovim repository for generate api.mpack")
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("genapi: ")
	flag.Parse()

	var mpack io.ReadCloser
	var err error
	if flag.NArg() == 0 {
		log.Println("gen api.mpack ...")
		r, err := genAPIMpack()
		if err != nil {
			log.Fatal(err)
		}
		mpack = ioutil.NopCloser(r)
	} else {
		mpack, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer mpack.Close()
	}

	if flagDump {
		if err := dumpJSON(mpack); err != nil {
			log.Fatal(err)
		}
	}
}

const neovimRepoURL = "https://github.com/neovim/neovim"

func genAPIMpack() (io.Reader, error) {
	tmpdir, err := ioutil.TempDir("", "go-client")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	git := &exec.Cmd{
		Path: gitPath,
		Args: []string{"git", "init"},
		Dir:  tmpdir,
	}
	if err := git.Run(); err != nil {
		return nil, fmt.Errorf("failed to git: %v", err)
	}

	log.Println("git remote add origin ...")
	git = &exec.Cmd{
		Path: gitPath,
		Args: []string{"git", "remote", "add", "origin", neovimRepoURL},
		Dir:  tmpdir,
	}
	if err := git.Run(); err != nil {
		return nil, fmt.Errorf("failed to git: %v", err)
	}

	log.Printf("git fetch --depth=1 origin %s ...\n", flagNvimCommit)
	git = &exec.Cmd{
		Path: gitPath,
		Args: []string{"git", "fetch", "--depth=1", "origin", flagNvimCommit},
		Dir:  tmpdir,
	}
	if err := git.Run(); err != nil {
		return nil, fmt.Errorf("failed to git: %v", err)
	}

	if flagNvimCommit == "master" {
		flagNvimCommit = path.Join("origin", "master")
	}
	log.Printf("git reset --hard %s ...\n", flagNvimCommit)
	git = &exec.Cmd{
		Path: gitPath,
		Args: []string{"git", "reset", "--hard", flagNvimCommit},
		Dir:  tmpdir,
	}
	if err := git.Run(); err != nil {
		return nil, fmt.Errorf("failed to git: %v", err)
	}

	python3, err := exec.LookPath("python3")
	if err != nil {
		return nil, err
	}

	genVimdoc := exec.Command(python3, []string{filepath.Join(tmpdir, "scripts", "gen_vimdoc.py")}...)
	genVimdoc.Dir = tmpdir
	genVimdoc.Env = append(os.Environ(), []string{"INCLUDE_C_DECL=true", "INCLUDE_DEPRECATED=true"}...)

	log.Println("exec scripts/gen_vimdoc.py ...")
	_ = genVimdoc.Run() // ignore python error

	data, err := ioutil.ReadFile(filepath.Join(tmpdir, "runtime", "doc", "api.mpack"))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = buf.Write(data)

	return &buf, err
}

func dumpJSON(r io.Reader) error {
	var api API
	if err := msgpack.NewDecoder(r).Decode(&api); err != nil {
		return fmt.Errorf("failed to parsing msppack: %w", err)
	}

	data, err := json.MarshalIndent(api, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal api to json: %w", err)
	}

	_, err = os.Stdout.Write(append(data, '\n'))
	return err
}
