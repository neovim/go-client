//go:build ignore
// +build ignore

// This file defines the Nvim remote API using Go syntax. Run the 'go generate'
// command to convert this file to the API implementation in api.go.
//
// The go generate command runs the apitool program. Run
//
//  go run api_tool.go --help
//
// to learn more about the apitool program.

package main

// autocmd.c

// Autocmds get all autocommands that match the corresponding {opts}.
//
// Note that when multiple patterns or events are provided, it will find all the autocommands that
// match any combination of them.
func Autocmds(opts map[string]interface{}) (result []*AutocmdType) {
	name(nvim_get_autocmds)
}

// CreateAutocmd create an autocommand.
//
// The API allows for two (mutually exclusive) types of actions to be executed when the autocommand
// triggers: a callback function (Lua or Vimscript), or a command (like regular autocommands).
func CreateAutocmd(event interface{}, opts map[string]interface{}) (id int) {
	name(nvim_create_autocmd)
}

// DeleteAutocmd delete an autocommand by id.
//
// NOTE: Only autocommands created via the API have an id.
func DeleteAutocmd(id int) {
	name(nvim_del_autocmd)
}

// ClearAutocmds clear all autocommands that match the corresponding {opts}.
//
// To delete a particular autocmd, see DeleteAutocmd.
func ClearAutocmds(opts map[string]interface{}) {
	name(nvim_clear_autocmds)
}

// CreateAugroup create or get an autocommand group(autocmd-groups).
func CreateAugroup(name string, opts map[string]interface{}) (id int) {
	name(nvim_create_augroup)
}

// DeleteAugroupByID delete an autocommand group by id.
func DeleteAugroupByID(id int) {
	name(nvim_del_augroup_by_id)
}

// DeleteAugroupByID delete an autocommand group by name.
func DeleteAugroupByName(name string) {
	name(nvim_del_augroup_by_name)
}

// ExecAutocmds execute all autocommands for {event} that match the corresponding {opts} autocmd-execute.
func ExecAutocmds(event interface{}, opts map[string]interface{}) {
	name(nvim_exec_autocmds)
}

// buffer.c

// BufferLineCount gets the buffer line count.
//
// The buffer arg is specific Buffer, or 0 for current buffer.
//
// The returns line count, or 0 for unloaded buffer.
func BufferLineCount(buffer Buffer) (count int) {
	name(nvim_buf_line_count)
}

// AttachBuffer activates buffer-update events on a channel.
//
// The buffer is specific Buffer, or 0 for current buffer.
//
// If sendBuffer is true, initial notification should contain the whole buffer.
// If false, the first notification will be a "nvim_buf_lines_event".
// Otherwise, the first notification will be a "nvim_buf_changedtick_event".
//
// Returns whether the updates couldn't be enabled because the buffer isn't loaded or opts contained an invalid key.
func AttachBuffer(buffer Buffer, sendBuffer bool, opts map[string]interface{}) (attached bool) {
	name(nvim_buf_attach)
}

// DetachBuffer deactivate updates from this buffer to the current channel.
//
// Returns whether the updates couldn't be disabled because the buffer isn't loaded.
func DetachBuffer(buffer Buffer) (detached bool) {
	name(nvim_buf_detach)
}

// BufferLines gets a line-range from the buffer.
//
// Indexing is zero-based, end-exclusive.
// Negative indices are interpreted as length+1+index: -1 refers to the index past the end.
// So to get the last element use start=-2 and end=-1.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strictIndexing is set.
func BufferLines(buffer Buffer, start, end int, strictIndexing bool) (lines [][]byte) {
	name(nvim_buf_get_lines)
}

// SetBufferLines sets or replaces a line-range in the buffer.
//
// Indexing is zero-based, end-exclusive.
// Negative indices are interpreted as length+1+index: -1 refers to the index past the end.
// So to change or delete the last element use start=-2 and end=-1.
//
// To insert lines at a given index, set start and end args to the same index.
//
// To delete a range of lines, set replacement arg to an empty array.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless
// strict_indexing arg is set to true.
func SetBufferLines(buffer Buffer, start, end int, strictIndexing bool, replacement [][]byte) {
	name(nvim_buf_set_lines)
}

// SetBufferText sets or replaces a range in the buffer.
//
// This is recommended over SetBufferLines when only modifying parts of a
// line, as extmarks will be preserved on non-modified parts of the touched
// lines.
//
// Indexing is zero-based and end-exclusive.
//
// To insert text at a given index, set startRow and endRow args ranges to the same index.
//
// To delete a range, set replacement arg to an array containing an empty string, or simply an empty array.
//
// Prefer SetBufferLines when adding or deleting entire lines only.
func SetBufferText(buffer Buffer, startRow, startCol, endRow, endCol int, replacement [][]byte) {
	name(nvim_buf_set_text)
}

// BufferText gets a range from the buffer.
//
// This differs from BufferLines in that it allows retrieving only
// portions of a line.
//
// Indexing is zero-based. Column indices are end-exclusive.
//
// Prefer BufferLines when retrieving entire lines.
//
// opts is optional parameters. Currently unused.
func BufferText(buffer Buffer, startRow, startCol, endRow, endCol int, opts map[string]interface{}) [][]byte {
	name(nvim_buf_get_text)
}

// BufferOffset returns the byte offset of a line (0-indexed).
//
// Line 1 (index=0) has offset 0. UTF-8 bytes are counted. EOL is one byte.
// "fileformat" and "fileencoding" are ignored.
//
// The line index just after the last line gives the total byte-count of the buffer.
// A final EOL byte is counted if it would be written, see ":help eol".
//
// Unlike "line2byte" vim function, throws error for out-of-bounds indexing.
//
// If Buffer is unloaded buffer, returns -1.
func BufferOffset(buffer Buffer, index int) (offset int) {
	name(nvim_buf_get_offset)
}

// BufferVar gets a buffer-scoped (b:) variable.
func BufferVar(buffer Buffer, name string) (value interface{}) {
	name(nvim_buf_get_var)
}

// BufferChangedTick gets a changed tick of a buffer.
func BufferChangedTick(buffer Buffer) (changedtick int) {
	name(nvim_buf_get_changedtick)
}

// BufferKeymap gets a list of buffer-local mapping definitions.
//
// The mode short-name ("n", "i", "v", ...).
func BufferKeyMap(buffer Buffer, mode string) []*Mapping {
	name(nvim_buf_get_keymap)
}

// SetBufferKeyMap sets a buffer-local mapping for the given mode.
func SetBufferKeyMap(buffer Buffer, mode, lhs, rhs string, opts map[string]bool) {
	name(nvim_buf_set_keymap)
}

// DeleteBufferKeyMap unmaps a buffer-local mapping for the given mode.
func DeleteBufferKeyMap(buffer Buffer, mode, lhs string) {
	name(nvim_buf_del_keymap)
}

// SetBufferVar sets a buffer-scoped (b:) variable.
func SetBufferVar(buffer Buffer, name string, value interface{}) {
	name(nvim_buf_set_var)
}

// DeleteBufferVar removes a buffer-scoped (b:) variable.
func DeleteBufferVar(buffer Buffer, name string) {
	name(nvim_buf_del_var)
}

