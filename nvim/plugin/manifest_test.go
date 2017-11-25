package plugin

import (
	"bytes"
	"testing"
)

func TestOverwriteManifest(t *testing.T) {

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
		output := overwriteManifest(host, []byte(vimFileContent), []byte(manifest))
		expected := []byte(vimFileContent + manifest)
		if !bytes.Equal(output, expected) {
			t.Error("failed")
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
		output := overwriteManifest(host, []byte(vimFileContent), []byte(manifest))
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
			t.Error("failed")
		}
	})
}
