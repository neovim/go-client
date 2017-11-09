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

// This file defines the Nvim remote API using Go syntax. Run the 'go generate'
// command to convert this file to the API implementation in apiimp.go.
//
// The go generate command runs the apitool program. Run
//
//  go run apitool.go --help
//
// to learn more about the apitool program.

package main

// BufferLineCount returns the number of lines in the buffer.
func BufferLineCount(buffer Buffer) int {
	name(nvim_buf_line_count)
}

// BufferLines retrieves a line range from a buffer.
//
// Indexing is zero-based, end-exclusive. Negative indices are interpreted as
// length+1+index, i e -1 refers to the index past the end. So to get the last
// element set start=-2 and end=-1.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strict
// = true.
func BufferLines(buffer Buffer, start int, end int, strict bool) [][]byte {
	name(nvim_buf_get_lines)
}

// SetBufferLines replaces a line range on a buffer.
//
// Indexing is zero-based, end-exclusive. Negative indices are interpreted as
// length+1+index, ie -1 refers to the index past the end. So to change or
// delete the last element set start=-2 and end=-1.
//
// To insert lines at a given index, set both start and end to the same index.
// To delete a range of lines, set replacement to an empty array.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strict
// = true.
func SetBufferLines(buffer Buffer, start int, end int, strict bool, replacement [][]byte) {
	name(nvim_buf_set_lines)
}

// BufferVar gets a buffer-scoped (b:) variable.
func BufferVar(buffer Buffer, name string) interface{} {
	name(nvim_buf_get_var)
}

// BufferChangedTick gets a changed tick of a buffer.
func BufferChangedTick(buffer Buffer) int {
	name(nvim_buf_get_changedtick)
}

// SetBufferVar sets a buffer-scoped (b:) variable.
func SetBufferVar(buffer Buffer, name string, value interface{}) {
	name(nvim_buf_set_var)
}

// DeleteBufferVar removes a buffer-scoped (b:) variable.
func DeleteBufferVar(buffer Buffer, name string) {
	name(nvim_buf_del_var)
}

// BufferOption gets a buffer option value.
func BufferOption(buffer Buffer, name string) interface{} {
	name(nvim_buf_get_option)
}

// SetBufferOption sets a buffer option value. The value nil deletes the option
// in the case where there's a global fallback.
func SetBufferOption(buffer Buffer, name string, value interface{}) {
	name(nvim_buf_set_option)
}

// BufferNumber gets a buffer's number.
//
// Deprecated: use int(buffer) to get the buffer's number as an integer.
func BufferNumber(buffer Buffer) int {
	name(nvim_buf_get_number)
	deprecatedSince(2)
}

// BufferName gets the full file name of a buffer.
func BufferName(buffer Buffer) string {
	name(nvim_buf_get_name)
}

// SetBufferName sets the full file name of a buffer.
// BufFilePre/BufFilePost are triggered.
func SetBufferName(buffer Buffer, name string) {
	name(nvim_buf_set_name)
}

// IsBufferValid returns true if the buffer is valid.
func IsBufferValid(buffer Buffer) bool {
	name(nvim_buf_is_valid)
}

// BufferMark returns the (row,col) of the named mark.
func BufferMark(buffer Buffer, name string) [2]int {
	name(nvim_buf_get_mark)
}

// AddBufferHighlight adds a highlight to buffer and returns the source id of
// the highlight.
//
// AddBufferHighlight can be used for plugins which dynamically generate
// highlights to a buffer (like a semantic highlighter or linter). The function
// adds a single highlight to a buffer. Unlike matchaddpos() highlights follow
// changes to line numbering (as lines are inserted/removed above the
// highlighted line), like signs and marks do.
//
// The srcID is useful for batch deletion/updating of a set of highlights. When
// called with srcID = 0, an unique source id is generated and returned.
// Successive calls can pass in it as srcID to add new highlights to the same
// source group. All highlights in the same group can then be cleared with
// ClearBufferHighlight. If the highlight never will be manually deleted pass
// in -1 for srcID.
//
// If hlGroup is the empty string no highlight is added, but a new srcID is
// still returned. This is useful for an external plugin to synchrounously
// request an unique srcID at initialization, and later asynchronously add and
// clear highlights in response to buffer changes.
//
// The startCol and endCol parameters specify the range of columns to
// highlight. Use endCol = -1 to highlight to the end of the line.
func AddBufferHighlight(buffer Buffer, srcID int, hlGroup string, line int, startCol int, endCol int) int {
	name(nvim_buf_add_highlight)
}