// BufferNumber gets a buffer's number.
//
// Deprecated: Use int(buffer) to get the buffer's number as an integer.
func BufferNumber(buffer Buffer) (number int) {
	name(nvim_buf_get_number)
	deprecatedSince(2)
}

// BufferName gets the full file name for the buffer.
func BufferName(buffer Buffer) (name string) {
	name(nvim_buf_get_name)
}

// SetBufferName sets the full file name for a buffer.
func SetBufferName(buffer Buffer, name string) {
	name(nvim_buf_set_name)
}

// IsBufferLoaded checks if a buffer is valid and loaded.
//
// See |help api-buffer| for more info about unloaded buffers.
func IsBufferLoaded(buffer Buffer) (loaded bool) {
	name(nvim_buf_is_loaded)
}

// DeleteBuffer deletes the buffer.
// See
//
//	:help :bwipeout
//
// The opts args is optional parameters.
//
//	force
//
// Force deletion and ignore unsaved changes. bool type.
//
//	unload
//
// Unloaded only, do not delete. See |help :bunload|. bool type.
func DeleteBuffer(buffer Buffer, opts map[string]bool) {
	name(nvim_buf_delete)
}

// IsBufferValid returns whether the buffer is valid.
//
// Note: Even if a buffer is valid it may have been unloaded.
// See |help api-buffer| for more info about unloaded buffers.
func IsBufferValid(buffer Buffer) (valied bool) {
	name(nvim_buf_is_valid)
}

// DeleteBufferMark deletes a named mark in the buffer.
// See |help mark-motions|.
func DeleteBufferMark(buffer Buffer, name string) (deleted bool) {
	name(nvim_buf_del_mark)
}

// SetBufferMark sets a named mark in the given buffer, all marks are allowed
// file/uppercase, visual, last change, etc.
// See |help mark-motions|.
//
// line and col are (1,0)-indexed.
//
// opts is optional parameters. Reserved for future use.
func SetBufferMark(buffer Buffer, name string, line, col int, opts map[string]interface{}) (set bool) {
	name(nvim_buf_set_mark)
}

// BufferMark return a tuple (row,col) representing the position of the named mark.
//
// Marks are (1,0)-indexed.
func BufferMark(buffer Buffer, name string) (pos [2]int) {
	name(nvim_buf_get_mark)
}

// ClearBufferHighlight clears highlights from a given source group and a range
// of lines.
//
// To clear a source group in the entire buffer, pass in 1 and -1 to startLine
// and endLine respectively.
//
// The lineStart and lineEnd parameters specify the range of lines to clear.
// The end of range is exclusive. Specify -1 to clear to the end of the file.
//
// Deprecated: Use ClearBufferNamespace instead.
func ClearBufferHighlight(buffer Buffer, srcID, startLine, endLine int) {
	name(nvim_buf_clear_highlight)
	deprecatedSince(7)
}

// SetBufferVirtualText set the virtual text (annotation) for a buffer line.
//
// By default (and currently the only option), the text will be placed after
// the buffer text.
//
// Virtual text will never cause reflow, rather virtual text will be truncated at the end of the screen line.
// The virtual text will begin one cell (|lcs-eol| or space) after the ordinary text.
//
// Namespaces are used to support batch deletion/updating of virtual text.
// To create a namespace, use CreateNamespace. Virtual text is cleared using ClearBufferNamespace.
//
// The same nsID can be used for both virtual text and highlights added by AddBufferHighlight,
// both can then be cleared with a single call to ClearBufferNamespace.
// If the virtual text never will be cleared by an API call, pass "nsID = -1".
//
// As a shorthand, "nsID = 0" can be used to create a new namespace for the
// virtual text, the allocated id is then returned.
//
// The opts arg is reserved for future use.
//
// Deprecated: Use SetBufferExtmark instead.
func SetBufferVirtualText(buffer Buffer, nsID, line int, chunks []TextChunk, opts map[string]interface{}) (id int) {
	name(nvim_buf_set_virtual_text)
	deprecatedSince(8)
}

// command.c

// ParseCmd parse command line.
//
// Doesn't check the validity of command arguments.
//
// str is command line string to parse. Cannot contain "\n".
//
// opts is optional parameters. Reserved for future use.
//
// Return dictionary containing command information, with these keys:
//
//	cmd: (string)
//
// Command name
//
//	range: ([]int)
//
// Command <range>. Can have 0-2 elements depending on how many items the range contains.
// Has no elements if command doesn't accept a range or if no range was specified,
// one element if only a single range item was specified and two elements if both range items were specified.
//
//	count: (int)
//
// Any |<count>| that was supplied to the command. -1 if command cannot take a count.
//
//	reg: (int)
//
// The optional command |<register>|, if specified. Empty string if not specified or if command cannot take a register.
//
//	bang: (bool)
//
// Whether command contains a |<bang>| (!) modifier.
//
//	args: ([]string)
//
// Command arguments.
//
//	addr: (string)
//
// Value of |:command-addr|. Uses short name.
//
//	nargs: (string)
//
// Value of |:command-nargs|.
//
//	nextcmd: (string)
//
// Next command if there are multiple commands separated by a |:bar|. Empty if there isn't a next command.
//
//	magic: (struct)
//
// Which characters have special meaning in the command arguments.
//
//	magic.file: (bool)
//
// The command expands filenames. Which means characters such as "%", "#" and wildcards are expanded.
//
//	magic.bar: (bool)
//
// The "|" character is treated as a command separator and the double quote character (\") is treated as the start of a comment.
//
//	mods: (struct)
//
// mods
//
//	mods.filter: (struct)
//
// |:filter|.
//
//	mods.filter.pattern: (string)
//
// Filter pattern. Empty string if there is no filter.
//
//	mods.filter.force: (bool)
//
// Whether filter is inverted or not.
//
//	mods.silent: (bool)
//
// silent
//
//	mods.emsg_silent: (bool)
//
// emsg_silent
//
//	mods.unsilent: (bool)
//
// unsilent
//
//	mods.sandbox: (bool)
//
// sandbox
//
//	mods.noautocmd: (bool)
//
// noautocmd
//
//	mods.browse: (bool)
//
// browse
//
//	mods.confirm: (bool)
//
// confirm
//
//	mods.hide: (bool)
//
// hide
//
//	mods.horizontal: (bool)
//
// horizontal
//
//	mods.keepalt: (bool)
//
// keepalt
//
//	mods.keepjumps: (bool)
//
// keepjumps
//
//	mods.keepmarks: (bool)
//
// keepmarks
//
//	mods.keeppatterns: (bool)
//
// keeppatterns
//
//	mods.lockmarks: (bool)
//
// lockmarks
//
//	mods.noswapfile: (bool)
//
// noswapfile
//
//	mods.tab: (int)
//
// -1 when omitted.
//
//	verbose: (int)
//
// -1 when omitted.
//
//	vertical: (bool)
//
// vertical
//
//	split: (string)
//
// Split modifier string, is an empty string when there's no split modifier. If there is a split modifier it can be one of:
//
//	aboveleft
//	belowright
//	topleft
//	botright
func ParseCmd(str string, opts map[string]interface{}) (cmd Cmd) {
	name(nvim_parse_cmd)
	returnPtr()
}

