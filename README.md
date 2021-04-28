# Go Build Tools

[![Go Report Card](https://goreportcard.com/badge/github.com/forewing/gobuild?style=flat-square)](https://goreportcard.com/report/github.com/forewing/gobuild)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/forewing/gobuild?style=flat-square)](https://github.com/forewing/gobuild/releases/latest)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/forewing/gobuild)](https://pkg.go.dev/github.com/forewing/gobuild)

Tools for building and distributing Go executables

## Example

[example/purego](./example/purego)

```
$ cd example/purego
$ go run .
$ tar xaf output/purego-XXX-linux-amd64.tar.gz
$ ./purego-XXX-linux-amd64
hello world
version: XXX
hash: YYY
```
