all: validate

validate:
	@go build ./...
	@go vet .
	@go tool vet -shadow .
	@golint .
	@ineffassign .