// Cmd executes an Ex command.
//
// Unlike Command() this command takes a structured Dictionary instead of a String. This
// allows for easier construction and manipulation of an Ex command. This also allows for things
// such as having spaces inside a command argument, expanding filenames in a command that otherwise
// doesn't expand filenames, etc.
//
// On execution error: fails with VimL error, updates v:errmsg.
// See Exec() and Command().
//
// cmd is the command to execute. Must be a Dictionary that can contain the same values
// as the return value of ParseCmd except "addr", "nargs" and "nextcmd" which are ignored if provided.
// All values except for "cmd" are optional.
//
// opts is the optional parameters.
//
//	output: (boolean, default false)
//
// Whether to return command output.
func Cmd(cmd *Cmd, opts map[string]bool) (output string) {
	name(nvim_cmd)
}

// CreateUserCommand create a new user command.
//
// name is name of the new user command. Must begin with an uppercase letter.
//
// command is replacement command to execute when this user command is executed.
// When called from Lua, the command can also be a Lua function.
//
// opts is optional command attributes. See |command-attributes| for more details.
//
// To use boolean attributes (such as |:command-bang| or |:command-bar|) set the value to "true".
// In addition to the string options listed in |:command-complete|,
// the "complete" key also accepts a Lua function which works like the "customlist" completion mode |:command-completion-customlist|.
//
//	desc (string)
//
// Used for listing the command when a Lua function is used for {command}.
//
//	force (bool, default true)
//
// Override any previous definition.
func CreateUserCommand(name string, command UserCommand, opts map[string]interface{}) {
	name(nvim_create_user_command)
}

// DeleteUserCommand delete a user-defined command.
func DeleteUserCommand(name string) {
	name(nvim_del_user_command)
}

// CreateBufferUserCommand create a new user command |user-commands| in the given buffer.
//
// Only commands created with |:command-buffer| or this function can be deleted with this function.
func CreateBufferUserCommand(buffer Buffer, name string, command UserCommand, opts map[string]interface{}) {
	name(nvim_buf_create_user_command)
}

// DeleteBufferUserCommand create a new user command |user-commands| in the given buffer.
//
// Only commands created with |:command-buffer| or this function can be deleted with this function.
func DeleteBufferUserCommand(buffer Buffer, name string) {
	name(nvim_buf_del_user_command)
}

// Commands gets a map of global (non-buffer-local) Ex commands.
// Currently only user-commands are supported, not builtin Ex commands.
//
// opts is optional parameters. Currently only supports:
//
//	{"builtin":false}
func Commands(opts map[string]interface{}) (commands map[string]*Command) {
	name(nvim_get_commands)
}

// BufferCommands gets a map of buffer-local user-commands.
//
// opts is optional parameters. Currently not used.
func BufferCommands(buffer Buffer, opts map[string]interface{}) map[string]*Command {
	name(nvim_buf_get_commands)
}

// tabpage.c

// TabpageWindows gets the windows in a tabpage.
func TabpageWindows(tabpage Tabpage) (windows []Window) {
	name(nvim_tabpage_list_wins)
}

// TabpageVar gets a tab-scoped (t:) variable.
func TabpageVar(tabpage Tabpage, name string) (value interface{}) {
	name(nvim_tabpage_get_var)
}

// SetTabpageVar sets a tab-scoped (t:) variable.
func SetTabpageVar(tabpage Tabpage, name string, value interface{}) {
	name(nvim_tabpage_set_var)
}

// DeleteTabpageVar removes a tab-scoped (t:) variable.
func DeleteTabpageVar(tabpage Tabpage, name string) {
	name(nvim_tabpage_del_var)
}

// TabpageWindow gets the current window in a tabpage.
func TabpageWindow(tabpage Tabpage) Window {
	name(nvim_tabpage_get_win)
}

// TabpageNumber gets the tabpage number.
func TabpageNumber(tabpage Tabpage) (number int) {
	name(nvim_tabpage_get_number)
}

// IsTabpageValid checks if a tabpage is valid.
func IsTabpageValid(tabpage Tabpage) (valid bool) {
	name(nvim_tabpage_is_valid)
}

// extmark.c

// CreateNamespace creates a new namespace, or gets an existing one.
//
// Namespaces are used for buffer highlights and virtual text, see
// AddBufferHighlight and SetBufferVirtualText.
//
// Namespaces can be named or anonymous. If "name" matches an existing namespace,
// the associated id is returned. If "name" is an empty string a new, anonymous
// namespace is created.
//
// The returns the namespace ID.
func CreateNamespace(name string) (nsID int) {
	name(nvim_create_namespace)
}

// BufferExtmarkByID beturns position for a given extmark id.
//
// opts is optional parameters.
//
//	details
//
// Whether to include the details dict. bool type.
func BufferExtmarkByID(buffer Buffer, nsID, id int, opt map[string]interface{}) (pos []int) {
	name(nvim_buf_get_extmark_by_id)
}

// BufferExtmarks gets extmarks in "traversal order" from a |charwise| region defined by
// buffer positions (inclusive, 0-indexed).
//
// Region can be given as (row,col) tuples, or valid extmark ids (whose
// positions define the bounds).
// 0 and -1 are understood as (0,0) and (-1,-1) respectively, thus the following are equivalent:
//
//	BufferExtmarks(0, myNS, 0, -1, {})
//	BufferExtmarks(0, myNS, [0,0], [-1,-1], {})
//
// If end arg is less than start arg, traversal works backwards.
// It useful with limit arg, to get the first marks prior to a given position.
//
// The start and end args is start or end of range, given as (row, col), or
// valid extmark id whose position defines the bound.
//
// opts is optional parameters.
//
//	limit
//
// Maximum number of marks to return. int type.
//
//	details
//
// Whether to include the details dict. bool type.
func BufferExtmarks(buffer Buffer, nsID int, start, end interface{}, opt map[string]interface{}) (marks []ExtMark) {
	name(nvim_buf_get_extmarks)
}

