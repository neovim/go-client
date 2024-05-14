//go:build ignore
// +build ignore

// Command api_tool generates api.go from api_def.go. The command also has
// an option to compare api_def.go to Nvim's current API meta data.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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
	UIOptions  UIOptions                `msgpack:"ui_options"`
	Version    Version                  `msgpack:"version"`
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
	ReturnName      string   `msgpack:"_"`
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

type UIOptions []string

type Version struct {
	APICompatible int  `msgpack:"api_compatible"`
	APILevel      int  `msgpack:"api_level"`
	APIPrerelease bool `msgpack:"api_prerelease"`
	Major         int  `msgpack:"major"`
	Minor         int  `msgpack:"minor"`
	Patch         int  `msgpack:"patch"`
}

var errorTypes = map[string]ErrorType{
	"Exception": {
		ID: 0,
	},
	"Validation": {
		ID: 1,
	},
}

var extensionTypes = map[string]ExtensionType{
	"Buffer": {
		ID:  0,
		Doc: `// Buffer represents a Nvim buffer.`,
	},
	"Window": {
		ID:  1,
		Doc: `// Window represents a Nvim window.`,
	},
	"Tabpage": {
		ID:  2,
		Doc: `// Tabpage represents a Nvim tabpage.`,
	},
}

func formatNode(fset *token.FileSet, node interface{}) string {
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
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

// parseAPIDef parses the file api_def.go.
func parseAPIDef() ([]*Function, []*Function, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "api_def.go", nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var functions []*Function
	var deprecated []*Function

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
			return nil, nil, fmt.Errorf("%s: more than one result for %s", fset.Position(fdecl.Pos()), m.Name)
		}

		if len(fields) == 1 {
			m.ReturnName = fields[0].Name
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
							if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
								m.DeprecatedSince, _ = strconv.Atoi(lit.Value)
							}
						case "returnPtr":
							m.ReturnPtr = true
						}
					}
				}
			}
		}

		if m.Name == "" {
			return nil, nil, fmt.Errorf("%s: service method not specified for %s", fset.Position(fdecl.Pos()), m.Name)
		}

		if m.DeprecatedSince > 0 {
			deprecated = append(deprecated, m)
			continue
		}
		functions = append(functions, m)
	}

	return functions, deprecated, nil
}

const genTemplate = `

{{define "doc" -}}
{{.Doc}}
//
// See: [{{.Name}}()]
// 
// [{{.Name}}()]: https://neovim.io/doc/user/api.html#{{.Name}}()
{{- end}}

{{range .Functions}}
{{if eq "interface{}" .ReturnType}}
{{template "doc" .}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} result interface{}) error {
    return v.call("{{.Name}}", result, {{range .Parameters}}{{.Name}},{{end}})
}

{{template "doc" .}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} result interface{}) {
    b.call("{{.Name}}", &result, {{range .Parameters}}{{.Name}},{{end}})
}

{{else if and .ReturnName .ReturnPtr}}
{{template "doc" .}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) ({{.ReturnName}} *{{.ReturnType}}, err error) {
	var result {{.ReturnType}}
	err = v.call("{{.Name}}", &result, {{range .Parameters}}{{.Name}},{{end}})
	return &result, err
}
{{template "doc" .}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} {{.ReturnName}} *{{.ReturnType}}) {
    b.call("{{.Name}}", {{.ReturnName}}, {{range .Parameters}}{{.Name}},{{end}})
}

{{else if and (.ReturnName) (not .ReturnPtr)}}
{{template "doc" .}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) ({{.ReturnName}} {{.ReturnType}}, err error) {
	err = v.call("{{.Name}}", &{{.ReturnName}}, {{range .Parameters}}{{.Name}},{{end}})
	return {{.ReturnName}}, err
}
{{template "doc" .}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} {{.ReturnName}} *{{.ReturnType}}) {
    b.call("{{.Name}}", {{.ReturnName}}, {{range .Parameters}}{{.Name}},{{end}})
}
{{else if .ReturnType}}
{{template "doc" .}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) ({{if .ReturnPtr}}*{{end}}{{.ReturnType}}, error) {
    var result {{.ReturnType}}
    err := v.call("{{.Name}}", &result, {{range .Parameters}}{{.Name}},{{end}})
    return {{if .ReturnPtr}}&{{end}}result, err
}
{{template "doc" .}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}} result *{{.ReturnType}}) {
    b.call("{{.Name}}", result, {{range .Parameters}}{{.Name}},{{end}})
}
{{else}}
{{template "doc" .}}
func (v *Nvim) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) error {
    return v.call("{{.Name}}", nil, {{range .Parameters}}{{.Name}},{{end}})
}
{{template "doc" .}}
func (b *Batch) {{.GoName}}({{range .Parameters}}{{.Name}} {{.Type}},{{end}}) {
    b.call("{{.Name}}", nil, {{range .Parameters}}{{.Name}},{{end}})
}
{{end}}
{{end}}
`

