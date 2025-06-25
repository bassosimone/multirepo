# multirepo: Manage repositories as a monorepo

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/multirepo)](https://pkg.go.dev/github.com/bassosimone/multirepo) [![Build Status](https://github.com/bassosimone/multirepo/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/multirepo/actions) [![codecov](https://codecov.io/gh/bassosimone/multirepo/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/multirepo)

This repository implements a tool to manage several repositories
as they were a single repository. This tool is still in the early
stages of development and is not yet ready for production use.

## Design

See [DESIGN.md](DESIGN.md).

## Examples

Creating a multirepo in the current directory:

```bash
multirepo init
```

Cloning a repository within the multirepo:

```bash
multirepo clone git@github.com:rbmk-project/rbmk
```

Executing a command within the multirepo:

```bash
multirepo foreach git status -v
```

Getting interactive help:

```bash
multirepo --help
```

## Minimum Supported Go Version

Go 1.24

## Installation

```bash
go install github.com/bassosimone/multirepo@latest
```

## Running Tests

You need GNU make installed.

```
make check
```

## Compiling for the current system

You need GNU make installed.

```
make multirepo
```

## Building a release

You need GNU make installed.

```
make release
```

## Dependencies

- [github.com/bassosimone/clip](https://pkg.go.dev/github.com/bassosimone/clip)

- [github.com/charmbracelet/lipgloss](https://pkg.go.dev/github.com/charmbracelet/lipgloss)

- [github.com/kballard/go-shellquote](https://pkg.go.dev/github.com/kballard/go-shellquote)

- [github.com/rogpeppe/go-internal](https://pkg.go.dev/github.com/rogpeppe/go-internal)

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```
