package plugin

import (
	"bytes"
	"testing"
)

func TestReplaceManifest(t *testing.T) {

	t.Run("not written manifest yet", func(t *testing.T) {
		vimFileContent := `if exists('g:loaded_hello')
  finish
endif
let g:loaded_hello = 1

function! s:RequireHello(host) abort
  return jobstart(['hello.nvim'], { 'rpc': v:true })
endfunction
`
		manifest := `call remote#host#RegisterPlugin('hello.nvim', '0', [
\ {'type': 'function', 'name': 'Ahello', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'Bhello', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},
\ ])
`
		host := "hello.nvim"
		output := replaceManifest(host, []byte(vimFileContent), []byte(manifest))
		expected := []byte(vimFileContent + manifest)
		if !bytes.Equal(output, expected) {
			t.Errorf("want %s, but get %s", string(expected), string(output))
		}
	})

	t.Run("already written manifest", func(t *testing.T) {
		vimFileContent := `if exists('g:loaded_hello')
  finish
endif
let g:loaded_hello = 1

function! s:RequireHello(host) abort
  return jobstart(['hello.nvim'], { 'rpc': v:true })
endfunction
call remote#host#RegisterPlugin('hello.nvim', '0', [
\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},
\ ])
`
		manifest := `call remote#host#RegisterPlugin('hello.nvim', '0', [
\ {'type': 'function', 'name': 'Ahello', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'Bhello', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},
\ ])
`
		host := "hello.nvim"
		output := replaceManifest(host, []byte(vimFileContent), []byte(manifest))
		expected := `if exists('g:loaded_hello')
  finish
endif
let g:loaded_hello = 1

function! s:RequireHello(host) abort
  return jobstart(['hello.nvim'], { 'rpc': v:true })
endfunction
call remote#host#RegisterPlugin('hello.nvim', '0', [
\ {'type': 'function', 'name': 'Ahello', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'Bhello', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'Hello', 'sync': 1, 'opts': {}},
\ ])
`
		if !bytes.Equal(output, []byte(expected)) {
			t.Errorf("want %s, but get %s", expected, string(output))
		}
	})
}
