# µC Compiler

[![Build Status](https://travis-ci.org/mewmew/uc.svg?branch=dev)](https://travis-ci.org/mewmew/uc)
[![Coverage Status](https://coveralls.io/repos/github/mewmew/uc/badge.svg?branch=dev)](https://coveralls.io/github/mewmew/uc?branch=dev)
[![GoDoc](https://godoc.org/github.com/mewmew/uc?status.svg)](https://godoc.org/github.com/mewmew/uc)

A compiler for the [µC programming language](https://www.it.uu.se/katalog/aleji304/CompilersProject/uc.html).

## Installation

1. [Install Go](https://golang.org/doc/install).
2. Install Gocc `go get github.com/goccmack/gocc`.

```
$ go get -u github.com/mewmew/uc
$ cd ${GOPATH}/src/github.com/mewmew/uc/gocc
$ make gen
$ go get github.com/mewmew/uc/...
$ go install github.com/mewmew/uc/cmd/ulex
$ go install github.com/mewmew/uc/cmd/uparse
$ go install github.com/mewmew/uc/cmd/uclang
```

## Usage

* [ulex](https://godoc.org/github.com/mewmew/uc/cmd/ulex): a lexer for the µC language which pretty-prints tokens to standard output.
* [uparse](https://godoc.org/github.com/mewmew/uc/cmd/uparse): a parser for the µC language which pretty-prints abstract syntax trees to standard output.
* [usem](https://godoc.org/github.com/mewmew/uc/cmd/usem): a static semantic checker for the µC language which validates the input and reports errors to standard error.
* [uclang](https://godoc.org/github.com/mewmew/uc/cmd/uclang): a compiler for the µC language which validates the input, generates corresponding LLVM IR assembly and reports errors to standard error.

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
