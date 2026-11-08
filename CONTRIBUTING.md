# Contributing to pigment

Thanks for contributing.

## Development requirements

- Go 1.21 or newer
- No external runtime dependencies

Before opening a pull request, run:

```sh
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

For parser changes, also run:

```sh
go test -fuzz=FuzzRenderNeverPanics -fuzztime=30s ./...
go test -run '^$' -bench BenchmarkRender ./...
```

Keep public API changes documented in `README.md` and `doc.go`. Add focused
tests for new behavior and bug fixes. Avoid unrelated refactors in the same
change.

## Pull requests

Describe the user-visible behavior, test commands, and any compatibility
impact. Pull requests must pass the GitHub Actions checks before merging.
