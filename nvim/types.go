package nvim

// Mode represents a Nvim's current mode.
type Mode struct {
	// Mode is the current mode.
	Mode string `msgpack:"mode"`

	// Blocking is true if Nvim is waiting for input.
	Blocking bool `msgpack:"blocking"`
}

// HLAttrs represents a highlight definitions.
type HLAttrs struct {
	Bold       bool `msgpack:"bold,omitempty"`
	Underline  bool `msgpack:"underline,omitempty"`
	Undercurl  bool `msgpack:"undercurl,omitempty"`
	Italic     bool `msgpack:"italic,omitempty"`
	Reverse    bool `msgpack:"reverse,omitempty"`
	Foreground int  `msgpack:"foreground,omitempty" empty:"-1"`
	Background int  `msgpack:"background,omitempty" empty:"-1"`
	Special    int  `msgpack:"special,omitempty" empty:"-1"`
}

// Mapping represents a nvim mapping options.
type Mapping struct {
	// LHS is the {lhs} of the mapping.
	LHS string `msgpack:"lhs"`

	// RHS is the {hrs} of the mapping as typed.
	RHS string `msgpack:"rhs"`

	// Silent is 1 for a |:map-silent| mapping, else 0.
	Silent int `msgpack:"silent"`

	// Noremap is 1 if the {rhs} of the mapping is not remappable.
	NoRemap int `msgpack:"noremap"`

	// Expr is  1 for an expression mapping.
	Expr int `msgpack:"expr"`

	// Buffer for a local mapping.
	Buffer int `msgpack:"buffer"`

	// SID is the script local ID, used for <sid> mappings.
	SID int `msgpack:"sid"`

	// Nowait is 1 if map does not wait for other, longer mappings.
	NoWait int `msgpack:"nowait"`

	// Mode specifies modes for which the mapping is defined.
	Mode string `msgpack:"string"`
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
	Max int `msgpack:",array"`
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
// (This could happen if a library first identifies the channel, and a plugin using that library later overrides that info)
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

	// Pty is the name of pseudoterminal, if one is used (optional).
	Pty string `msgpack:"pty,omitempty"`

	// Buffer is the buffer with connected terminal instance (optional).
	Buffer Buffer `msgpack:"buffer,omitempty"`

	// Client is the information about the client on the other end of the RPC channel, if it has added it using nvim_set_client_info (optional).
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
	Bang        bool   `msgpack:"bang"`
	Complete    string `msgpack:"complete,omitempty"`
	Nargs       string `msgpack:"nargs"`
	Range       string `msgpack:"range,omitempty"`
	Name        string `msgpack:"name"`
	ScriptID    int    `msgpack:"script_id"`
	Bar         bool   `msgpack:"bar"`
	Register    bool   `msgpack:"register"`
	Addr        string `msgpack:"addr,omitempty"`
	Count       string `msgpack:"count,omitempty"`
	CompleteArg string `msgpack:"complete_arg,omitempty"`
	Definition  string `msgpack:"definition"`
}

// VirtualTextChunk represents a virtual text chunk.
type VirtualTextChunk struct {
	Text    string `msgpack:",array"`
	HLGroup string `msgpack:",array"`
}
