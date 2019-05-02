package plugin

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/neovim/go-client/nvim"
)

// Plugin represents a remote plugin.
type Plugin struct {
	Nvim        *nvim.Nvim
	pluginSpecs []*pluginSpec

	// Event/pattern counters used to generate unique paths for autocmds.
	eventPathCounts map[string]int
}

// New returns an intialized plugin.
func New(v *nvim.Nvim) *Plugin {
	p := &Plugin{
		Nvim:            v,
		eventPathCounts: make(map[string]int),
	}

	// Disable support for "specs" method until path mechanism for supporting
	// binary exectables with Nvim is worked out.
	// err := v.RegisterHandler("specs", func(path string) ([]*pluginSpec, error) {
	//  return p.pluginSpecs, nil
	// })

	return p
}

type pluginSpec struct {
	sm   string
	Type string            `msgpack:"type"`
	Name string            `msgpack:"name"`
	Sync bool              `msgpack:"sync"`
	Opts map[string]string `msgpack:"opts"`
}

func (spec *pluginSpec) path() string {
	if i := strings.Index(spec.sm, ":"); i > 0 {
		return spec.sm[:i]
	}
	return ""
}

func isSync(f interface{}) bool {
	t := reflect.TypeOf(f)
	return t.Kind() == reflect.Func && t.NumOut() > 0
}

func (p *Plugin) handle(fn interface{}, spec *pluginSpec) {
	p.pluginSpecs = append(p.pluginSpecs, spec)
	if p.Nvim == nil {
		return
	}
	if err := p.Nvim.RegisterHandler(spec.sm, fn); err != nil {
		panic(err)
	}
}

// Handle registers fn as a MessagePack RPC handler for the specified method
// name. The function signature for fn is one of
//
//  func([v *nvim.Nvim,] {args}) ({resultType}, error)
//  func([v *nvim.Nvim,] {args}) error
//  func([v *nvim.Nvim,] {args})
//
// where {args} is zero or more arguments and {resultType} is the type of of a
// return value. Call the handler from Nvim using the rpcnotify and rpcrequest
// functions:
//
//  :help rpcrequest()
//  :help rpcnotify()
func (p *Plugin) Handle(method string, fn interface{}) {
	if p.Nvim == nil {
		return
	}
	if err := p.Nvim.RegisterHandler(method, fn); err != nil {
		panic(err)
	}
}

// FunctionOptions specifies function options.
type FunctionOptions struct {
	// Name is the name of the function in Nvim. The name must be made of
	// alphanumeric characters and '_', and must start with a capital letter.
	Name string

	// Eval is an expression evaluated in Nvim. The result is passed the
	// handler function.
	Eval string
}

// HandleFunction registers fn as a handler for a Nvim function. The function
// signature for fn is one of
//
//  func([v *nvim.Nvim,] args {arrayType} [, eval {evalType}]) ({resultType}, error)
//  func([v *nvim.Nvim,] args {arrayType} [, eval {evalType}]) error
//
// where {arrayType} is a type that can be unmarshaled from a MessagePack
// array, {evalType} is a type compatible with the Eval option expression and
// {resultType} is the type of function result.
//
// If options.Eval == "*", then HandleFunction constructs the expression to
// evaluate in Nvim from the type of fn's last argument. The last argument is
// assumed to be a pointer to a struct type with 'eval' field tags set to the
// expression to evaluate for each field. Nested structs are supported. The
// expression for the function
//
//  func example(eval *struct{
//      GOPATH string `eval:"$GOPATH"`
//      Cwd    string `eval:"getcwd()"`
//  })
//
// is
//
//  {'GOPATH': $GOPATH, Cwd: getcwd()}
func (p *Plugin) HandleFunction(options *FunctionOptions, fn interface{}) {
	m := make(map[string]string)
	if options.Eval != "" {
		m["eval"] = eval(options.Eval, fn)
	}
	p.handle(fn, &pluginSpec{
		sm:   "0:function:" + options.Name,
		Type: "function",
		Name: options.Name,
		Sync: isSync(fn),
		Opts: m,
	})
}

