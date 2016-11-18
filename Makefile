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

sanitize: build vet ineffassign misspell

build:
	@go build ./...

vet:
	@go vet ./cmd/eye/
	@go vet ./cmd/soma/
	@go vet ./cmd/somaadm/
	@go vet ./cmd/somadbctl/
	@go vet ./lib/auth/
	@go vet ./lib/proto/
	@go vet ./internal/adm/
	@go vet ./internal/cmpl/
	@go vet ./internal/db/
	@go vet ./internal/help/
	@go vet ./internal/msg/
	@go vet ./internal/stmt/
	@go vet ./internal/tree/
	@go tool vet -shadow ./cmd/eye/
	@go tool vet -shadow ./cmd/soma/
	@go tool vet -shadow ./cmd/somaadm/
	@go tool vet -shadow ./cmd/somadbctl/
	@go tool vet -shadow ./lib/auth/
	@go tool vet -shadow ./lib/proto/
	@go tool vet -shadow ./internal/adm/
	@go tool vet -shadow ./internal/cmpl/
	@go tool vet -shadow ./internal/db/
	@go tool vet -shadow ./internal/help/
	@go tool vet -shadow ./internal/msg/
	@go tool vet -shadow ./internal/stmt/
	@go tool vet -shadow ./internal/tree/

ineffassign:
	@ineffassign ./cmd
	@ineffassign ./lib
	@ineffassign ./internal

misspell:
	@misspell ./cmd
	@misspell ./lib
	@misspell ./internal
	@misspell ./docs