// ClearBufferHighlight clears highlights from a given source group and a range
// of lines.
//
// To clear a source group in the entire buffer, pass in 1 and -1 to startLine
// and endLine respectively.
//
// The lineStart and lineEnd parameters specify the range of lines to clear.
// The end of range is exclusive. Specify -1 to clear to the end of the file.
func ClearBufferHighlight(buffer Buffer, srcID int, startLine int, endLine int) {
	name(nvim_buf_clear_highlight)
}

// TabpageWindows returns the windows in a tabpage.
func TabpageWindows(tabpage Tabpage) []Window {
	name(nvim_tabpage_list_wins)
}

// TabpageVar gets a tab-scoped (t:) variable.
func TabpageVar(tabpage Tabpage, name string) interface{} {
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

// TabpageWindow gets the current window in a tab page.
func TabpageWindow(tabpage Tabpage) Window {
	name(nvim_tabpage_get_win)
}

// TabpageNumber gets the tabpage number from the tabpage handle.
func TabpageNumber(tabpage Tabpage) int {
	name(nvim_tabpage_get_number)
}

// IsTabpageValid checks if a tab page is valid.
func IsTabpageValid(tabpage Tabpage) bool {
	name(nvim_tabpage_is_valid)
}

// AttachUI registers the client as a remote UI. After this method is called,
// the client will receive redraw notifications.
//
//  :help rpc-remote-ui
//
// The redraw notification method has variadic arguments. Register a handler
// for the method like this:
//
//  v.RegisterHandler("redraw", func(updates ...[]interface{}) {
//      for _, update := range updates {
//          // handle update
//      }
//  })
func AttachUI(width int, height int, options map[string]interface{}) {
	name(nvim_ui_attach)
}

// DetachUI unregisters the client as a remote UI.
func DetachUI() {
	name(nvim_ui_detach)
}

// TryResizeUI notifies Nvim that the client window has resized. If possible,
// Nvim will send a redraw request to resize.
func TryResizeUI(width int, height int) {
	name(nvim_ui_try_resize)
}

// SetUIOption sets a UI option.
func SetUIOption(name string, value interface{}) {
	name(nvim_ui_set_option)
}

// Command executes a single ex command.
func Command(cmd string) {
	name(nvim_command)
}

// FeedKeys Pushes keys to the Nvim user input buffer. Options can be a string
// with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//     mapping. This matters for undo, opening folds, etc.
func FeedKeys(keys string, mode string, escapeCSI bool) {
	name(nvim_feedkeys)
}

// Input pushes bytes to the Nvim low level input buffer.
//
// Unlike FeedKeys, this uses the lowest level input buffer and the call is not
// deferred. It returns the number of bytes actually written(which can be less
// than what was requested if the buffer is full).
func Input(keys string) int {
	name(nvim_input)
}

// ReplaceTermcodes replaces any terminal code strings by byte sequences. The
// returned sequences are Nvim's internal representation of keys, for example:
//
//  <esc> -> '\x1b'
//  <cr>  -> '\r'
//  <c-l> -> '\x0c'
//  <up>  -> '\x80ku'
//
// The returned sequences can be used as input to feedkeys.
func ReplaceTermcodes(str string, fromPart bool, doLT bool, special bool) string {
	name(nvim_replace_termcodes)
}

// CommandOutput executes a single ex command and returns the output.
func CommandOutput(cmd string) string {
	name(nvim_command_output)
}

// Eval evaluates the expression expr using the Vim internal expression
// evaluator.
//
//  :help expression
func Eval(expr string) interface{} {
	name(nvim_eval)
}

// StringWidth returns the number of display cells the string occupies. Tab is
// counted as one cell.
func StringWidth(s string) int {
	name(nvim_strwidth)
}

// RuntimePaths returns a list of paths contained in the runtimepath option.
func RuntimePaths() []string {
	name(nvim_list_runtime_paths)
}

// SetCurrentDirectory changes the Vim working directory.
func SetCurrentDirectory(dir string) {
	name(nvim_set_current_dir)
}

// CurrentLine gets the current line in the current buffer.
func CurrentLine() []byte {
	name(nvim_get_current_line)
}

// SetCurrentLine sets the current line in the current buffer.
func SetCurrentLine(line []byte) {
	name(nvim_set_current_line)
}

// DeleteCurrentLine deletes the current line in the current buffer.
func DeleteCurrentLine() {
	name(nvim_del_current_line)
}

// Var gets a global (g:) variable.
func Var(name string) interface{} {
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

// VVar gets a vim (v:) variable.
func VVar(name string) interface{} {
	name(nvim_get_vvar)
}

// Option gets an option.
func Option(name string) interface{} {
	name(nvim_get_option)
}

// SetOption sets an option.
func SetOption(name string, value interface{}) {
	name(nvim_set_option)
}

// WriteOut writes a message to vim output buffer. The string is split and
// flushed after each newline. Incomplete lines are kept for writing later.
func WriteOut(str string) {
	name(nvim_out_write)
}

// WriteErr writes a message to vim error buffer. The string is split and
// flushed after each newline. Incomplete lines are kept for writing later.
func WriteErr(str string) {
	name(nvim_err_write)
}

// WritelnErr writes prints str and a newline as an error message.
func WritelnErr(str string) {
	name(nvim_err_writeln)
}

// Buffers returns the current list of buffers.
func Buffers() []Buffer {
	name(nvim_list_bufs)
}

// CurrentBuffer returns the current buffer.
func CurrentBuffer() Buffer {
	name(nvim_get_current_buf)
}

// SetCurrentBuffer sets the current buffer.
func SetCurrentBuffer(buffer Buffer) {
	name(nvim_set_current_buf)
}

// Windows returns the current list of windows.
func Windows() []Window {
	name(nvim_list_wins)
}

// CurrentWindow returns the current window.
func CurrentWindow() Window {
	name(nvim_get_current_win)
}

// SetCurrentWindow sets the current window.
func SetCurrentWindow(window Window) {
	name(nvim_set_current_win)
}

// Tabpages returns the current list of tabpages.
func Tabpages() []Tabpage {
	name(nvim_list_tabpages)
}

// CurrentTabpage returns the current tabpage.
func CurrentTabpage() Tabpage {
	name(nvim_get_current_tabpage)
}

// SetCurrentTabpage sets the current tabpage.
func SetCurrentTabpage(tabpage Tabpage) {
	name(nvim_set_current_tabpage)
}

// Subscribe subscribes to a Nvim event.
func Subscribe(event string) {
	name(nvim_subscribe)
}

// Unsubscribe unsubscribes to a Nvim event.
func Unsubscribe(event string) {
	name(nvim_unsubscribe)
}

func ColorByName(name string) int {
	name(nvim_get_color_by_name)
}

func ColorMap() map[string]interface{} {
	// TODO: should this functino return map[string]int?
	name(nvim_get_color_map)
}

// Mode gets Nvim's current mode.
func Mode() Mode {
	name(nvim_get_mode)
	returnPtr()
}

func APIInfo() []interface{} {
	name(nvim_get_api_info)
}

// WindowBuffer returns the current buffer in a window.
func WindowBuffer(window Window) Buffer {
	name(nvim_win_get_buf)
}

// WindowCursor returns the cursor position in the window.
func WindowCursor(window Window) [2]int {
	name(nvim_win_get_cursor)
}

// SetWindowCursor sets the cursor position in the window to the given position.
func SetWindowCursor(window Window, pos [2]int) {
	name(nvim_win_set_cursor)
}

// WindowHeight returns the window height.
func WindowHeight(window Window) int {
	name(nvim_win_get_height)
}

// SetWindowHeight sets the window height.
func SetWindowHeight(window Window, height int) {
	name(nvim_win_set_height)
}

// WindowWidth returns the window width.
func WindowWidth(window Window) int {
	name(nvim_win_get_width)
}

// SetWindowWidth sets the window width.
func SetWindowWidth(window Window, width int) {
	name(nvim_win_set_width)
}

// WindowVar gets a window-scoped (w:) variable.
func WindowVar(window Window, name string) interface{} {
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

// WindowOption gets a window option.
func WindowOption(window Window, name string) interface{} {
	name(nvim_win_get_option)
}

// SetWindowOption sets a window option.
func SetWindowOption(window Window, name string, value interface{}) {
	name(nvim_win_set_option)
}

// WindowPosition gets the window position in display cells. First position is zero.
func WindowPosition(window Window) [2]int {
	name(nvim_win_get_position)
}

// WindowTabpage gets the tab page that contains the window.
func WindowTabpage(window Window) Tabpage {
	name(nvim_win_get_tabpage)
}

// WindowNumber gets the window number from the window handle.
func WindowNumber(window Window) int {
	name(nvim_win_get_number)
}

// IsWindowValid returns true if the window is valid.
func IsWindowValid(window Window) bool {
	name(nvim_win_is_valid)
}
