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

// +build ignore

// Command apitool generates apiimp.go from apidef.go. The command also has
// an option to compare apidef.go to Nvim's current API meta data.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/neovim/go-client/msgpack"
)

// APIInfo represents the output from nvim --api-info
type APIInfo struct {
	ErrorTypes map[string]ErrorType     `msgpack:"error_types"`
	Types      map[string]ExtensionType `msgpack:"types"`
	Functions  []*Function              `msgpack:"functions"`
}

type ErrorType struct {
	ID int `msgpack:"id"`
}

type ExtensionType struct {
	ID  int    `msgpack:"id"`
	Doc string `msgpack:"-"`
}

type Function struct {
	Name            string   `msgpack:"name"`
	Parameters      []*Field `msgpack:"parameters"`
	ReturnType      string   `msgpack:"return_type"`
	DeprecatedSince int      `msgpack:"deprecated_since"`
	Doc             string   `msgpack:"-"`
	GoName          string   `msgpack:"-"`
	ReturnPtr       bool     `msgpack:"-"`
}

type Field struct {
	Type string `msgpack:",array"`
	Name string
}

var errorTypes = map[string]ErrorType{
	"Exception":  {ID: 0},
	"Validation": {ID: 1},
}

var extensionTypes = map[string]ExtensionType{
	"Buffer":  {ID: 0, Doc: `// Buffer represents a remote Nvim buffer.`},
	"Window":  {ID: 1, Doc: `// Window represents a remote Nvim window.`},
	"Tabpage": {ID: 2, Doc: `// Tabpage represents a remote Nvim tabpage.`},
}

func formatNode(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		panic(err)
	}
	return buf.String()
}

func parseFields(fset *token.FileSet, fl *ast.FieldList) []*Field {
	if fl == nil {
		return nil
	}
	var fields []*Field
	for _, f := range fl.List {
		typ := formatNode(fset, f.Type)
		if len(f.Names) == 0 {
			fields = append(fields, &Field{Type: typ})
		} else {
			for _, id := range f.Names {
				fields = append(fields, &Field{Name: id.Name, Type: typ})
			}
		}
	}
	return fields
}

// parseAPIDef parses the file apidef.go.
func parseAPIDef() ([]*Function, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "apidef.go", nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var functions []*Function
	for _, decl := range file.Decls {
		fdecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		var doc []byte
		if cg := fdecl.Doc; cg != nil {
			for i, c := range cg.List {
				if i > 0 {
					doc = append(doc, '\n')
				}
				doc = append(doc, c.Text...)
			}
		}
		m := &Function{
			GoName:     fdecl.Name.Name,
			Doc:        string(doc),
			Parameters: parseFields(fset, fdecl.Type.Params),
		}
		fields := parseFields(fset, fdecl.Type.Results)
		if len(fields) > 1 {
			return nil, fmt.Errorf("%s: more than one result for %s", fset.Position(fdecl.Pos()), m.Name)
		} else if len(fields) == 1 {
			m.ReturnType = fields[0].Type
		}
		for _, n := range fdecl.Body.List {
			if expr, ok := n.(*ast.ExprStmt); ok {
				if call, ok := expr.X.(*ast.CallExpr); ok {
					if id, ok := call.Fun.(*ast.Ident); ok {
						switch id.Name {
						case "name":
							if len(call.Args) == 1 {
								if id, ok := call.Args[0].(*ast.Ident); ok {
									m.Name = id.Name
								}
							}
						case "deprecatedSince":
							if len(call.Args) == 1 {
								if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
									m.DeprecatedSince, _ = strconv.Atoi(lit.Value)
								}
							}
						case "returnPtr":
							m.ReturnPtr = true
						}
					}
				}
			}
		}
		if m.Name == "" {
			return nil, fmt.Errorf("%s: service method not specified for %s", fset.Position(fdecl.Pos()), m.Name)
		}
		functions = append(functions, m)
	}
	return functions, nil
}

var implementationTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`// Copyright 2016 Gary Burd
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

// Code generated by running "go generate" in github.com/neovim/go-client/nvim. DO NOT EDIT.

package nvim

import (
    "fmt"

    "github.com/neovim/go-client/msgpack"
    "github.com/neovim/go-client/msgpack/rpc"
)

const (
    {{- range $name, $type := .ErrorTypes}}
    {{- lower $name}}Error = {{$type.ID}}
    {{  end -}}
)

func withExtensions() rpc.Option {
	return rpc.WithExtensions(msgpack.ExtensionMap{
{{range $name, $type := .Types}}
		{{$type.ID}}: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return {{$name}}(x), err
		},
{{end}}
	})
}

{{range $name, $type := .Types}}
{{$type.Doc}}
type {{$name}} int

func (x *{{$name}}) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	n, err := unmarshalExt(dec, {{$type.ID}}, x)
	*x = {{$name}}(n)
	return err
}

func (x {{$name}}) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackExtension({{$type.ID}}, encodeExt(int(x)))
}

func (x {{$name}}) String() string {
	return fmt.Sprintf("{{$name}}:%d", int(x))
}
{{end}}

{{range .Functions}}
{{if eq "interface{}" .ReturnType}}
{{.Doc}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} result interface{}) error {
    return v.call("{{.Name}}", result, {{range .Parameters}}{{.Name}},{{end}})
}

{{.Doc}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} result interface{}) {
    b.call("{{.Name}}", result, {{range .Parameters}}{{.Name}},{{end}})
}

{{else if .ReturnType}}
{{.Doc}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) ({{if .ReturnPtr}}*{{end}}{{.ReturnType}}, error) {
    var result {{.ReturnType}}
    err := v.call("{{.Name}}", &result, {{range .Parameters}}{{.Name}},{{end}})
    return {{if .ReturnPtr}}&{{end}}result, err
}
{{.Doc}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} result *{{.ReturnType}}) {
    b.call("{{.Name}}", result, {{range .Parameters}}{{.Name}},{{end}})
}
{{else}}
{{.Doc}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) error {
    return v.call("{{.Name}}", nil, {{range .Parameters}}{{.Name}},{{end}})
}
{{.Doc}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) {
    b.call("{{.Name}}", nil, {{range .Parameters}}{{.Name}},{{end}})
}
{{end}}
{{end}}
`))

func printImplementation(functions []*Function, outFile string) error {
	var buf bytes.Buffer
	if err := implementationTemplate.Execute(&buf, &APIInfo{
		Functions:  functions,
		Types:      extensionTypes,
		ErrorTypes: errorTypes,
	}); err != nil {
		return err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		for i, p := range bytes.Split(buf.Bytes(), []byte("\n")) {
			fmt.Fprintf(os.Stderr, "%d: %s\n", i+1, p)
		}
		return fmt.Errorf("error formating source: %v", err)
	}

	if outFile != "" {
		return ioutil.WriteFile(outFile, out, 0666)
	}
	_, err = os.Stdout.Write(out)
	return err
}

