# kbcli

CLI tool for knowledge base management. Uses Cobra for commands, Bubbletea for interactive TUI, and SQLite for local storage.

## Commands

Run tests:

    make test
    # Runs: go test -race -count=1 ./...

Clean build artifacts:

    go clean ./...

Tidy modules:

    make mod-tidy

Build (macOS):

    make build-macos-amd-64

## Go Conventions

- Follow standard Go formatting: `gofmt`/`goimports`
- Wrap errors with context: `fmt.Errorf("doing X: %w", err)`
- Use table-driven tests with `t.Run`
- Define interfaces in the consumer package, not the provider
- Avoid global state; pass dependencies explicitly