// SetBufferExtmark creates or updates an extmark.
//
// To create a new extmark, pass id=0. The extmark id will be returned.
// To move an existing mark, pass its id.
//
// It is also allowed to create a new mark by passing in a previously unused
// id, but the caller must then keep track of existing and unused ids itself.
// (Useful over RPC, to avoid waiting for the return value.)
//
// Using the optional arguments, it is possible to use this to highlight
// a range of text, and also to associate virtual text to the mark.
//
// The opts arg is optional parameters.
//
//	id
//
// ID of the extmark to edit.
//
//	end_line
//
// Ending line of the mark, 0-based inclusive.
//
//	end_col
//
// Ending col of the mark, 0-based inclusive.
//
//	hl_group
//
// Name of the highlight group used to highlight this mark.
//
//	virt_text
//
// Virtual text to link to this mark.
//
//	virt_text_pos
//
// Positioning of virtual text.
// Possible values:
//
//	eol
//
// right after eol character (default)
//
//	overlay
//
// display over the specified column, without shifting the underlying text.
//
//	virt_text_win_col
//
// position the virtual text at a fixed window column (starting from the first text column)
//
//	virt_text_hide
//
// Hide the virtual text when the background text is selected or hidden due to horizontal scroll "nowrap".
//
//	hl_mode
//
// Control how highlights are combined with the highlights of the text. Currently only affects
// virt_text highlights, but might affect "hl_group" in later versions.
// Possible values:
//
//	replace
//
// only show the virt_text color. This is the default.
//
//	combine
//
// combine with background text color
//
//	blend
//
// blend with background text color.
//
//	hl_eol
//
// when true, for a multiline highlight covering the EOL of a line, continue the highlight for the rest
// of the screen line (just like for diff and cursorline highlight).
//
//	ephemeral
//
// For use with "nvim_set_decoration_provider" callbacks. The mark will only be used for the current redraw cycle,
// and not be permantently stored in the buffer.
//
//	right_gravity
//
// Boolean that indicates the direction the extmark will be shifted in when new text is
// inserted (true for right, false for left).  defaults to true.
//
//	end_right_gravity
//
// Boolean that indicates the direction the extmark end position (if it exists) will be
// shifted in when new text is inserted (true for right, false for left). Defaults to false.
//
//	priority
//
// A priority value for the highlight group. For example treesitter highlighting uses a value of 100.
func SetBufferExtmark(buffer Buffer, nsID, line, col int, opts map[string]interface{}) (id int) {
	name(nvim_buf_set_extmark)
}

// DeleteBufferExtmark removes an extmark.
//
// THe returns whether the extmark was found.
func DeleteBufferExtmark(buffer Buffer, nsID, extmarkID int) (deleted bool) {
	name(nvim_buf_del_extmark)
}

// AddBufferHighlight adds a highlight to buffer.
//
// IT useful for plugins that dynamically generate highlights to a buffer like a semantic highlighter or linter.
//
// The function adds a single highlight to a buffer.
// Unlike |matchaddpos()| vim function, highlights follow changes to line numbering as lines are
// inserted/removed above the highlighted line, like signs and marks do.
//
// Namespaces are used for batch deletion/updating of a set of highlights.
// To create a namespace, use CreateNamespace which returns a namespace id.
// Pass it in to this function as nsID to add highlights to the namespace.
// All highlights in the same namespace can then be cleared with single call to ClearBufferNamespace.
// If the highlight never will be deleted by an API call, pass nsID = -1.
//
// As a shorthand, "srcID = 0" can be used to create a new namespace for the
// highlight, the allocated id is then returned.
//
// If hlGroup arg is the empty string, no highlight is added, but a new `nsID` is still returned.
// This is supported for backwards compatibility, new code should use CreateNamespaceto create a new empty namespace.
func AddBufferHighlight(buffer Buffer, srcID int, hlGroup string, line, startCol, endCol int) (id int) {
	name(nvim_buf_add_highlight)
}

// ClearBufferNamespace clears namespaced objects (highlights, extmarks, virtual text) from a region.
// Lines are 0-indexed.
//
// To clear the namespace in the entire buffer, specify line_start=0 and line_end=-1.
func ClearBufferNamespace(buffer Buffer, nsID, lineStart, lineEnd int) {
	name(nvim_buf_clear_namespace)
}

// options.c

// OptionValue gets the value of an option.
//
// The behavior of this function matches that of |:set|: the local value of an option is returned if it exists; otherwise,
// the global value is returned.
// Local values always correspond to the current buffer or window.
//
// To get a buffer-local or window-local option for a specific buffer or window, use BufferOption() or WindowOption().
//
// name is the option name.
//
// opts is the Optional parameters.
//
//	scope
//
// Analogous to |:setglobal| and |:setlocal|, respectively.
func OptionValue(name string, opts map[string]OptionValueScope) (optionValue interface{}) {
	name(nvim_get_option_value)
}

// SetOptionValue sets the value of an option. The behavior of this function matches that of
// |:set|: for global-local options, both the global and local value are set
// unless otherwise specified with {scope}.
// name is the option name.
//
// opts is the Optional parameters.
//
//	scope
//
// Analogous to |:setglobal| and |:setlocal|, respectively.
func SetOptionValue(name string, value interface{}, opts map[string]OptionValueScope) {
	name(nvim_set_option_value)
}

// AllOptionsInfo gets the option information for all options.
//
// The dictionary has the full option names as keys and option metadata
// dictionaries as detailed at OptionInfo.
//
// Resulting map has keys:
//
//	name
//
// Name of the option (like "filetype").
//
//	shortname
//
// Shortened name of the option (like "ft").
//
//	type
//
// type of option ("string", "number" or "boolean").
//
//	default
//
// The default value for the option.
//
//	was_set
//
// Whether the option was set.
//
//	last_set_sid
//
// Last set script id (if any).
//
//	last_set_linenr
//
// line number where option was set.
//
//	last_set_chan
//
// Channel where option was set (0 for local).
//
//	scope
//
// One of "global", "win", or "buf".
//
//	global_local
//
// Whether win or buf option has a global value.
//
//	commalist
//
// List of comma separated values.
//
//	flaglist
//
// List of single char flags.
func AllOptionsInfo() (opinfo OptionInfo) {
	name(nvim_get_all_options_info)
	returnPtr()
}

// OptionInfo gets the option information for one option.
//
// Resulting map has keys:
//
//	name
//
// Name of the option (like "filetype").
//
//	shortname
//
// Shortened name of the option (like "ft").
//
//	type
//
// type of option ("string", "number" or "boolean").
//
//	default
//
// The default value for the option.
//
//	was_set
//
// Whether the option was set.
//
//	last_set_sid
//
// Last set script id (if any).
//
//	last_set_linenr
//
// line number where option was set.
//
//	last_set_chan
//
// Channel where option was set (0 for local).
//
//	scope
//
// One of "global", "win", or "buf".
//
//	global_local
//
// Whether win or buf option has a global value.
//
//	commalist
//
// List of comma separated values.
//
//	flaglist
//
// List of single char flags.
func OptionInfo(name string, opts map[string]interface{}) (opinfo OptionInfo) {
	name(nvim_get_option_info2)
	returnPtr()
}

// SetOption sets an option value.
func SetOption(name string, value interface{}) {
	name(nvim_set_option)
}

// Option gets an option value string.
func Option(name string) (option interface{}) {
	name(nvim_get_option)
}

// BufferOption gets a buffer option value.
func BufferOption(buffer Buffer, name string) (value interface{}) {
	name(nvim_buf_get_option)
}

// SetBufferOption sets a buffer option value.
//
// Passing nil as value arg to deletes the option (only works if there's a global fallback).
func SetBufferOption(buffer Buffer, name string, value interface{}) {
	name(nvim_buf_set_option)
}

// WindowOption gets a window option value.
func WindowOption(window Window, name string) (value interface{}) {
	name(nvim_win_get_option)
}

// SetWindowOption sets a window option value. Passing "nil" as value deletes the option(only works if there's a global fallback).
func SetWindowOption(window Window, name string, value interface{}) {
	name(nvim_win_set_option)
}

// ui.c

