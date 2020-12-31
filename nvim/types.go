package nvim

// QuickfixError represents an item in a quickfix list.
type QuickfixError struct {
	// Buffer number
	Bufnr int `msgpack:"bufnr,omitempty"`

	// Line number in the file.
	LNum int `msgpack:"lnum,omitempty"`

	// Search pattern used to locate the error.
	Pattern string `msgpack:"pattern,omitempty"`

	// Column number (first column is 1).
	Col int `msgpack:"col,omitempty"`

	// When Vcol is != 0,  Col is visual column.
	VCol int `msgpack:"vcol,omitempty"`

	// Error number.
	Nr int `msgpack:"nr,omitempty"`

	// Description of the error.
	Text string `msgpack:"text,omitempty"`

	// Single-character error type, 'E', 'W', etc.
	Type string `msgpack:"type,omitempty"`

	// Name of a file; only used when bufnr is not present or it is invalid.
	FileName string `msgpack:"filename,omitempty"`

	// Valid is non-zero if this is a recognized error message.
	Valid int `msgpack:"valid,omitempty"`

	// Module name of a module. If given it will be used in quickfix error window instead of the filename.
	Module string `msgpack:"module,omitempty"`
}

// CommandCompletionArgs represents the arguments to a custom command line
// completion function.
//
//  :help :command-completion-custom
type CommandCompletionArgs struct {
	// ArgLead is the leading portion of the argument currently being completed
	// on.
	ArgLead string `msgpack:",array"`

	// CmdLine is the entire command line.
	CmdLine string

	// CursorPosString is decimal representation of the cursor position in
	// bytes.
	CursorPosString int
}

// CursorPos returns the cursor position.
func (a *CommandCompletionArgs) CursorPos() int {
	n := a.CursorPosString
	return n
}

// Mode represents a Nvim's current mode.
type Mode struct {
	// Mode is the current mode.
	Mode string `msgpack:"mode"`

	// Blocking is true if Nvim is waiting for input.
	Blocking bool `msgpack:"blocking"`
}

// HLAttrs represents a highlight definitions.
type HLAttrs struct {
	// Bold is the bold font style.
	Bold bool `msgpack:"bold,omitempty"`

	// Underline is the underline font style.
	Underline bool `msgpack:"underline,omitempty"`

	// Undercurl is the curly underline font style.
	Undercurl bool `msgpack:"undercurl,omitempty"`

	// Italic is the italic font style.
	Italic bool `msgpack:"italic,omitempty"`

	// Reverse is the reverse to foreground and background.
	Reverse bool `msgpack:"reverse,omitempty"`

	// Inverse same as Reverse.
	Inverse bool `msgpack:"inverse,omitempty"`

	// Standout is the standout font style.
	Standout int `msgpack:"standout,omitempty"`

	// Nocombine override attributes instead of combining them.
	Nocombine int `msgpack:"nocombine,omitempty"`

	// Foreground use normal foreground color.
	Foreground int `msgpack:"foreground,omitempty" empty:"-1"`

	// Background use normal background color.
	Background int `msgpack:"background,omitempty" empty:"-1"`

	// Special is used for undercurl and underline.
	Special int `msgpack:"special,omitempty" empty:"-1"`

	// Blend override the blend level for a highlight group within the popupmenu
	// or floating windows.
	//
	// Only takes effect if 'pumblend' or 'winblend' is set for the menu or window.
	// See the help at the respective option.
	Blend int `msgpack:"blend,omitempty"`
}

// Mapping represents a nvim mapping options.
type Mapping struct {
	// LHS is the {lhs} of the mapping.
	LHS string `msgpack:"lhs,omitempty"`

	// RHS is the {hrs} of the mapping as typed.
	RHS string `msgpack:"rhs,omitempty"`

	// Silent is 1 for a :map-silent mapping, else 0.
	Silent int `msgpack:"silent,omitempty"`

	// Noremap is 1 if the {rhs} of the mapping is not remappable.
	NoRemap int `msgpack:"noremap,omitempty"`

	// Expr is  1 for an expression mapping.
	Expr int `msgpack:"expr,omitempty"`

	// Buffer for a local mapping.
	Buffer int `msgpack:"buffer,omitempty"`

	// SID is the script local ID, used for <sid> mappings.
	SID int `msgpack:"sid,omitempty"`

	// Nowait is 1 if map does not wait for other, longer mappings.
	NoWait int `msgpack:"nowait,omitempty"`

	// Mode specifies modes for which the mapping is defined.
	Mode string `msgpack:"string,omitempty"`
}