// CommandOptions specifies command options.
type CommandOptions struct {
	// Name is the name of the command in Nvim. The name must be made of
	// alphanumeric characters and '_', and must start with a capital
	// letter.
	Name string

	// NArgs specifies the number command arguments.
	//
	//  0   No arguments are allowed
	//  1   Exactly one argument is required, it includes spaces
	//  *   Any number of arguments are allowed (0, 1, or many),
	//      separated by white space
	//  ?   0 or 1 arguments are allowed
	//  +   Arguments must be supplied, but any number are allowed
	NArgs string

	// Range specifies that the command accepts a range.
	//
	//  .   Range allowed, default is current line. The value
	//      "." is converted to "" for Nvim.
	//  %   Range allowed, default is whole file (1,$)
	//  N   A count (default N) which is specified in the line
	//      number position (like |:split|); allows for zero line
	//	    number.
	//
	//  :help :command-range
	Range string

	// Count specfies that thecommand accepts a count.
	//
	//  N   A count (default N) which is specified either in the line
	//	    number position, or as an initial argument (like |:Next|).
	//      Specifying -count (without a default) acts like -count=0
	//
	//  :help :command-count
	Count string

	// Addr sepcifies the domain for the range option
	//
	//  lines           Range of lines (this is the default)
	//  arguments       Range for arguments
	//  buffers         Range for buffers (also not loaded buffers)
	//  loaded_buffers  Range for loaded buffers
	//  windows         Range for windows
	//  tabs            Range for tab pages
	//
	//  :help command-addr
	Addr string

	// Bang specifies that the command can take a ! modifier (like :q or :w).
	Bang bool

	// Register specifes that the first argument to the command can be an
	// optional register name (like :del, :put, :yank).
	Register bool

	// Eval is evaluated in Nvim and the result is passed as an argument.
	Eval string

	// Bar specifies that the command can be followed by a "|" and another
	// command.  A "|" inside the command argument is not allowed then. Also
	// checks for a " to start a comment.
	Bar bool

	// Complete specifies command completion.
	//
	//  :help :command-complete
	Complete string
}

// HandleCommand registers fn as a handler for a Nvim command. The arguments
// to the function fn are:
//
//  v *nvim.Nvim        optional
//  args []string       when options.NArgs != ""
//  range [2]int        when options.Range == "." or Range == "%"
//  range int           when options.Range == N or Count != ""
//  bang bool           when options.Bang == true
//  register string     when options.Register == true
//  eval interface{}    when options.Eval != ""
//
// The function fn must return an error.
//
// If options.Eval == "*", then HandleCommand constructs the expression to
// evaluate in Nvim from the type of fn's last argument. See the
// HandleFunction documentation for information on how the expression is
// generated.
func (p *Plugin) HandleCommand(options *CommandOptions, fn interface{}) {
	m := make(map[string]string)

	if options.NArgs != "" {
		m["nargs"] = options.NArgs
	}

	if options.Range != "" {
		if options.Range == "." {
			options.Range = ""
		}
		m["range"] = options.Range
	} else if options.Count != "" {
		m["count"] = options.Count
	}

	if options.Bang {
		m["bang"] = ""
	}

	if options.Register {
		m["register"] = ""
	}

	if options.Eval != "" {
		m["eval"] = eval(options.Eval, fn)
	}

	if options.Addr != "" {
		m["addr"] = options.Addr
	}

	if options.Bar {
		m["bar"] = ""
	}

	if options.Complete != "" {
		m["complete"] = options.Complete
	}

	p.handle(fn, &pluginSpec{
		sm:   "0:command:" + options.Name,
		Type: "command",
		Name: options.Name,
		Sync: isSync(fn),
		Opts: m,
	})
}

// AutocmdOptions specifies autocmd options.
type AutocmdOptions struct {
	// Event is the event name.
	Event string

	// Group specifies the autocmd group.
	Group string

	// Pattern specifies an autocmd pattern.
	//
	//  :help autocmd-patterns
	Pattern string

	// Nested allows nested autocmds.
	//
	//  :help autocmd-nested
	Nested bool

	// Eval is evaluated in Nvim and the result is passed the the handler
	// function.
	Eval string
}

