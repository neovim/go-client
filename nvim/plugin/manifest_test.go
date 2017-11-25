package plugin

import (
	"bytes"
	"testing"
)

func TestReplaceManifest(t *testing.T) {

	t.Run("not written manifest yet", func(t *testing.T) {
		vimFileContentLines := []string{
			`if exists('g:loaded_hello')`,
			`  finish`,
			`endif`,
			`let g:loaded_hello = 1`,
			``,
			`function! s:RequireHello(host) abort`,
			`  return jobstart(['hello.nvim'], { 'rpc': v:true })`,
			`endfunction`,
		}
		var vimFileContent string
		for _, l := range vimFileContentLines {
			vimFileContent = vimFileContent + l + "\n"
		}

		manifestLines := []string{
			`call remote#host#RegisterPlugin('hello.nvim', '0', [`,
			`\ {'type': 'function', 'name': 'Ahello', 'sync': 1, 'opts': {}},`,
			`\ {'type': 'function', 'name': 'Bhello', 'sync': 1, 'opts': {}},`,
			`\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},`,
			`\ ])`,
		}
		var manifest string
		for _, l := range manifestLines {
			manifest = manifest + l + "\n"
		}

		host := "hello.nvim"
		output := replaceManifest(host, []byte(vimFileContent), []byte(manifest))
		expected := []byte(vimFileContent + manifest)
		if !bytes.Equal(output, expected) {
			t.Errorf("want %s, but get %s", string(expected), string(output))
		}
	})

	t.Run("already written manifest", func(t *testing.T) {
		vimFileContentLines := []string{
			`if exists('g:loaded_hello')`,
			`  finish`,
			`endif`,
			`let g:loaded_hello = 1`,
			``,
			`function! s:RequireHello(host) abort`,
			`  return jobstart(['hello.nvim'], { 'rpc': v:true })`,
			`endfunction`,
			`call remote#host#RegisterPlugin('hello.nvim', '0', [`,
			`\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},`,
			`\ ])`,
		}
		var vimFileContent string
		for _, l := range vimFileContentLines {
			vimFileContent = vimFileContent + l + "\n"
		}

		manifestLines := []string{
			`call remote#host#RegisterPlugin('hello.nvim', '0', [`,
			`\ {'type': 'function', 'name': 'Ahello', 'sync': 1, 'opts': {}},`,
			`\ {'type': 'function', 'name': 'Bhello', 'sync': 1, 'opts': {}},`,
			`\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},`,
			`\ ])`,
		}
		var manifest string
		for _, l := range manifestLines {
			manifest = manifest + l + "\n"
		}

		host := "hello.nvim"
		output := replaceManifest(host, []byte(vimFileContent), []byte(manifest))

		expectedLines := []string{
			`if exists('g:loaded_hello')`,
			`  finish`,
			`endif`,
			`let g:loaded_hello = 1`,
			``,
			`function! s:RequireHello(host) abort`,
			`  return jobstart(['hello.nvim'], { 'rpc': v:true })`,
			`endfunction`,
			`call remote#host#RegisterPlugin('hello.nvim', '0', [`,
			`\ {'type': 'function', 'name': 'Ahello', 'sync': 1, 'opts': {}},`,
			`\ {'type': 'function', 'name': 'Bhello', 'sync': 1, 'opts': {}},`,
			`\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},`,
			`\ ])`,
		}
		var expected string
		for _, l := range expectedLines {
			expected = expected + l + "\n"
		}

		if !bytes.Equal(output, []byte(expected)) {
			t.Errorf("want %s, but get %s", expected, string(output))
		}
	})
}