// ClientVersion represents a version of client for nvim.
type ClientVersion struct {
	// Major major version. (defaults to 0 if not set, for no release yet)
	Major int `msgpack:"major,omitempty" empty:"0"`

	// Minor minor version.
	Minor int `msgpack:"minor,omitempty"`

	// Patch patch number.
	Patch int `msgpack:"patch,omitempty"`

	// Prerelease string describing a prerelease, like "dev" or "beta1".
	Prerelease string `msgpack:"prerelease,omitempty"`

	// Commit hash or similar identifier of commit.
	Commit string `msgpack:"commit,omitempty"`
}

// ClientType type of client information.
type ClientType string

const (
	// RemoteClientType for the client library.
	RemoteClientType ClientType = "remote"

	// UIClientType for the gui frontend.
	UIClientType ClientType = "ui"

	// EmbedderClientType for the application using nvim as a component, for instance IDE/editor implementing a vim mode.
	EmbedderClientType ClientType = "embedder"

	// HostClientType for the plugin host. Typically started by nvim.
	HostClientType ClientType = "host"

	// PluginClientType for the single plugin. Started by nvim.
	PluginClientType ClientType = "plugin"
)

// ClientMethod builtin methods in the client.
//
// For a host, this does not include plugin methods which will be discovered later.
// The key should be the method name, the values are dicts with the following (optional) keys. See below.
//
// Further keys might be added in later versions of nvim and unknown keys are thus ignored.
// Clients must only use keys defined in this or later versions of nvim.
type ClientMethod struct {
	// Async is defines whether the uses notification request or blocking request.
	//
	// If true, send as a notification.
	// If false, send as a blocking request.
	Async bool `msgpack:"async"`

	// NArgs is the number of method arguments.
	NArgs ClientMethodNArgs `msgpack:"nargs"`
}

// ClientMethodNArgs is the number of arguments. Could be a single integer or an array two integers, minimum and maximum inclusive.
type ClientMethodNArgs struct {
	// Min is the minimum number of method arguments.
	Min int `msgpack:",array"`

	// Max is the maximum number of method arguments.
	Max int
}

// ClientAttributes informal attributes describing the client. Clients might define their own keys, but the following are suggested.
type ClientAttributes map[string]string

const (
	// ClientAttributeKeyWebsite Website of client (for instance github repository).
	ClientAttributeKeyWebsite = "website"

	// ClientAttributeKeyLicense Informal description of the license, such as "Apache 2", "GPLv3" or "MIT".
	ClientAttributeKeyLicense = "license"

	// ClientoAttributeKeyLogo URI or path to image, preferably small logo or icon. .png or .svg format is preferred.
	ClientoAttributeKeyLogo = "logo"
)

// Client represents a identify the client for nvim.
//
// Can be called more than once, but subsequent calls will remove earlier info, which should be resent if it is still valid.
// (This could happen if a library first identifies the channel, and a plugin using that library later overrides that info).
type Client struct {
	// Name is short name for the connected client.
	Name string `msgpack:"name,omitempty"`

	// Version describes the version, with the following possible keys (all optional).
	Version ClientVersion `msgpack:"version,omitempty"`

	// Type is the client type. Must be one of the ClientType type values.
	Type ClientType `msgpack:"type,omitempty"`

	// Methods builtin methods in the client.
	Methods map[string]*ClientMethod `msgpack:"methods,omitempty"`

	// Attributes is informal attributes describing the client.
	Attributes ClientAttributes `msgpack:"attributes,omitempty"`
}