// HandleAutocmd registers fn as a handler an autocmnd event.
//
// If options.Eval == "*", then HandleAutocmd constructs the expression to
// evaluate in Nvim from the type of fn's last argument. See the HandleFunction
// documentation for information on how the expression is generated.
func (p *Plugin) HandleAutocmd(options *AutocmdOptions, fn interface{}) {
	pattern := ""
	m := make(map[string]string)
	if options.Group != "" {
		m["group"] = options.Group
	}
	if options.Pattern != "" {
		m["pattern"] = options.Pattern
		pattern = options.Pattern
	}
	if options.Nested {
		m["nested"] = "1"
	}
	if options.Eval != "" {
		m["eval"] = eval(options.Eval, fn)
	}

	// Compute unique path for event and pattern.
	ep := options.Event + ":" + pattern
	i := p.eventPathCounts[ep]
	p.eventPathCounts[ep] = i + 1
	sm := fmt.Sprintf("%d:autocmd:%s", i, ep)

	p.handle(fn, &pluginSpec{
		sm:   sm,
		Type: "autocmd",
		Name: options.Event,
		Sync: isSync(fn),
		Opts: m,
	})
}

// RegisterForTests registers the plugin with Nvim. Use this method for testing
// plugins in an embedded instance of Nvim.
func (p *Plugin) RegisterForTests() error {
	specs := make(map[string][]*pluginSpec)
	for _, spec := range p.pluginSpecs {
		specs[spec.path()] = append(specs[spec.path()], spec)
	}
	const host = "nvim-go-test"
	for path, specs := range specs {
		if err := p.Nvim.Call("remote#host#RegisterPlugin", nil, host, path, specs); err != nil {
			return err
		}
	}
	err := p.Nvim.Call("remote#host#Register", nil, host, "x", p.Nvim.ChannelID())
	return err
}

func eval(eval string, f interface{}) string {
	if eval != "*" {
		return eval
	}
	ft := reflect.TypeOf(f)
	if ft.Kind() != reflect.Func || ft.NumIn() < 1 {
		panic(`Eval: "*" option requires function with at least one argument`)
	}
	argt := ft.In(ft.NumIn() - 1)
	if argt.Kind() != reflect.Ptr || argt.Elem().Kind() != reflect.Struct {
		panic(`Eval: "*" option requires function with pointer to struct as last argument`)
	}
	return structEval(argt.Elem())
}

func structEval(t reflect.Type) string {
	buf := []byte{'{'}
	sep := ""
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Anonymous {
			panic(`Eval: "*" does not support anonymous fields`)
		}

		eval := sf.Tag.Get("eval")
		if eval == "" {
			ft := sf.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				eval = structEval(ft)
			}
		}

		if eval == "" {
			continue
		}

		name := strings.Split(sf.Tag.Get("msgpack"), ",")[0]
		if name == "" {
			name = sf.Name
		}

		buf = append(buf, sep...)
		buf = append(buf, "'"...)
		buf = append(buf, name...)
		buf = append(buf, "': "...)
		buf = append(buf, eval...)
		sep = ", "
	}
	buf = append(buf, '}')
	return string(buf)
}

type byServiceMethod []*pluginSpec

func (a byServiceMethod) Len() int           { return len(a) }
func (a byServiceMethod) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byServiceMethod) Less(i, j int) bool { return a[i].sm < a[j].sm }

func (p *Plugin) Manifest(host string) []byte {
	var buf bytes.Buffer

	// Sort for consistent order on output.
	sort.Sort(byServiceMethod(p.pluginSpecs))
	escape := strings.NewReplacer("'", "''").Replace

	prevPath := ""
	for _, spec := range p.pluginSpecs {
		path := spec.path()
		if path != prevPath {
			if prevPath != "" {
				fmt.Fprintf(&buf, "\\ )")
			}
			fmt.Fprintf(&buf, "call remote#host#RegisterPlugin('%s', '%s', [\n", host, path)
			prevPath = path
		}

		sync := "0"
		if spec.Sync {
			sync = "1"
		}

		fmt.Fprintf(&buf, "\\ {'type': '%s', 'name': '%s', 'sync': %s, 'opts': {", spec.Type, spec.Name, sync)

		var keys []string
		for k := range spec.Opts {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		optDelim := ""
		for _, k := range keys {
			fmt.Fprintf(&buf, "%s'%s': '%s'", optDelim, k, escape(spec.Opts[k]))
			optDelim = ", "
		}

		fmt.Fprintf(&buf, "}},\n")
	}
	if prevPath != "" {
		fmt.Fprintf(&buf, "\\ ])\n")
	}
	return buf.Bytes()
}