var implementationTemplate = template.Must(template.New("implementation").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`// Code generated by running "go generate" in github.com/neovim/go-client/nvim. DO NOT EDIT.

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
{{- range $name, $type := .Types}}
		{{$type.ID}}: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return {{$name}}(x), err
		},
{{end -}}
	})
}

{{range $name, $type := .Types}}
{{$type.Doc}}
type {{$name}} int

// MarshalMsgPack implements msgpack.Marshaler.
func (x {{$name}}) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackExtension({{$type.ID}}, encodeExt(int(x)))
}

// UnmarshalMsgPack implements msgpack.Unmarshaler.
func (x *{{$name}}) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	n, err := unmarshalExt(dec, {{$type.ID}}, x)
	*x = {{$name}}(n)
	return err
}

// String returns a string representation of the {{$name}}.
func (x {{$name}}) String() string {
	return fmt.Sprintf("{{$name}}:%d", int(x))
}
{{end}}
` + genTemplate))

var deprecatedTemplate = template.Must(template.New("deprecated").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`// Code generated by running "go generate" in github.com/neovim/go-client/nvim. DO NOT EDIT.

package nvim

// EmbedOptions specifies options for starting an embedded instance of Nvim.
//
// Deprecated: Use ChildProcessOption instead.
type EmbedOptions struct {
	// Logf log function for rpc.WithLogf.
	Logf func(string, ...interface{})

	// Dir specifies the working directory of the command. The working
	// directory in the current process is used if Dir is "".
	Dir string

	// Path is the path of the command to run. If Path = "", then
	// StartEmbeddedNvim searches for "nvim" on $PATH.
	Path string

	// Args specifies the command line arguments. Do not include the program
	// name (the first argument) or the --embed option.
	Args []string

	// Env specifies the environment of the Nvim process. The current process
	// environment is used if Env is nil.
	Env []string
}

// NewEmbedded starts an embedded instance of Nvim using the specified options.
//
// The application must call Serve() to handle RPC requests and responses.
//
// Deprecated: Use NewChildProcess instead.
func NewEmbedded(options *EmbedOptions) (*Nvim, error) {
	if options == nil {
		options = &EmbedOptions{}
	}
	path := options.Path
	if path == "" {
		path = "nvim"
	}

	return NewChildProcess(
		ChildProcessArgs(append([]string{"--embed"}, options.Args...)...),
		ChildProcessCommand(path),
		ChildProcessEnv(options.Env),
		ChildProcessDir(options.Dir),
		ChildProcessServe(false))
}

// ExecuteLua executes a Lua block.
//
// Deprecated: Use ExecLua instead.
func (v *Nvim) ExecuteLua(code string, result interface{}, args ...interface{}) error {
	if args == nil {
		args = emptyArgs
	}
	return v.call("nvim_execute_lua", result, code, args)
}

// ExecuteLua executes a Lua block.
//
// Deprecated: Use ExecLua instead.
func (b *Batch) ExecuteLua(code string, result interface{}, args ...interface{}) {
	if args == nil {
		args = emptyArgs
	}
	b.call("nvim_execute_lua", result, code, args)
}
` + genTemplate))

func printImplementation(functions []*Function, tmpl *template.Template, outFile string) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &APIInfo{
		Functions:  functions,
		Types:      extensionTypes,
		ErrorTypes: errorTypes,
	}); err != nil {
		return fmt.Errorf("falied to Execute implementationTemplate: %w", err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		for i, p := range bytes.Split(buf.Bytes(), []byte("\n")) {
			fmt.Fprintf(os.Stderr, "%d: %s\n", i+1, p)
		}
		return fmt.Errorf("error formating source: %w", err)
	}

	if outFile != "" {
		return os.WriteFile(outFile, out, 0666)
	}
	_, err = os.Stdout.Write(out)
	return err
}

func readAPIInfo(cmdName string) (*APIInfo, error) {
	const cmdArgs = "--api-info"
	output, err := exec.Command(cmdName, cmdArgs).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execuce %s %s: %w", cmdName, cmdArgs, err)
	}

	var info APIInfo
	if err := msgpack.NewDecoder(bytes.NewReader(output)).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode APIInfo: %w", err)
	}
	return &info, nil
}

