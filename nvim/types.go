package nvim

const (
	// EventBufChangedtick event name of "nvim_buf_changedtick_event".
	EventBufChangedtick = "nvim_buf_changedtick_event"

	// EventBufLines event name of "nvim_buf_lines_event".
	EventBufLines = "nvim_buf_lines_event"

	// EventBufDetach event name of "nvim_buf_detach_event".
	EventBufDetach = "nvim_buf_detach_event"
)

// ChangedtickEvent represents a EventBufChangedtick type.
type ChangedtickEvent struct {
	Buffer     Buffer `msgpack:"buffer,omitempty"`
	Changetick int64  `msgpack:"changetick,omitempty"`
}

// BufLinesEvent represents a EventBufLines type.
type BufLinesEvent struct {
	Buffer      Buffer   `msgpack:"buffer,omitempty"`
	Changetick  int64    `msgpack:"changetick,omitempty"`
	FirstLine   int64    `msgpack:"firstLine,omitempty"`
	LastLine    int64    `msgpack:"lastLine,omitempty"`
	LineData    []string `msgpack:",array"`
	IsMultipart bool     `msgpack:"isMultipart,omitempty"`
}

// BufDetachEvent represents a EventBufDetach type.
type BufDetachEvent struct {
	Buffer Buffer `msgpack:"buffer,omitempty"`
}

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
//	:help :command-completion-custom
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
	return a.CursorPosString
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

	// Standout is the standout font style.
	Standout int `msgpack:"standout,omitempty"`

	// Underline is the underline font style.
	Underline bool `msgpack:"underline,omitempty"`

	// Undercurl is the curly underline font style.
	Undercurl bool `msgpack:"undercurl,omitempty"`

	// Italic is the italic font style.
	Italic bool `msgpack:"italic,omitempty"`

	// Reverse is the reverse to foreground and background.
	Reverse bool `msgpack:"reverse,omitempty"`

	// Strikethrough is the strikethrough font style.
	Strikethrough bool `msgpack:"strikethrough,omitempty"`

	ForegroundIndexed bool `msgpack:"fg_indexed,omitempty"`

	BackgroundIndexed bool `msgpack:"bg_indexed,omitempty"`

	// Foreground is foreground color of RGB color.
	Foreground int `msgpack:"foreground,omitempty" empty:"-1"`

	// Background is background color of RGB color.
	Background int `msgpack:"background,omitempty" empty:"-1"`

	// Special is used for undercurl and underline.
	Special int `msgpack:"special,omitempty" empty:"-1"`

	// Blend override the blend level for a highlight group within the popupmenu
	// or floating windows.
	//
	// Only takes effect if 'pumblend' or 'winblend' is set for the menu or window.
	// See the help at the respective option.
	Blend int `msgpack:"blend,omitempty"`

	// Nocombine override attributes instead of combining them.
	Nocombine bool `msgpack:"nocombine,omitempty"`

	// Default don't override existing definition, like "hi default".
	//
	// This value is used only SetHighlight.
	Default bool `msgpack:"default,omitempty"`

	// Cterm is cterm attribute map. Sets attributed for cterm colors.
	//
	// Note thet by default cterm attributes are same as attributes of gui color.
	//
	// This value is used only SetHighlight.
	Cterm *HLAttrs `msgpack:"cterm,omitempty"`

	// CtermForeground is the foreground of cterm color.
	//
	// This value is used only SetHighlight.
	CtermForeground int `msgpack:"ctermfg,omitempty" empty:"-1"`

	// CtermBackground is the background of cterm color.
	//
	// This value is used only SetHighlight.
	CtermBackground int `msgpack:"ctermbg,omitempty" empty:"-1"`
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
	NArgs ClientMethodNArgs
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