// AttachUI registers the client as a remote UI. After this method is called,
// the client will receive redraw notifications.
//
//	:help rpc-remote-ui
//
// The redraw notification method has variadic arguments. Register a handler
// for the method like this:
//
//	v.RegisterHandler("redraw", func(updates ...[]interface{}) {
//	    for _, update := range updates {
//	        // handle update
//	    }
//	})
func AttachUI(width, height int, options map[string]interface{}) {
	name(nvim_ui_attach)
}

// SetFocusUI tells the nvim server if focus was gained or lost by the GUI.
func SetFocusUI(gained bool) {
	name(nvim_ui_set_focus)
}

// DetachUI unregisters the client as a remote UI.
func DetachUI() {
	name(nvim_ui_detach)
}

// TryResizeUI notifies Nvim that the client window has resized. If possible,
// Nvim will send a redraw request to resize.
func TryResizeUI(width, height int) {
	name(nvim_ui_try_resize)
}

// SetUIOption sets a UI option.
func SetUIOption(name string, value interface{}) {
	name(nvim_ui_set_option)
}

// TryResizeUIGrid tell Nvim to resize a grid. Triggers a grid_resize event with the requested
// grid size or the maximum size if it exceeds size limits.
//
// On invalid grid handle, fails with error.
func TryResizeUIGrid(grid, width, height int) {
	name(nvim_ui_try_resize_grid)
}

// SetPumHeight tells Nvim the number of elements displaying in the popumenu, to decide
// <PageUp> and <PageDown> movement.
//
// height is popupmenu height, must be greater than zero.
func SetPumHeight(height int) {
	name(nvim_ui_pum_set_height)
}

// SetPumBounds tells Nvim the geometry of the popumenu, to align floating windows with an
// external popup menu.
//
// Note that this method is not to be confused with SetPumHeight,
// which sets the number of visible items in the popup menu, while this
// function sets the bounding box of the popup menu, including visual
// elements such as borders and sliders.
//
// Floats need not use the same font size, nor be anchored to exact grid corners, so one can set floating-point
// numbers to the popup menu geometry.
func SetPumBounds(width, height, row, col float64) {
	name(nvim_ui_pum_set_bounds)
}

// vimscript.c

// Exec executes Vimscript (multiline block of Ex-commands), like anonymous source.
//
// Unlike Command, this function supports heredocs, script-scope (s:), etc.
//
// When fails with VimL error, does not update "v:errmsg".
func Exec(src string, opts map[string]interface{}) (out map[string]interface{}) {
	name(nvim_exec2)
}

// Command executes an ex-command.
//
// When fails with VimL error, does not update "v:errmsg".
func Command(cmd string) {
	name(nvim_command)
}

// ParseExpression parse a VimL expression.
func ParseExpression(expr, flags string, highlight bool) (expression map[string]interface{}) {
	name(nvim_parse_expression)
}

// vim.c

// HL gets a highlight definition by name.
//
// nsID get highlight groups for namespace ns_id [Namespaces]. Use 0 to get global highlight groups |:highlight|.
//
// opts dict:
//
//	name
//
// Get a highlight definition by name.
//
//	id
//
// Get a highlight definition by id.
//
//	link
//
// Show linked group name instead of effective definition.
//
// The returned HLAttrs highlight groups as a map from group name to a highlight definition map as in SetHighlight, or only a single highlight definition map if requested by name or id.
func HL(nsID int, opts map[string]interface{}) (highlight HLAttrs) {
	name(nvim_get_hl)
	returnPtr()
}

// HLByID gets a highlight definition by name.
//
// hlID is the highlight id as returned by HLIDByName.
//
// rgb is the whether the export RGB colors.
//
// The returned highlight is the highlight definition.
func HLByID(hlID int, rgb bool) (highlight HLAttrs) {
	name(nvim_get_hl_by_id)
	returnPtr()
	deprecatedSince(9)
}

// HLIDByName gets a highlight group by name.
//
// name is the Highlight group name.
//
// The returns hlID is the highlight id.
//
// This function similar to HLByID, but allocates a new ID if not present.
func HLIDByName(name string) (hlID int) {
	name(nvim_get_hl_id_by_name)
}

// HLByName gets a highlight definition by id.
//
// name is Highlight group name.
//
// rgb is whether the export RGB colors.
//
// The returned highlight is the highlight definition.
func HLByName(name string, rgb bool) (highlight HLAttrs) {
	name(nvim_get_hl_by_name)
	returnPtr()
	deprecatedSince(9)
}

// SetHighlight sets a highlight group.
//
// nsID is number of namespace for this highlight.
//
// name is highlight group name, like "ErrorMsg".
//
// val is highlight definiton map, like HLByName.
//
// in addition the following keys are also recognized:
//
//	default
//
// don't override existing definition, like "hi default".
func SetHighlight(nsID int, name string, val *HLAttrs) {
	name(nvim_set_hl)
}

// SetHighlightNamespace set active namespace for highlights. This can be set for a single window,
//
// See SetWindowHeightNamespace.
func SetHighlightNamespace(nsID int) {
	name(nvim_set_hl_ns)
}

// SetFastHighlightNamespace set active namespace for highlights while redrawing.
//
// This function meant to be called while redrawing, primarily from
// nvim_set_decoration_provider on_win and on_line callbacks, which
// are allowed to change the namespace during a redraw cycle.
func SetFastHighlightNamespace(nsID int) {
	name(nvim_set_hl_ns_fast)
}

// FeedKeys input-keys to Nvim, subject to various quirks controlled by "mode"
// flags. Unlike Input, this is a blocking call.
//
// This function does not fail, but updates "v:errmsg".
//
// If need to input sequences like <C-o> use ReplaceTermcodes to
// replace the termcodes and then pass the resulting string to nvim_feedkeys.
// You'll also want to enable escape_csi.
//
// mode is following character flags:
//
//	m
//
// Remap keys. This is default.
//
//	n
//
// Do not remap keys.
//
//	t
//
// Handle keys as if typed; otherwise they are handled as if coming from a mapping.
// This matters for undo, opening folds, etc.
//
// escapeCSI is whether the escape K_SPECIAL/CSI bytes in keys.
func FeedKeys(keys, mode string, escapeCSI bool) {
	name(nvim_feedkeys)
}

// Input queues raw user-input.
//
// Unlike FeedKeys, this uses a low-level input buffer and the call
// is non-blocking (input is processed asynchronously by the eventloop).
//
// This function does not fail but updates "v:errmsg".
//
// keys is to be typed.
//
// Note: "keycodes" like "<CR>" are translated, so "<" is special. To input a literal "<", send "<LT>".
//
// Note: For mouse events use InputMouse. The pseudokey form "<LeftMouse><col,row>" is deprecated.
//
// The returned written is number of bytes actually written (can be fewer than
// requested if the buffer becomes full).
func Input(keys string) (written int) {
	name(nvim_input)
}