// nvimTypes maps Go types to Nvim API types.
var nvimTypes = map[string]string{
	"":              "void",
	"[]byte":        "String",
	"[]uint":        "Array",
	"[]interface{}": "Array",
	"bool":          "Boolean",
	"int":           "Integer",
	"interface{}":   "Object",
	"string":        "String",
	"float64":       "Float",

	"ClientType":  "String",
	"Process":     "Object",
	"UserCommand": "Object",

	"Cmd":                         "Dictionary",
	"*Cmd":                        "Dictionary",
	"Channel":                     "Dictionary",
	"*Channel":                    "Dictionary",
	"ClientVersion":               "Dictionary",
	"HLAttrs":                     "Dictionary",
	"*HLAttrs":                    "Dictionary",
	"[]*HLAttrs":                  "Dictionary",
	"WindowConfig":                "Dictionary",
	"*WindowConfig":               "Dictionary",
	"ClientAttributes":            "Dictionary",
	"ClientMethods":               "Dictionary",
	"map[string]*ClientMethod":    "Dictionary",
	"map[string]*Command":         "Dictionary",
	"map[string][]string":         "Dictionary",
	"map[string]bool":             "Dictionary",
	"map[string]int":              "Dictionary",
	"map[string]interface{}":      "Dictionary",
	"map[string]OptionValueScope": "Dictionary",
	"Mode":                        "Dictionary",
	"OptionInfo":                  "Dictionary",

	"[]*AutocmdType": "Array",
	"[]*Channel":     "Array",
	"[]*Process":     "Array",
	"[]*UI":          "Array",
	"[]ExtMark":      "Array",
	"[]TextChunk":    "Array",
	"Mark":           "Array",

	"[2]int":     "ArrayOf(Integer, 2)",
	"[]*Mapping": "ArrayOf(Dictionary)",
	"[][]byte":   "ArrayOf(String)",
	"[]Buffer":   "ArrayOf(Buffer)",
	"[]int":      "ArrayOf(Integer)",
	"[]string":   "ArrayOf(String)",
	"[]Tabpage":  "ArrayOf(Tabpage)",
	"[]Window":   "ArrayOf(Window)",
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

// hiddenAPIs list of hidden API.
var hiddenAPIs = map[string]bool{
	"nvim__set_hl_ns": true,
}

// specialAPIs lists API calls that are implemented by hand.
var specialAPIs = map[string]bool{
	"nvim_call_atomic":             true,
	"nvim_call_function":           true,
	"nvim_call_dict_function":      true,
	"nvim_execute_lua":             true,
	"nvim_exec_lua":                true,
	"nvim_buf_call":                true,
	"nvim_set_decoration_provider": true,
	"nvim_chan_send":               true, // FUNC_API_LUA_ONLY
	"nvim_win_call":                true, // FUNC_API_LUA_ONLY
	"nvim_notify":                  true, // implements underling nlua(vim.notify)
	"nvim_get_option_info":         true, // deprecated
}

func compareFunctions(cmdName string, functions []*Function) error {
	info, err := readAPIInfo(cmdName)
	if err != nil {
		return fmt.Errorf("failed to real APIInfo :%w", err)
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
			if !hiddenAPIs[a.Name] {
				data.Extra = append(data.Extra, a)
			}
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

	if err := compareTemplate.Execute(os.Stdout, &data); err != nil {
		return fmt.Errorf("falied to Execute compareTemplate: %w", err)
	}
	return nil
}

func dumpAPI(cmdName string) error {
	output, err := exec.Command(cmdName, "--api-info").Output()
	if err != nil {
		return fmt.Errorf("error getting API info: %w", err)
	}

	var v interface{}
	if err := msgpack.NewDecoder(bytes.NewReader(output)).Decode(&v); err != nil {
		return fmt.Errorf("error parsing msppack: %w", err)
	}

	p, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return nil
	}

	os.Stdout.Write(append(p, '\n'))
	return nil
}

var (
	flagNvim       string
	flagGenerate   string
	flagDeprecated string
	flagCompare    bool
	flagDump       bool
)

func main() {
	log.SetFlags(log.Lshortfile)

	flag.StringVar(&flagNvim, "nvim", "nvim", "nvim binary path")
	flag.StringVar(&flagGenerate, "generate", "", "Generate implementation from api_def.go and write to `file`")
	flag.StringVar(&flagDeprecated, "deprecated", "", "Generate deprecated implementation from api_def.go and write to `file`")
	flag.BoolVar(&flagCompare, "compare", false, "Compare api_def.go to the output of nvim --api-info")
	flag.BoolVar(&flagDump, "dump", false, "Print nvim --api-info as JSON")
	flag.Parse()

	if flagDump {
		if err := dumpAPI(flagNvim); err != nil {
			log.Fatal(err)
		}
		return
	}

	functions, deprecated, err := parseAPIDef()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case flagCompare:
		functions = append(functions, deprecated...)
		if err := compareFunctions(flagNvim, functions); err != nil {
			log.Fatal(err)
		}

	case flagGenerate != "":
		if flagDeprecated == "" {
			functions = append(functions, deprecated...)
		}
		if err := printImplementation(functions, implementationTemplate, flagGenerate); err != nil {
			log.Fatal(err)
		}

		if flagDeprecated != "" {
			if err := printImplementation(deprecated, deprecatedTemplate, flagDeprecated); err != nil {
				log.Fatal(err)
			}
		}
	}
}