// Channel information about a channel.
type Channel struct {
	// Stream is the stream underlying the channel.
	Stream string `msgpack:"stream,omitempty"`

	// Mode is the how data received on the channel is interpreted.
	Mode string `msgpack:"mode,omitempty"`

	// Pty is the name of pseudoterminal, if one is used.
	Pty string `msgpack:"pty,omitempty"`

	// Buffer is the buffer with connected terminal instance.
	Buffer Buffer `msgpack:"buffer,omitempty"`

	// Client is the information about the client on the other end of the RPC channel, if it has added it using SetClientInfo.
	Client *Client `msgpack:"client,omitempty"`
}

// Process represents a Proc and ProcChildren functions return type.
type Process struct {
	// Name is the name of process command.
	Name string `msgpack:"name,omitempty"`

	// PID is the process ID.
	PID int `msgpack:"pid,omitempty"`

	// PPID is the parent process ID.
	PPID int `msgpack:"ppid,omitempty"`
}

// UI represents a nvim ui options.
type UI struct {
	// Height requested height of the UI
	Height int `msgpack:"height,omitempty"`

	// Width requested width of the UI
	Width int `msgpack:"width,omitempty"`

	// RGB whether the UI uses rgb colors (false implies cterm colors)
	RGB bool `msgpack:"rgb,omitempty"`

	// ExtPopupmenu externalize the popupmenu.
	ExtPopupmenu bool `msgpack:"ext_popupmenu,omitempty"`

	// ExtTabline externalize the tabline.
	ExtTabline bool `msgpack:"ext_tabline,omitempty"`

	// ExtCmdline externalize the cmdline.
	ExtCmdline bool `msgpack:"ext_cmdline,omitempty"`

	// ExtWildmenu externalize the wildmenu.
	ExtWildmenu bool `msgpack:"ext_wildmenu,omitempty"`

	// ExtNewgrid use new revision of the grid events.
	ExtNewgrid bool `msgpack:"ext_newgrid,omitempty"`

	// ExtHlstate use detailed highlight state.
	ExtHlstate bool `msgpack:"ext_hlstate,omitempty"`

	// ChannelID channel id of remote UI (not present for TUI)
	ChannelID int `msgpack:"chan,omitempty"`
}

// Command represents a Neovim Ex command.
type Command struct {
	// Name is the name of command.
	Name string `msgpack:"name"`

	// Nargs is the command-nargs.
	// See :help :command-nargs.
	Nargs string `msgpack:"nargs"`

	// Complete is the specifying one or the other of the following attributes.
	// See :help :command-completion.
	Complete string `msgpack:"complete,omitempty"`

	// CompleteArg is the argument completion name.
	CompleteArg string `msgpack:"complete_arg,omitempty"`

	// Range is the specify that the command does take a range, or that it takes an arbitrary count value.
	Range string `msgpack:"range,omitempty"`

	// Count is a count (default N) which is specified either in the line number position, or as an initial argument, like `:Next`.
	// Specifying -count (without a default) acts like -count=0
	Count string `msgpack:"count,omitempty"`

	// Addr is the special characters in the range like `.`, `$` or `%` which by default correspond to the current line,
	// last line and the whole buffer, relate to arguments, (loaded) buffers, windows or tab pages.
	Addr string `msgpack:"addr,omitempty"`

	// Bang is the command can take a ! modifier, like `:q` or `:w`.
	Bang bool `msgpack:"bang"`

	// Bar is the command can be followed by a `|` and another command.
	// A `|` inside the command argument is not allowed then. Also checks for a `"` to start a comment.
	Bar bool `msgpack:"bar"`

	// Register is the first argument to the command can be an optional register name, like `:del`, `:put`, `:yank`.
	Register bool `msgpack:"register"`

	// ScriptID is the line number in the script sid.
	ScriptID int `msgpack:"script_id"`

	// Definition is the command's replacement string.
	Definition string `msgpack:"definition"`
}

// VirtualTextChunk represents a virtual text chunk.
type VirtualTextChunk struct {
	// Text is VirtualText text.
	Text string `msgpack:",array"`

	// HLGroup is VirtualText highlight group.
	HLGroup string
}