func readAPIInfo() (*APIInfo, error) {
	output, err := exec.Command("nvim", "--api-info").Output()
	if err != nil {
		return nil, err
	}

	var info APIInfo
	if err := msgpack.NewDecoder(bytes.NewReader(output)).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// nvimTypes maps Go types to Nvim API types.
var nvimTypes = map[string]string{
	"":              "void",
	"[]byte":        "String",
	"[]interface{}": "Array",
	"bool":          "Boolean",
	"int":           "Integer",
	"interface{}":   "Object",
	"string":        "String",

	"map[string]interface{}": "Dictionary",

	"[2]int":    "ArrayOf(Integer, 2)",
	"[]Buffer":  "ArrayOf(Buffer)",
	"[]Tabpage": "ArrayOf(Tabpage)",
	"[]Window":  "ArrayOf(Window)",
	"[][]byte":  "ArrayOf(String)",
	"[]string":  "ArrayOf(String)",

	"Mode":       "Dictionary",
	"*HLAttrs":   "Dictionary",
	"[]*Mapping": "ArrayOf(Dictionary)",
}

func convertToNvimTypes(f *Function) *Function {
	if t, ok := nvimTypes[f.ReturnType]; ok {
		f.ReturnType = t
	}
	for _, p := range f.Parameters {
		if t, ok := nvimTypes[p.Type]; ok {
			p.Type = t
		}
	}
	return f
}

type byName []*Function

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

var compareTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`
{{- range .Extra}}< {{template "f" .}}{{end}}
{{- range .Missing}}> {{template "f" .}}{{end}}
{{- range .Different}}----
< {{template "f" index . 0}}> {{template "f" index . 1}}{{end}}
{{- define "f"}}{{.Name}}({{range $i, $p := .Parameters}}{{if $i}}, {{end}}{{$p.Name}} {{$p.Type}}{{end}}){{with .ReturnType}} {{.}}{{end}}
    {{- print " {"}} name({{.Name}}){{with .DeprecatedSince}}; deprecatedSince({{.}});{{end}}{{print " }"}}
{{end}}`))

// specialAPIs lists API calls that are implemented by hand.
var specialAPIs = map[string]bool{
	"nvim_call_atomic":   true,
	"nvim_call_function": true,
	"nvim_execute_lua":   true,
}

func compareFunctions(functions []*Function) error {
	info, err := readAPIInfo()
	if err != nil {
		return err
	}

	sort.Sort(byName(functions))
	sort.Sort(byName(info.Functions))

	var data struct {
		Extra     []*Function
		Missing   []*Function
		Different [][2]*Function
	}

	i := 0
	j := 0
	for i < len(functions) && j < len(info.Functions) {
		a := convertToNvimTypes(functions[i])
		b := info.Functions[j]
		if a.Name < b.Name {
			data.Extra = append(data.Extra, a)
			i++
			continue
		}
		if b.Name < a.Name {
			if b.DeprecatedSince == 0 && !specialAPIs[b.Name] {
				data.Missing = append(data.Missing, b)
			}
			j++
			continue
		}
		equal := len(a.Parameters) == len(b.Parameters) && a.ReturnType == b.ReturnType && a.DeprecatedSince == b.DeprecatedSince
		if equal {
			for i := range a.Parameters {
				if a.Parameters[i].Type != b.Parameters[i].Type {
					equal = false
					break
				}
			}
		}
		if !equal {
			data.Different = append(data.Different, [2]*Function{a, b})
		}
		i++
		j++
	}

	for i < len(functions) {
		a := convertToNvimTypes(functions[i])
		data.Extra = append(data.Extra, a)
		i++
	}

	for j < len(info.Functions) {
		b := info.Functions[j]
		if b.DeprecatedSince == 0 {
			data.Missing = append(data.Missing, b)
		}
		j++
	}

	return compareTemplate.Execute(os.Stdout, &data)
}

func dumpAPI() error {
	output, err := exec.Command("nvim", "--api-info").Output()
	if err != nil {
		return fmt.Errorf("error getting API info: %v", err)
	}

	var v interface{}
	if err := msgpack.NewDecoder(bytes.NewReader(output)).Decode(&v); err != nil {
		return fmt.Errorf("error parsing msppack: %v", err)
	}

	p, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return nil
	}

	os.Stdout.Write(append(p, '\n'))
	return nil
}

func main() {
	log.SetFlags(0)

	generateFlag := flag.String("generate", "", "Generate implementation from apidef.go and write to `file`")
	compareFlag := flag.Bool("compare", false, "Compare apidef.go to the output of nvim --api-info")
	dumpFlag := flag.Bool("dump", false, "Print nvim --api-info as JSON")
	flag.Parse()

	if *dumpFlag {
		if err := dumpAPI(); err != nil {
			log.Fatal(err)
		}
		return
	}

	functions, err := parseAPIDef()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *compareFlag:
		err = compareFunctions(functions)
	default:
		err = printImplementation(functions, *generateFlag)
	}
	if err != nil {
		log.Fatal(err)
	}
}
