# go-client

[![pkg.go.dev][pkg.go.dev-badge]][pkg.go.dev]
[![Github Actions][Github Actions Badge]][Github Actions]
[![codecov.io][codecov-badge]][codecov]

Neovim/go-client is a [Neovim](https://neovim.io/) client for [Go](https://golang.org/).

Update API
----------

The API can be updated with the api tool.

```sh
cd nvim

# The compare command indicates what functions need to be changed. The results
# should be added to `api_def.go`.
go run api_tool.go -compare

# After applying the chages from compare to `api_def.go`, generate the `api.go`
# and `api_deprecated.go` files.
go generate
```

Release
-------

Follow the standard Go process for [publishing a module](https://go.dev/doc/modules/publishing). 

License
-------

[Apache License 2.0](https://github.com/neovim/go-client/blob/master/LICENSE)

<!-- badge links -->
[pkg.go.dev]: https://pkg.go.dev/github.com/neovim/go-client
[Github Actions]: https://github.com/neovim/go-client/actions
[codecov]: https://app.codecov.io/gh/neovim/go-client

[pkg.go.dev-badge]: https://pkg.go.dev/badge/github.com/neovim/go-client.svg
[Github Actions Badge]: https://img.shields.io/github/workflow/status/go-clang/gen/Test/main?label=test&logo=github&style=flat-square
[codecov-badge]: https://img.shields.io/codecov/c/github/neovim/go-client/master?logo=codecov&style=flat-square