// InputMouse Send mouse event from GUI.
//
// This API is non-blocking. It does not wait on any result, but queues the event to be
// processed soon by the event loop.
//
// button is mouse button. One of
//
//	left
//	right
//	middle
//	wheel
//
// action is for ordinary buttons. One of
//
//	press
//	drag
//	release
//
// For the wheel, One of
//
//	up
//	down
//	left
//	right
//
// modifier is string of modifiers each represented by a single char.
// The same specifiers are used as for a key press, except
// that the "-" separator is optional, so "C-A-", "c-a"
// and "CA" can all be used to specify "Ctrl+Alt+Click".
//
// grid is grid number if the client uses "ui-multigrid", else 0.
//
// row is mouse row-position (zero-based, like redraw events).
//
// col is mouse column-position (zero-based, like redraw events).
func InputMouse(button, action, modifier string, grid, row, col int) {
	name(nvim_input_mouse)
}

// ReplaceTermcodes replaces terminal codes and "keycodes" (<CR>, <Esc>, ...) in a string with
// the internal representation.
//
// str is string to be converted.
//
// fromPart is legacy Vim parameter. Usually true.
//
// doLT is also translate <lt>. Ignored if "special" is false.
//
// special is replace "keycodes", e.g. "<CR>" becomes a "\n" char.
//
// The returned sequences are Nvim's internal representation of keys, for example:
//
//	<esc> -> '\x1b'
//	<cr>  -> '\r'
//	<c-l> -> '\x0c'
//	<up>  -> '\x80ku'
//
// The returned sequences can be used as input to feedkeys.
func ReplaceTermcodes(str string, fromPart, doLT, special bool) (input string) {
	name(nvim_replace_termcodes)
}

// CommandOutput executes a single ex command and returns the output.
//
// Deprecated: Use Exec instead.
func CommandOutput(cmd string) (out string) {
	name(nvim_command_output)
	deprecatedSince(7)
}

// Eval evaluates a VimL expression.
//
// Dictionaries and Lists are recursively expanded.
//
// Fails with VimL error, does not update "v:errmsg".
//
// expr is VimL expression string.
//
//	:help expression
func Eval(expr string) (result interface{}) {
	name(nvim_eval)
}

// StringWidth calculates the number of display cells occupied by "text".
//
// "<Tab>" counts as one cell.
func StringWidth(s string) (width int) {
	name(nvim_strwidth)
}

// RuntimePaths gets the paths contained in "runtimepath".
func RuntimePaths() (paths []string) {
	name(nvim_list_runtime_paths)
}

// RuntimeFiles find files in runtime directories.
//
// name is can contain wildcards.
//
// For example,
//
//	RuntimeFiles("colors/*.vim", true)
//
// will return all color scheme files.
//
// Always use forward slashes (/) in the search pattern for subdirectories regardless of platform.
//
// It is not an error to not find any files, returned an empty array.
//
// To find a directory, name must end with a forward slash, like
// "rplugin/python/".
// Without the slash it would instead look for an ordinary file called "rplugin/python".
//
// all is whether to return all matches or only the first.
func RuntimeFiles(name string, all bool) (files []string) {
	name(nvim_get_runtime_file)
}

// SetCurrentDirectory changes the global working directory.
func SetCurrentDirectory(dir string) {
	name(nvim_set_current_dir)
}

// CurrentLine gets the current line.
func CurrentLine() (line []byte) {
	name(nvim_get_current_line)
}

// SetCurrentLine sets the current line.
func SetCurrentLine(line []byte) {
	name(nvim_set_current_line)
}

// DeleteCurrentLine deletes the current line.
func DeleteCurrentLine() {
	name(nvim_del_current_line)
}

// Var gets a global (g:) variable.
func Var(name string) (value interface{}) {
	name(nvim_get_var)
}

// SetVar sets a global (g:) variable.
func SetVar(name string, value interface{}) {
	name(nvim_set_var)
}

// DeleteVar removes a global (g:) variable.
func DeleteVar(name string) {
	name(nvim_del_var)
}

// VVar gets a v: variable.
func VVar(name string) (value interface{}) {
	name(nvim_get_vvar)
}

// SetVVar sets a v: variable, if it is not readonly.
func SetVVar(name string, value interface{}) {
	name(nvim_set_vvar)
}

// Echo echo a message.
//
// chunks is a list of [text, hl_group] arrays, each representing a
// text chunk with specified highlight. hl_group element can be omitted for no highlight.
//
// If history is true, add to "message-history".
//
// opts is optional parameters. Reserved for future use.
func Echo(chunks []TextChunk, history bool, opts map[string]interface{}) {
	name(nvim_echo)
}

// WriteOut writes a message to the Vim output buffer.
//
// Does not append "\n", the message is buffered (won't display) until a linefeed is written.
func WriteOut(str string) {
	name(nvim_out_write)
}

// WriteErr writes a message to the Vim error buffer.
//
// Does not append "\n", the message is buffered (won't display) until a linefeed is written.
func WriteErr(str string) {
	name(nvim_err_write)
}

// WritelnErr writes a message to the Vim error buffer.
//
// Appends "\n", so the buffer is flushed and displayed.
func WritelnErr(str string) {
	name(nvim_err_writeln)
}

// Buffers gets the current list of buffer handles.
//
// Includes unlisted (unloaded/deleted) buffers, like ":ls!". Use IsBufferLoaded to check if a buffer is loaded.
func Buffers() (buffers []Buffer) {
	name(nvim_list_bufs)
}

// CurrentBuffer gets the current buffer.
func CurrentBuffer() (buffer Buffer) {
	name(nvim_get_current_buf)
}

// SetCurrentBuffer sets the current buffer.
func SetCurrentBuffer(buffer Buffer) {
	name(nvim_set_current_buf)
}

// Windows gets the current list of window handles.
func Windows() (windows []Window) {
	name(nvim_list_wins)
}

// CurrentWindow gets the current window.
func CurrentWindow() (window Window) {
	name(nvim_get_current_win)
}

// SetCurrentWindow sets the current window.
func SetCurrentWindow(window Window) {
	name(nvim_set_current_win)
}

// CreateBuffer creates a new, empty, unnamed buffer.
//
// listed is sets buflisted buffer opttion. If false, sets "nobuflisted".
//
// scratch is creates a "throwaway" for temporary work (always 'nomodified').
//
//	bufhidden=hide
//	buftype=nofile
//	noswapfile
//	nomodeline
func CreateBuffer(listed, scratch bool) (buffer Buffer) {
	name(nvim_create_buf)
}

// OpenTerm opens a terminal instance in a buffer.
//
// By default (and currently the only option) the terminal will not be
// connected to an external process. Instead, input send on the channel
// will be echoed directly by the terminal. This is useful to disply
// ANSI terminal sequences returned as part of a rpc message, or similar.
//
// Note that to directly initiate the terminal using the right size, display the
// buffer in a configured window before calling this. For instance, for a
// floating display, first create an empty buffer using CreateBuffer,
// then display it using OpenWindow, and then call this function.
// Then "nvim_chan_send" cal be called immediately to process sequences
// in a virtual terminal having the intended size.
//
// buffer is the buffer to use (expected to be empty).
//
// opts is optional parameters. Reserved for future use.
func OpenTerm(buffer Buffer, opts map[string]interface{}) (channel int) {
	name(nvim_open_term)
}

