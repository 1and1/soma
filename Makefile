# vim: set ft=make ffs=unix fenc=utf8:
# vim: set noet ts=4 sw=4 tw=72 list:
#
SOMAVER != cat `git rev-parse --show-toplevel`/VERSION
BRANCH != git rev-parse --symbolic-full-name --abbrev-ref HEAD
GITHASH != git rev-parse --short HEAD

all: install

install: install_freebsd install_linux

install_freebsd: generate
	@env GOOS=freebsd GOARCH=amd64 go install -ldflags "-X main.somaVersion=$(SOMAVER)-$(GITHASH)/$(BRANCH)" ./...

install_linux: generate
	@env GOOS=linux GOARCH=amd64 go install -ldflags "-X main.somaVersion=$(SOMAVER)-$(GITHASH)/$(BRANCH)" ./...

generate:
	@go generate ./cmd/...