// Cmd represents a Neovim Ex Cmd args.
type Cmd struct {
	// Cmd command name.
	Cmd string `msgpack:"cmd"`

	// Range command <range>.
	//
	// Can have 0-2 elements depending on how many items the range contains.
	//
	// Has no elements if command doesn't accept a range or if no range was specified,
	// one element if only a single range item was specified and two elements if both range items were specified.
	Range []int `msgpack:"range,omitempty"`

	// Count any <count> that was supplied to the command. -1 if command cannot take a count.
	Count int `msgpack:"count,omitempty"`

	// Reg is the optional command |<register>|, if specified. Empty string if not specified or if command cannot take a register.
	Reg string `msgpack:"reg,omitempty"`

	// Bang Whether command contains a |<bang>| (!) modifier.
	Bang bool `msgpack:"bang,omitempty"`

	// Args Command arguments.
	Args []string `msgpack:"args,omitempty"`

	// Magic which characters have special meaning in the command arguments.
	Magic *CmdMagic `msgpack:"magic,omitempty"`

	Mods *CmdMods `msgpack:"mods,omitempty"`

	// Nargs value of command-nargs.
	Nargs string `msgpack:"nargs,omitempty" empty:"*"`

	// Addr value of command-addr. Uses short name.
	Addr string `msgpack:"addr,omitempty" empty:"none"`

	// Nextcmd next command if there are multiple commands separated by a |:bar|. Empty if there isn't a next command.
	Nextcmd string `msgpack:"nextcmd,omitempty"`
}

// CmdMagic which characters have special meaning in the command arguments.
type CmdMagic struct {
	// File is the command expands filenames. Which means characters such as "%", "#" and wildcards are expanded.
	File bool `msgpack:"file,omitempty"`
	// Bar is the "|" character is treated as a command separator and the double quote character (\") is treated as the start of a comment.
	Bar bool `msgpack:"bar,omitempty"`
}

// CmdMods command-modifiers.
type CmdMods struct {
	// Silent is the |:help :silent|
	Silent bool `msgpack:"silent,omitempty"`
	// emsg_silent is the |:help :silent|
	EmsgSilent bool `msgpack:"emsg_silent,omitempty"`
	// Unsilent is the |:help :unsilent|
	Unsilent bool `msgpack:"unsilent,omitempty"`
	// Filter is filter.
	Filter *CmdModsFilter `msgpack:"filter,omitempty"`
	// Sandbox is the |:help :sandbox|
	Sandbox bool `msgpack:"sandbox,omitempty"`
	// Noautocmd is the |:help :noautocmd|
	Noautocmd bool `msgpack:"noautocmd,omitempty"`
	// Browse is the |:help :browse|
	Browse bool `msgpack:"browse,omitempty"`
	// Confirm is the |:help :confirm|
	Confirm bool `msgpack:"confirm,omitempty"`
	// Hide is the |:help :hide|
	Hide bool `msgpack:"hide,omitempty"`
	// Horizontal is the |:help :horizontal|
	Horizontal bool `msgpack:"horizontal,omitempty"`
	// Keepalt is the |:help :keepalt|
	Keepalt bool `msgpack:"keepalt,omitempty"`
	// Keepjumps is the |:help :keepjumps|
	Keepjumps bool `msgpack:"keepjumps,omitempty"`
	// Keepmarks is the |:help :keepmarks|
	Keepmarks bool `msgpack:"keepmarks,omitempty"`
	// Keeppatterns is the |:help :keeppatterns|
	Keeppatterns bool `msgpack:"keeppatterns,omitempty"`
	// Lockmarks is the |:help :lockmarks|
	Lockmarks bool `msgpack:"lockmarks,omitempty"`
	// Noswapfile is the |:help :noswapfile|
	Noswapfile bool `msgpack:"noswapfile,omitempty"`
	// Tab is number of Tab. -1 when omitted.
	Tab int `msgpack:"tab,omitempty" empty:"-1"`
	// Verbose is number of verbose. -1 when omitted.
	Verbose int `msgpack:"verbose,omitempty" empty:"-1"`
	// Vertical is the |:help :vertical|
	Vertical bool `msgpack:"vertical,omitempty"`
	// Split split modifier string, is an empty string when there's no split modifier. If there is a split modifier it can be one of:
	//	aboveleft
	//	belowright
	//	topleft"
	//	botright"
	Split string `msgpack:"split,omitempty"`
}