// OpenWindow open a new window.
//
// Currently this is used to open floating and external windows.
// Floats are windows that are drawn above the split layout, at some anchor
// position in some other window.
// Floats can be drawn internally or by external GUI with the "ui-multigrid" extension.
// External windows are only supported with multigrid GUIs, and are displayed as separate top-level windows.
//
// For a general overview of floats, see
//
//	:help api-floatwin
//
// Exactly one of "external" and "relative" must be specified.
// The "width" and "height" of the new window must be specified.
//
// With relative=editor (row=0,col=0) refers to the top-left corner of the
// screen-grid and (row=Lines-1,col=Columns-1) refers to the bottom-right
// corner.
// Fractional values are allowed, but the builtin implementation
// (used by non-multigrid UIs) will always round down to nearest integer.
//
// Out-of-bounds values, and configurations that make the float not fit inside
// the main editor, are allowed.
// The builtin implementation truncates values so floats are fully within the main screen grid.
// External GUIs could let floats hover outside of the main window like a tooltip, but
// this should not be used to specify arbitrary WM screen positions.
func OpenWindow(buffer Buffer, enter bool, config *WindowConfig) (window Window) {
	name(nvim_open_win)
}

// Tabpages gets the current list of tabpage handles.
func Tabpages() (tabpages []Tabpage) {
	name(nvim_list_tabpages)
}

// CurrentTabpage gets the current tabpage.
func CurrentTabpage() (tabpage Tabpage) {
	name(nvim_get_current_tabpage)
}

// SetCurrentTabpage sets the current tabpage.
func SetCurrentTabpage(tabpage Tabpage) {
	name(nvim_set_current_tabpage)
}

// Namespaces gets existing, non-anonymous namespaces.
//
// The return dict that maps from names to namespace ids.
func Namespaces() (namespaces map[string]int) {
	name(nvim_get_namespaces)
}

// Paste pastes at cursor, in any mode.
//
// Invokes the "vim.paste" handler, which handles each mode appropriately.
// Sets redo/undo. Faster than Input(). Lines break at LF ("\n").
//
// Errors ("nomodifiable", "vim.paste()" "failure" ...) are reflected in `err`
// but do not affect the return value (which is strictly decided by `vim.paste()`).
//
// On error, subsequent calls are ignored ("drained") until the next paste is initiated (phase 1 or -1).
//
//	data
//
// multiline input. May be binary (containing NUL bytes).
//
//	crlf
//
// also break lines at CR and CRLF.
//
//	phase
//
// -1 is paste in a single call (i.e. without streaming).
//
// To stream a paste, call Paste sequentially with these phase args:
//
//	1
//
// starts the paste (exactly once)
//
//	2
//
// continues the paste (zero or more times)
//
//	3
//
// ends the paste (exactly once)
//
// The returned boolean state is:
//
//	true
//
// Client may continue pasting.
//
//	false
//
// Client must cancel the paste.
func Paste(data string, crlf bool, phase int) (state bool) {
	name(nvim_paste)
}

// Put puts text at cursor, in any mode.
//
// Compare :put and p which are always linewise.
//
// lines is readfile() style list of lines.
//
// typ is edit behavior: any getregtype() result, or:
//
//	 b
//	blockwise-visual mode (may include width, e.g. "b3")
//	 c
//	characterwise mode
//	 l
//	linewise mode
//	 ""
//
// guess by contents, see |setreg()|.
//
// After is insert after cursor (like `p`), or before (like `P`).
//
// follow arg is place cursor at end of inserted text.
func Put(lines []string, typ string, after, follow bool) {
	name(nvim_put)
}

// Subscribe subscribes to event broadcasts.
func Subscribe(event string) {
	name(nvim_subscribe)
}

// Unsubscribe unsubscribes to event broadcasts.
func Unsubscribe(event string) {
	name(nvim_unsubscribe)
}

// ColorByName Returns the 24-bit RGB value of a ColorMap color name or "#rrggbb" hexadecimal string.
//
// Example:
//
//	ColorByName("Pink")
//	ColorByName("#cbcbcb")
func ColorByName(name string) (color int) {
	name(nvim_get_color_by_name)
}

// ColorMap returns a map of color names and RGB values.
//
// Keys are color names (e.g. "Aqua") and values are 24-bit RGB color values (e.g. 65535).
//
// The returns map is color names and RGB values.
func ColorMap() (colorMap map[string]int) {
	name(nvim_get_color_map)
}

// Context gets a map of the current editor state.
// This API still under development.
//
// The opts arg is optional parameters.
// Key is "types".
//
// List of context-types to gather, or empty for "all" context.
//
//	regs
//	jumps
//	bufs
//	gvars
//	funcs
//	sfuncs
func Context(opts map[string][]string) (context map[string]interface{}) {
	name(nvim_get_context)
}

// LoadContext Sets the current editor state from the given context map.
func LoadContext(context map[string]interface{}) (contextMap interface{}) {
	name(nvim_load_context)
}

// Mode gets the current mode.
//
// |mode()| "blocking" is true if Nvim is waiting for input.
func Mode() (mode Mode) {
	name(nvim_get_mode)
	returnPtr()
}

// KeyMap gets a list of global (non-buffer-local) |mapping| definitions.
//
// The mode arg is the mode short-name, like "n", "i", "v" or etc.
func KeyMap(mode string) (maps []*Mapping) {
	name(nvim_get_keymap)
}

// SetKeyMap sets a global mapping for the given mode.
//
// To set a buffer-local mapping, use SetBufferKeyMap().
//
// Unlike :map, leading/trailing whitespace is accepted as part of the {lhs}
// or {rhs}.
// Empty {rhs} is <Nop>. keycodes are replaced as usual.
//
//	mode
//
// mode short-name (map command prefix: "n", "i", "v", "x", …) or "!" for :map!, or empty string for :map.
//
//	lhs
//
// Left-hand-side {lhs} of the mapping.
//
//	rhs
//
// Right-hand-side {rhs} of the mapping.
//
//	opts
//
// Optional parameters map. Accepts all ":map-arguments" as keys excluding "buffer" but including "noremap".
// Values are Booleans. Unknown key is an error.
func SetKeyMap(mode, lhs, rhs string, opts map[string]bool) {
	name(nvim_set_keymap)
}

// DeleteKeyMap unmaps a global mapping for the given mode.
//
// To unmap a buffer-local mapping, use DeleteBufferKeyMap().
func DeleteKeyMap(mode, lhs string) {
	name(nvim_del_keymap)
}

// APIInfo returns a 2-tuple (Array), where item 0 is the current channel id and item
// 1 is the "api-metadata" map (Dictionary).
//
// Returns 2-tuple [{channel-id}, {api-metadata}].
func APIInfo() (apiInfo []interface{}) {
	name(nvim_get_api_info)
}

// SetClientInfo self-identifies the client.
//
// The client/plugin/application should call this after connecting, to provide
// hints about its identity and purpose, for debugging and orchestration.
//
// Can be called more than once; the caller should merge old info if
// appropriate. Example: library first identifies the channel, then a plugin
// using that library later identifies itself.
func SetClientInfo(name string, version ClientVersion, typ ClientType, methods map[string]*ClientMethod, attributes ClientAttributes) {
	name(nvim_set_client_info)
}

