# toc

A simple table of contents generator for markdowns written in golang

## Contents
1. [Install](#install)
2. [Usage](#usage)
    1. [Output to stdin](#output-to-stdin)
    2. [Inplace update](#inplace-update)
    3. [Custom heading IDs](#custom)

## Install

```sh
$ go get github.com/tvastar/toc/toc.go
```

## Usage

The table of contents of this README.md is itself maintaained using this tool.

### Output to stdin

```sh
$ toc $GOPATH/github.com/tvastar/toc/README.md
```sh

### Inplace update

```sh
$ cd $GOPATH/github.com/tvastar/toc
$ toc -o README.md -h Contents README.md
```sh

### Custom heading IDs {#custom}

As with regular markdown, a custom heading ID may be specified via `{#id}` suffix to the heading line.  This is properly parsed and used for headings.

# Ignored

Top level sections are ignored by default
