# What is SOMA?

SOMA is a configuration service for monitoring systems. It is currently
in public beta-ish. The next release is targeted to contain the full
documentation at which point it will be a lot easier to use.

# Build instructions

```
  go generate ./cmd/...
  go build ./...
  go install ./...
```

# Notes

* the `go generate` phase requires `go-bindata` to be installed and in
  $PATH. It is available at https://github.com/jteeuwen/go-bindata
* running `go generate ./...` will also run the generate stages inside
  `vendor/`. This may or may not be what you intended.