// WindowConfig represents a configs of OpenWindow.
//
// Relative is the specifies the type of positioning method used for the floating window.
// The positioning method string keys names:
//
//  editor
// The global editor grid.
//  win
// Window given by the `win` field, or current window by default.
//  cursor
// Cursor position in current window.
//
// Anchor is the decides which corner of the float to place at row and col.
//
//  NW
// northwest (default)
//  NE
// northeast
//  SW
// southwest
//  SE
// southeast
//
// BufPos places float relative to buffer text only when Relative == "win".
// Takes a tuple of zero-indexed [line, column].
//
// Row and Col if given are applied relative to this position, else they default to Row=1 and Col=0 (thus like a tooltip near the buffer text).
//
// Row is the row position in units of "screen cell height", may be fractional.
//
// Col is the column position in units of "screen cell width", may be fractional.
//
// Focusable whether the enable focus by user actions (wincmds, mouse events).
// Defaults to true. Non-focusable windows can be entered by SetCurrentWindow.
//
// External is the GUI should display the window as an external top-level window.
// Currently accepts no other positioning configuration together with this.
//
// Style is the Configure the appearance of the window.
// Currently only takes one non-empty value:
//
//  minimal
// Nvim will display the window with many UI options disabled.
// This is useful when displaying a temporary float where the text should not be edited.
//
// Disables `number`, `relativenumber`, `cursorline`, `cursorcolumn`, `foldcolumn`, `spell` and `list` options.
// And, `signcolumn` is changed to `auto` and `colorcolumn` is cleared.
//
// The end-of-buffer region is hidden by setting `eob` flag of `fillchars` to a space char, and clearing the EndOfBuffer region in `winhighlight`.
type WindowConfig struct {
	// Relative is the specifies the type of positioning method used for the floating window.
	Relative string `msgpack:"relative,omitempty"`

	// Win is the Window for relative="win".
	Win Window `msgpack:"win,omitempty"`

	// Anchor is the decides which corner of the float to place at row and col.
	Anchor string `msgpack:"anchor,omitempty" empty:"NW"`

	// Width is the window width (in character cells). Minimum of 1.
	Width int `msgpack:"width" empty:"1"`

	// Height is the window height (in character cells). Minimum of 1.
	Height int `msgpack:"height" empty:"1"`

	// BufPos places float relative to buffer text only when relative="win".
	BufPos [2]int `msgpack:"bufpos,omitempty"`

	// Row is the row position in units of "screen cell height", may be fractional.
	Row float64 `msgpack:"row,omitempty"`

	// Col is the column position in units of "screen cell width", may be fractional.
	Col float64 `msgpack:"col,omitempty"`

	// Focusable whether the enable focus by user actions (wincmds, mouse events).
	Focusable bool `msgpack:"focusable,omitempty" empty:"true"`

	// External is the GUI should display the window as an external top-level window.
	External bool `msgpack:"external,omitempty"`

	// Style is the Configure the appearance of the window.
	Style string `msgpack:"style,omitempty"`
}

// ExtMarks represents a extmarks type.
type ExtMarks struct {
	// ID is the extmarks ID.
	ID int `msgpack:",array"`

	// Row is the extmark row position.
	Row int

	// Col is the extmark column position.
	Col int
}

// OptionInfo represents a option information.
type OptionInfo struct {
	// Name is the name of the option (like 'filetype').
	Name string `msgpack:"name"`

	// ShortName is the shortened name of the option (like 'ft').
	ShortName string `msgpack:"shortname"`

	// Type is the type of option ("string", "number" or "boolean").
	Type string `msgpack:"type"`

	// Default is the default value for the option.
	Default interface{} `msgpack:"default"`

	// WasSet whether the option was set.
	WasSet bool `msgpack:"was_set"`

	// LastSetSid is the last set script id (if any).
	LastSetSid int `msgpack:"last_set_sid"`

	// LastSetLinenr is the line number where option was set.
	LastSetLinenr int `msgpack:"last_set_linenr"`

	// LastSetChan is the channel where option was set (0 for local).
	LastSetChan int `msgpack:"last_set_chan"`

	// Scope one of "global", "win", or "buf".
	Scope string `msgpack:"scope"`

	// GlobalLocal whether win or buf option has a global value.
	GlobalLocal bool `msgpack:"global_local"`

	// CommaList whether the list of comma separated values.
	CommaList bool `msgpack:"commalist"`

	// FlagList whether the list of single char flags.
	FlagList bool `msgpack:"flaglist"`
}