type CmdModsFilter struct {
	// Pattern filter pattern. Empty string if there is no filter.
	Pattern string `msgpack:"pattern,omitempty"`
	// Force whether filter is inverted or not.
	Force bool `msgpack:"force,omitempty"`
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

// UserCommand represesents a user command.
type UserCommand interface {
	command()
}

// UserVimCommand is a user Vim command executed at UserCommand.
type UserVimCommand string

// make sure UserVimCommand implements the UserCommand interface.
var _ UserCommand = (*UserVimCommand)(nil)

// command implements UserCommand.command.
func (UserVimCommand) command() {}

// UserLuaCommand is a user Lua command executed at UserCommand.
type UserLuaCommand struct {
	// Args passed to the command, if any.
	Args string `msgpack:"args,omitempty"`

	// Bang true if the command was executed with a ! modifier.
	Bang bool `msgpack:"bang"`

	// StartLine is the starting line of the command range.
	StartLine int `msgpack:"line1,omitempty"`

	// FinalLine is the final line of the command range.
	FinalLine int `msgpack:"line2,omitempty"`

	// Range is the number of items in the command range: 0, 1, or 2.
	Range int `msgpack:"range,omitempty"`

	// Count is the any count supplied.
	Count int `msgpack:"count,omitempty"`

	// Reg is the optional register, if specified.
	Reg string `msgpack:"reg,omitempty"`

	// Mode is the command modifiers, if any.
	Mode string `msgpack:"mode,omitempty"`
}

// make sure UserLuaCommand implements the UserCommand interface.
var _ UserCommand = (*UserLuaCommand)(nil)

// command implements UserCommand.command.
func (UserLuaCommand) command() {}

// TextChunk represents a text chunk.
type TextChunk struct {
	// Text is text.
	Text string `msgpack:",array"`

	// HLGroup is text highlight group.
	HLGroup string
}

// WindowConfig represents a configs of OpenWindow.
//
// Relative is the specifies the type of positioning method used for the floating window.
// The positioning method string keys names:
//
//	editor
//
// The global editor grid.
//
//	win
//
// Window given by the `win` field, or current window by default.
//
//	cursor
//
// Cursor position in current window.
//
// Win is window ID for Relative="win".
//
// Anchor is the decides which corner of the float to place at row and col.
//
//	NW
//
// northwest (default)
//
//	NE
//
// northeast
//
//	SW
//
// southwest
//
//	SE
//
// southeast
//
// BufPos places float relative to buffer text only when Relative == "win".
// Takes a tuple of zero-indexed [line, column].
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
// ZIndex is stacking order. floats with higher "zindex" go on top on floats with lower indices. Must be larger than zero.
// The default value for floats are 50. In general, values below 100 are recommended, unless there is a good reason to overshadow builtin elements.
//
// Style is the Configure the appearance of the window.
// Currently only takes one non-empty value:
//
//	minimal
//
// Nvim will display the window with many UI options disabled.
// This is useful when displaying a temporary float where the text should not be edited.
//
// Disables "number", "relativenumber", "cursorline", "cursorcolumn", "foldcolumn", "spell" and "list" options.
// And, "signcolumn" is changed to "auto" and "colorcolumn" is cleared.
// The end-of-buffer region is hidden by setting "eob" flag of "fillchars" to a space char, and clearing the EndOfBuffer region in "winhighlight".
//
//	border
//
// Style of (optional) window border. This can either be a string or an array.
// The string values are:
//
//	none
//
// No border. This is the default.
//
//	single
//
// A single line box.
//
//	double
//
// A double line box.
//
//	rounded
//
// Like "single", but with rounded corners ("╭" etc.).
//
//	solid
//
// Adds padding by a single whitespace cell.
//
//	shadow
//
// A drop shadow effect by blending with the background.
//
// If it is an array it should be an array of eight items or any divisor of
// eight. The array will specifify the eight chars building up the border
// in a clockwise fashion starting with the top-left corner.
// As, an example, the double box style could be specified as:
//
//	[ "╔", "═" ,"╗", "║", "╝", "═", "╚", "║" ]
//
// If the number of chars are less than eight, they will be repeated.
// Thus an ASCII border could be specified as:
//
//	[ "/", "-", "\\", "|" ]
//
// Or all chars the same as:
//
//	[ "x" ]
//
// An empty string can be used to turn off a specific border, for instance,
//
//	[ "", "", "", ">", "", "", "", "<" ]
//
// By default "FloatBorder" highlight is used which links to "VertSplit"
// when not defined.
// It could also be specified by character:
//
//	[ {"+", "MyCorner"}, {"x", "MyBorder"} ]
//
// NoAutocmd is if true then no buffer-related autocommand events such as BufEnter, BufLeave or BufWinEnter may fire from calling this function.
type WindowConfig struct {
	// Relative is the specifies the type of positioning method used for the floating window.
	Relative string `msgpack:"relative,omitempty"`

	// Win is the Window for relative="win".
	Win Window `msgpack:"win,omitempty"`

	// Anchor is the decides which corner of the float to place at row and col.
	Anchor string `msgpack:"anchor,omitempty"`

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

	// ZIndex stacking order. floats with higher `zindex` go on top on floats with lower indices. Must be larger than zero.
	ZIndex int `msgpack:"zindex,omitempty" empty:"50"`

	// Style is the Configure the appearance of the window.
	Style string `msgpack:"style,omitempty"`

	// Border is the style of window border.
	Border any `msgpack:"border,omitempty"`

	// NoAutocmd whether the fire buffer-related autocommand events
	NoAutocmd bool `msgpack:"noautocmd,omitempty"`
}

// BorderStyle represents a WindowConfig.Border style.
type BorderStyle string

// list of BorderStyle.
const (
	// BorderStyleNone is the no border. This is the default.
	BorderStyleNone = BorderStyle("none")
	// BorderStyleSingle is a single line box.
	BorderStyleSingle = BorderStyle("single")
	// BorderStyleDouble a double line box.
	BorderStyleDouble = BorderStyle("double")
	// BorderStyleRounded like "single", but with rounded corners ("╭" etc.).
	BorderStyleRounded = BorderStyle("rounded")
	// BorderStyleSolid adds padding by a single whitespace cell.
	BorderStyleSolid = BorderStyle("solid")
	// BorderStyleShadow a drop shadow effect by blending with the background.
	BorderStyleShadow = BorderStyle("shadow")
)

// ExtMark represents a extmarks type.
type ExtMark struct {
	// ID is the extmarks ID.
	ID int `msgpack:",array"`

	// Row is the extmark row position.
	Row int

	// Col is the extmark column position.
	Col int
}

// Mark represents a mark.
type Mark struct {
	Row        int `msgpack:",array"`
	Col        int
	Buffer     Buffer
	BufferName string
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
	Default any `msgpack:"default"`

	// Scope one of "global", "win", or "buf".
	Scope string `msgpack:"scope"`

	// LastSetSid is the last set script id (if any).
	LastSetSid int `msgpack:"last_set_sid"`

	// LastSetLinenr is the line number where option was set.
	LastSetLinenr int `msgpack:"last_set_linenr"`

	// LastSetChan is the channel where option was set (0 for local).
	LastSetChan int `msgpack:"last_set_chan"`

	// WasSet whether the option was set.
	WasSet bool `msgpack:"was_set"`

	// GlobalLocal whether win or buf option has a global value.
	GlobalLocal bool `msgpack:"global_local"`

	// CommaList whether the list of comma separated values.
	CommaList bool `msgpack:"commalist"`

	// FlagList whether the list of single char flags.
	FlagList bool `msgpack:"flaglist"`
}

// OptionValueScope represents a OptionValue scope optional parameter value.
type OptionValueScope string

// list of OptionValueScope.
const (
	GlobalScope = OptionValueScope("global")
	LocalScope  = OptionValueScope("local")
)

// LogLevel represents a nvim log level.
type LogLevel int

// list of LogLevels.
//
// Should kept sync neovim LogLevel.
const (
	LogTraceLevel LogLevel = iota
	LogDebugLevel
	LogInfoLevel
	LogWarnLevel
	LogErrorLevel
)

// String returns a string representation of the LogLevel.
func (level LogLevel) String() string {
	switch level {
	case LogTraceLevel:
		return "TraceLevel"
	case LogDebugLevel:
		return "DebugLevel"
	case LogInfoLevel:
		return "InfoLevel"
	case LogWarnLevel:
		return "WarnLevel"
	case LogErrorLevel:
		return "ErrorLevel"
	default:
		return "unknown Level"
	}
}

// AutocmdArg is the type of autocmd API arguments.
type AutocmdArg struct {
	Group   int `msgpack:"group,omitempty"`
	Event   any `msgpack:"event,omitempty"`
	Pattern any `msgpack:"pattern,omitempty"`
}

// AutocmdType is the type of autocmd API return type.
type AutocmdType struct {
	// ID is the autocommand id (only when defined with the API).
	ID int `msgpack:"int"`

	// Group is the autocommand group id.
	Group int `msgpack:"group"`

	// Desc is the autocommand description.
	Desc string `msgpack:"desc"`

	// Event is the autocommand event.
	Event string `msgpack:"event"`

	// Command is the autocommand command.
	Command string `msgpack:"command"`

	// Once whether the autocommand is only run once.
	Once bool `msgpack:"once"`

	// Pattern is the autocommand pattern.
	Pattern string `msgpack:"pattern"`

	// BufLocal is the true if the autocommand is buffer local.
	BufLocal bool `msgpack:"buflocal,omitempty"`

	// Buffer is the buffer number.
	Buffer Buffer `msgpack:"buffer,omitempty"`
}
