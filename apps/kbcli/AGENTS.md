# kbcli

CLI tool for knowledge base management. Uses Cobra for commands, Bubbletea for interactive TUI, and SQLite for local storage.

## Build & Run

### Prerequisites

- Go: 1.25+ (check individual service `go.mod` for specific versions)
- Docker & Docker Compose
- Make

### Quality Gates

Run tests:

    ```bash
    make test
    # Runs: go test -race -count=1 ./...
    ```

Tidy modules:

    ```bash
    make mod-tidy
    ```

Lint:

    ```bash
    make lint
    ```

Build (macOS):

    make build-macos-amd-64

## Coding Practices

### Version Control

- Use Conventional Commits format
- Branch naming: `feat/`, `fix/`, `chore/`

### Go Conventions

- Follow standard Go formatting: `gofmt`/`goimports`
- Wrap errors with context: `fmt.Errorf("doing X: %w", err)`
- Use table-driven tests with `t.Run`
- Define interfaces in the consumer package, not the provider
- Avoid global state; pass dependencies explicitly

### Testing

- Unit tests required for all new code (co-located with source: `*_test.go`)
- Use table-driven tests where appropriate