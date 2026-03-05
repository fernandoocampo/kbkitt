# Contributing to kbcli

## Setup

1. Install git hooks:

```sh
make install-hooks
```

2. Verify tests and linter pass:

```sh
make test
make lint
```

## Claude Code Commands

This project includes custom [Claude Code](https://claude.com/claude-code) slash commands that automate the issue-to-PR workflow.

### `/resolve-issue <issue-number>`

Autonomously resolves a GitHub issue end-to-end. Example:

```
/resolve-issue 69
```

**What it does:**

1. Reads the issue from `fernandoocampo/kbkitt` (title, body, labels)
2. Determines the target app from `app=<name>` labels (defaults to `kbcli`)
3. Picks a branch type from labels (`enhancement` → `feat`, `bug` → `fix`, `documentation` → `docs`, `chore` → `chore`)
4. Creates an isolated worktree and branch (`<type>/<issue-number>`)
5. Follows TDD: writes failing tests → implements → refactors
6. Enforces quality gates: `make test`, `make lint`, `make coverage` (≥ 50%)
7. Commits using [Conventional Commits](#commit-conventions), pushes, and creates a PR

### `/pr-feedback <pr-number>`

Addresses review comments on an existing PR. Example:

```
/pr-feedback 42
```

**What it does:**

1. Reads the PR, review comments, and inline comments
2. Checks out the PR branch
3. Addresses each review comment, running `make test` + `make lint` after each change
4. Commits and pushes the fixes

## Commit Conventions

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>
```

- **Types:** `feat`, `fix`, `chore`, `docs`, `test`, `refactor`
- **Scope:** the app name (e.g., `kbcli`)
- **Description:** imperative mood, lowercase, no trailing period, max 72 chars
- **Footer:** `Resolves #<issue-number>` when closing an issue

Example:

```
feat(kbcli): add sync command for remote knowledge bases

Resolves #69
```

Git hooks must never be skipped (`--no-verify` is forbidden).

## Code Quality

| Gate | Command | Threshold |
|------|---------|-----------|
| Tests | `make test` | All pass |
| Lint | `make lint` | Zero warnings |
| Coverage | `make coverage` | ≥ 50% |
