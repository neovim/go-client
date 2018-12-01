// Copyright 2018 Gary Burd
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

package nvim

type Mode struct {
	// Mode is the current mode.
	Mode string `msgpack:"mode"`

	// Blocking is true if Nvim is waiting for input.
	Blocking bool `msgpack:"blocking"`
}

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

type VersionType string

const (
	// VersionMajor major version. (defaults to 0 if not set, for no release yet)
	VersionMajor VersionType = "major"
	// VersionMinor minor version.
	VersionMinor VersionType = "minor"
	// VersionPatch patch number.
	VersionPatch VersionType = "patch"
	// VersionPrerelease string describing a prerelease, like "dev" or "beta1".
	VersionPrerelease VersionType = "prerelease"
	// VersionCommit hash or similar identifier of commit.
	VersionCommit VersionType = "commit"
)

// Version type of describing the version, with the following.
type Version map[VersionType]string

// ClientType type of client type.
type ClientType string

const (
	RemoteClient   ClientType = "remote"
	UIClient       ClientType = "ui"
	EmbedderClient ClientType = "embedder"
	HostClient     ClientType = "host"
	PluginClient   ClientType = "plugin"
)

type MethodsType string

const (
	// Async is if true, send as a notification. If false or unspecified, use a blocking request.
	MethodsAsync MethodsType = "async"
	// Nargs is the number of arguments. Could be a single integer or an array two integers, minimum and maximum inclusive.
	MethodsNargs MethodsType = "nargs"
)

// Methods builtin methods in the client.
//
// For a host, this does not include plugin methods which will be discovered later.
// The key should be the method name, the values are dicts with the following (optional) keys. See below.
//
// Further keys might be added in later versions of nvim and unknown keys are thus ignored.
// Clients must only use keys defined in this or later versions of nvim.
type Methods map[MethodsType]string

type AttributesType string

const (
	// AttributesWebsite Website of client (for instance github repository)
	AttributesWebsite AttributesType = "website"
	// AttributesLicense Informal description of the license, such as "Apache 2", "GPLv3" or "MIT"
	AttributesLicense AttributesType = "license"
	// AttributesLogo URI or path to image, preferably small logo or icon. .png or .svg format is preferred.
	AttributesLogo AttributesType = "logo"
)

// Attributes informal attributes describing the client. Clients might define their own keys, but the following are suggested.
type Attributes map[AttributesType]string

// Client represents a identify the client for nvim.
//
// Can be called more than once, but subsequent calls will remove earlier info, which should be resent if it is still valid.
// (This could happen if a library first identifies the channel, and a plugin using that library later overrides that info)
type Client struct {
	// Name is short name for the connected client.
	Name string `msgpack:"name,omitempty"`
	// Version describes the version, with the following possible keys (all optional).
	Version Version `msgpack:"version,omitempty"`
	// Type is the client type. Must be one of the ClientType type values.
	Type ClientType `msgpack:"type,omitempty"`
	// Methods builtin methods in the client.
	Methods Methods `msgpack:"methods,omitempty"`
	// Attributes is informal attributes describing the client.
	Attributes Attributes `msgpack:"attributes,omitempty"`
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

type Process struct {
	// Name is the name of process command.
	Name string `msgpack:"name,omitempty"`
	// PID is the process ID.
	PID int `msgpack:"pid,omitempty"`
	// PPID is the parent process ID.
	PPID int `msgpack:"ppid,omitempty"`
}

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