// ChannelInfo get information about a channel.
//
// Rreturns Dictionary describing a channel, with these keys:
//
//	stream
//
// The stream underlying the channel. value are:
//
//	stdio
//
// stdin and stdout of this Nvim instance.
//
//	stderr
//
// stderr of this Nvim instance.
//
//	socket
//
// TCP/IP socket or named pipe.
//
//	job
//
// job with communication over its stdio.
//
//	mode
//
// How data received on the channel is interpreted. value are:
//
//	bytes
//
// send and receive raw bytes.
//
//	terminal
//
// A terminal instance interprets ASCII sequences.
//
//	rpc
//
// RPC communication on the channel is active.
//
//	pty
//
// Name of pseudoterminal, if one is used (optional).
// On a POSIX system, this will be a device path like /dev/pts/1.
// Even if the name is unknown, the key will still be present to indicate a pty is used.
// This is currently the case when using winpty on windows.
//
//	buffer
//
// Buffer with connected |terminal| instance (optional).
//
//	client
//
// Information about the client on the other end of the RPC channel, if it has added it using SetClientInfo() (optional).
func ChannelInfo(channelID int) (channel Channel) {
	name(nvim_get_chan_info)
	returnPtr()
}

// Channels get information about all open channels.
func Channels() (channels []*Channel) {
	name(nvim_list_chans)
}

// UIs gets a list of dictionaries representing attached UIs.
func UIs() (uis []*UI) {
	name(nvim_list_uis)
}

// ProcChildren gets the immediate children of process `pid`.
func ProcChildren(pid int) (processes []uint) {
	name(nvim_get_proc_children)
}

// Proc gets info describing process "pid".
func Proc(pid int) (process Process) {
	name(nvim_get_proc)
}

// SelectPopupmenuItem selects an item in the completion popupmenu.
//
// If |ins-completion| is not active this API call is silently ignored.
// Useful for an external UI using |ui-popupmenu| to control the popupmenu
// with the mouse. Can also be used in a mapping; use <cmd> |:map-cmd| to
// ensure the mapping doesn't end completion mode.
//
// opts optional parameters. Reserved for future use.
func SelectPopupmenuItem(item int, insert, finish bool, opts map[string]interface{}) {
	name(nvim_select_popupmenu_item)
}

// DeleteMark deletes a uppercase/file named mark.
// See |help mark-motions|.
func DeleteMark(name string) (deleted bool) {
	name(nvim_del_mark)
}

// Mark returns a tuple (row, col, buffer, buffername) representing the position of
// the uppercase/file named mark.
// See |help mark-motions|.
//
// opts is optional parameters. Reserved for future use.
func Mark(name string, opts map[string]interface{}) (mark Mark) {
	name(nvim_get_mark)
	returnPtr()
}

// EvalStatusLine evaluates statusline string.
//
// opts optional parameters.
//
//	winid (int)
//
// Window ID of the window to use as context for statusline.
//
//	maxwidth (int)
//
// Maximum width of statusline.
//
//	fillchar (string)
//
// Character to fill blank spaces in the statusline (see 'fillchars').
//
//	highlights (bool)
//
// Return highlight information.
//
//	use_tabline (bool)
//
// Evaluate tabline instead of statusline. When true, {winid} is ignored.
func EvalStatusLine(name string, opts map[string]interface{}) (statusline map[string]interface{}) {
	name(nvim_eval_statusline)
}

// window.c

// WindowBuffer gets the current buffer in a window.
func WindowBuffer(window Window) (buffer Buffer) {
	name(nvim_win_get_buf)
}

// SetBufferToWindow Sets the current buffer in a window, without side-effects.
func SetBufferToWindow(window Window, buffer Buffer) {
	name(nvim_win_set_buf)
}

// WindowCursor gets the (1,0)-indexed cursor position in the window.
func WindowCursor(window Window) (pos [2]int) {
	name(nvim_win_get_cursor)
}

// SetWindowCursor sets the (1,0)-indexed cursor position in the window.
func SetWindowCursor(window Window, pos [2]int) {
	name(nvim_win_set_cursor)
}

// WindowHeight returns the window height.
func WindowHeight(window Window) (height int) {
	name(nvim_win_get_height)
}

// SetWindowHeight Sets the window height. This will only succeed if the screen is split horizontally.
func SetWindowHeight(window Window, height int) {
	name(nvim_win_set_height)
}

// WindowWidth returns the window width.
func WindowWidth(window Window) (width int) {
	name(nvim_win_get_width)
}

// SetWindowWidth Sets the window width. This will only succeed if the screen is split vertically.
func SetWindowWidth(window Window, width int) {
	name(nvim_win_set_width)
}

// WindowVar gets a window-scoped (w:) variable.
func WindowVar(window Window, name string) (value interface{}) {
	name(nvim_win_get_var)
}

// SetWindowVar sets a window-scoped (w:) variable.
func SetWindowVar(window Window, name string, value interface{}) {
	name(nvim_win_set_var)
}

// DeleteWindowVar removes a window-scoped (w:) variable.
func DeleteWindowVar(window Window, name string) {
	name(nvim_win_del_var)
}

// WindowPosition gets the window position in display cells. First position is zero.
func WindowPosition(window Window) (pos [2]int) {
	name(nvim_win_get_position)
}

// WindowTabpage gets the window tabpage.
func WindowTabpage(window Window) (tabpage Tabpage) {
	name(nvim_win_get_tabpage)
}

// WindowNumber gets the window number.
func WindowNumber(window Window) (number int) {
	name(nvim_win_get_number)
}

// IsWindowValid checks if a window is valid.
func IsWindowValid(window Window) (valid bool) {
	name(nvim_win_is_valid)
}

// SetWindowConfig configure window position. Currently this is only used to configure
// floating and external windows (including changing a split window to these types).
//
// When reconfiguring a floating window, absent option keys will not be
// changed. "row"/"col" and "relative" must be reconfigured together.
//
// See documentation at OpenWindow, for the meaning of parameters.
func SetWindowConfig(window Window, config *WindowConfig) {
	name(nvim_win_set_config)
}

// WindowConfig return window configuration.
//
// The returned value may be given to OpenWindow.
//
// Relative will be an empty string for normal windows.
func WindowConfig(window Window) (config WindowConfig) {
	name(nvim_win_get_config)
	returnPtr()
}

// HideWindow closes the window and hide the buffer it contains (like ":hide" with a
// windowID).
//
// Like ":hide" the buffer becomes hidden unless another window is editing it,
// or "bufhidden" is "unload", "delete" or "wipe" as opposed to ":close" or
// CloseWindow, which will close the buffer.
func HideWindow(window Window) {
	name(nvim_win_hide)
}

// CloseWindow Closes the window (like ":close" with a window-ID).
func CloseWindow(window Window, force bool) {
	name(nvim_win_close)
}

// SetWindowHeightNamespace set highlight namespace for a window. This will use highlights defined in
// this namespace, but fall back to global highlights (ns=0) when missing.
//
// This takes predecence over the 'winhighlight' option.
func SetWindowHeightNamespace(window Window, nsID int) {
	name(nvim_win_set_hl_ns)
}
