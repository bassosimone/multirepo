# SPDX-License-Identifier: GPL-3.0-or-later

#doc:
#doc: usage: make [target]
#doc:
#doc: We support the following targets:
.PHONY: help
help:
	@cat GNUmakefile | grep -E '^#doc:' | sed -e 's/^#doc: //g' -e 's/^#doc://'

#doc:
#doc: - `all`: builds `multiple` for current platform
.PHONY: all
all: multirepo

# Common variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.Version=$(VERSION)
TAGS := netgo
GOENV := CGO_ENABLED=0
EXE ?=

#doc:
#doc: - `multirepo`: build multirepo in the current directory
#doc:
#doc: Use GOOS and GOARCH to force a specific architecture. For example, to
#doc: build for Linux on an ARM64 machine, run:
#doc:
#doc:     GOOS=linux GOARCH=arm64 make multirepo
#doc:
#doc: The resulting binary will be named `multirepo-linux-arm64`.
.PHONY: multirepo
multirepo:
	$(GOENV) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o multirepo-$(GOOS)-$(GOARCH)${EXE} -ldflags '$(LDFLAGS)' -tags $(TAGS) .

#doc:
#doc: - `release`: cross compile all platform/variant combinations
.PHONY: release
release:
	GOOS=linux GOARCH=amd64 $(MAKE) multirepo
	GOOS=linux GOARCH=arm64 $(MAKE) multirepo
	GOOS=windows GOARCH=amd64 EXE=.exe $(MAKE) multirepo
	GOOS=darwin GOARCH=arm64 $(MAKE) multirepo

#doc:
#doc: - `check`: run tests
.PHONY: check
check:
	go test -race -count 1 -cover ./...

#doc:
#doc: - `clean`: remove build artifacts
.PHONY: clean
clean:
	rm -f multirepo multirepo-*

#doc:
#doc: - `install`: install multirepo into the system
#doc:
#doc: Installs multirepo for the current platform.
#doc: Use PREFIX to specify installation prefix (default: `/usr/local`).
#doc: For staged installations, use DESTDIR as usual.
#doc:
#doc: Examples:
#doc:     make install
#doc:     make PREFIX=/opt/multirepo install
#doc:     make DESTDIR=/tmp/stage PREFIX=/usr/local install
.PHONY: install
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

install: multirepo
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 multirepo-$(GOOS)-$(GOARCH) $(DESTDIR)$(BINDIR)/multirepo

#doc:
