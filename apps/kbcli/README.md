# kbcli

`kbcli` is the CLI application for **kbkitt** — a knowledge base management system. It lets you store, search, and retrieve snippets of knowledge (commands, concepts, prompts, bookmarks, quotes, etc.) locally via SQLite with optional synchronization to a central kbkitt server.

## Table of Contents

- [How to run?](#how-to-run)
- [How to test?](#how-to-test)
- [How to setup?](#how-to-setup)
- [Commands](#commands)
  - [add](#add)
  - [get](#get)
  - [update](#update)
  - [import](#import)
  - [export](#export)
  - [sync](#sync)
  - [version](#version)
- [Knowledge Base Data Model](#knowledge-base-data-model)
- [Interactive UI Keyboard Shortcuts](#interactive-ui-keyboard-shortcuts)
- [Media Management](#media-management)
- [Technology Stack](#technology-stack)
- [References](#references)

---

## How to run?

* Show help by default.

```sh
make run
```

* Or using go command.

```sh
go run cmd/kbcli/main.go
```

## How to test?

```sh
make test
```

---

## How to setup?

Run `configure` command to setup your environment.

```sh
kbkitt configure
```

You will be asked to configure kbkitt in your machine, typically in `{HOME_DIR}/.kbkitt`.

```sh
.
├── config.yaml
├── kbkitt.db
├── media
│   └── btc.jpeg
└── sync.yaml
```

`config.yaml` contains basic information about where the remote server is, in case you want to centralize your kbs.

```yaml
version: 0.1.0
fileForSyncPath: {HOME_DIR}/.kbkitt/sync.yaml
dirForMediaPath: {HOME_DIR}/.kbkitt/media
server:
    url: http://localhost:3030
```

* `fileForSyncPath` — file that keeps KBs you could not send to the central server (offline queue).
* `dirForMediaPath` — directory that stores media resources (images, docs, videos, etc.) saved as KBs.
* `server.url` — kbkitt remote server URL.
* `kbkitt.db` — local SQLite database with full-text search support.

---

## How to contribute?

1. install githooks first

```sh
make install-githooks
```

## Commands

### add

Add a new knowledge base entry.

```sh
kbkitt add --help

add a new knowledge base such as: concepts, commands, prompts, etc.

Usage:
  kb add [flags]

Flags:
  -c, --category string    category of knowledge base
  -h, --help               help for add
  -k, --key string         knowledge base key
  -n, --namespace string   namespace of knowledge base (default: "default")
  -o, --notes string       knowledge base notes
  -r, --reference string   author or reference of this kb
  -t, --tags strings       comma separated tags for this kb
  -u, --ux                 add KB in interactive mode
  -v, --value string       knowledge base value
```

The application will prompt for any parameters you did not provide.

```sh
kbkitt add
```

```yaml
key:
value:
notes:
category:
reference:
tags:
```

For an improved interactive GUI, use the `-u`/`--ux` flag.

```sh
kbkitt add --ux
```

```sh
 Adding a new KB:

Key

Category

Value
...

Notes

Reference

Tags
keyword1 keyword2 keyword3 keywordN

Continue ->

• tab fields • shift+tab fields • ctrl+c: quit
```

To add a KB in a single command:

```sh
kbkitt add \
  -k btc -v crypto -n currencies -c crypto \
  -t btc,crypto,currencies,blockchain \
  -r dementor

Key: btc
Value: crypto
Notes: currencies
Category: crypto
Reference: dementor
Tags: [btc crypto currencies blockchain]

> do you want to save it? [y/n]:
```

---

### get

Search and retrieve knowledge bases.

```sh
kbkitt get --help

Usage:
  kb get [flags]

Flags:
  -c, --category string    filter by category
  -h, --help               help for get
  -i, --id string          knowledge base id
  -k, --key string         filter by key
  -l, --limit int          max number of results (default 5)
  -n, --namespace string   filter by namespace
  -o, --offset int         pagination offset (default 0)
      --random-quote       get a random KB from the "quote" category
  -w, --keyword string     search by keyword (full-text search on tags)
```

**Basic search:**

```sh
# Search by key
kbkitt get -k btc

# Search by keyword (full-text search on tags)
kbkitt get -w blockchain

# Filter by category and namespace
kbkitt get -c crypto -n default

# Get a random quote
kbkitt get --random-quote

# Paginate results
kbkitt get -c crypto -l 10 -o 0
```

**Interactive search UI:**

Running `kbkitt get` without flags launches an interactive TUI with:
- A filter panel (toggle with `Ctrl+F`) to set category, namespace, key, and keyword
- A results table with pagination
- A detail viewer for the selected KB with markdown rendering

---

### update

Update an existing knowledge base entry.

```sh
kbkitt update --help

Usage:
  kb update [flags]

Flags:
  -h, --help        help for update
  -i, --id string   knowledge base id to update
  -u, --ux          update KB in interactive mode
```

```sh
# Update by ID in interactive mode
kbkitt update -i <kb-id> --ux
```

---

### import

Import knowledge bases from a YAML file.

```sh
kbkitt import --help

Usage:
  kb import [flags]

Flags:
  -f, --file string       path to YAML file to import
  -h, --help              help for import
      --show-added-kbs    print successfully imported KBs
      --show-failed-kbs   print KBs that failed to import
```

```sh
kbkitt import -f my-kbs.yaml --show-added-kbs --show-failed-kbs
```

The YAML file format matches the export format, making it easy to move KBs between environments.

---

### export

Export knowledge bases to YAML format (stdout).

```sh
kbkitt export --help

Usage:
  kb export [flags]

Flags:
  -c, --category string    filter by category
  -h, --help               help for export
  -n, --namespace string   filter by namespace
```

```sh
# Export all KBs
kbkitt export

# Export a specific category
kbkitt export -c crypto

# Export and save to file
kbkitt export -c crypto > crypto-kbs.yaml
```

---

### sync

Sync locally saved KBs to the central server.

When the kbkitt server is unreachable, new KBs are saved to `~/.kbkitt/sync.yaml`. The `sync` command retries sending those queued entries to the server.

```sh
kbkitt sync --help

Usage:
  kb sync [flags]

Flags:
  -h, --help              help for sync
      --show-added-kbs    print successfully synced KBs
      --show-failed-kbs   print KBs that failed to sync
```

```sh
kbkitt sync --show-added-kbs
```

---

### version

Display build version information.

```sh
kbkitt version

version: 0.1.0
commit: ab78de2
built at: 2024-01-01T00:00:00Z
```

---

## Knowledge Base Data Model

Each KB entry has the following fields:

| Field | Description | Constraints |
|-------|-------------|-------------|
| `id` | Unique identifier (UUID) | Auto-generated |
| `key` | Short identifier | Lowercase, alphanumeric |
| `value` | Main content | Up to 700 characters |
| `notes` | Additional notes | Up to 700 characters |
| `category` | Classification | Lowercase (e.g., `quote`, `media`, `bookmark`, `command`) |
| `namespace` | Organization scope | Lowercase, default: `default` |
| `reference` | Author or source attribution | Free text |
| `tags` | Search keywords | Alphanumeric + hyphens, deduplicated and sorted |

**Tags** power the full-text search — use descriptive tags to make KBs easy to find later.

---

## Interactive UI Keyboard Shortcuts

### Get / Search Mode

| Shortcut | Action |
|----------|--------|
| `Ctrl+F` | Toggle filter panel |
| `Ctrl+C` | Copy selected KB value to clipboard |
| `Ctrl+O` | Open selected KB URL in browser (bookmarks) |
| `Ctrl+R` | Return to results table from detail view |
| `↑ / ↓` | Navigate rows in results table |
| `← / →` | Previous / next page of results |
| `Enter` | View selected KB detail |
| `Esc / Ctrl+Q` | Quit |

### Add / Update Mode

| Shortcut | Action |
|----------|--------|
| `Tab / Ctrl+N` | Next field |
| `Shift+Tab / Ctrl+P` | Previous field |
| `Enter` (on Continue) | Submit form |
| `Ctrl+C` | Quit |

---

## Media Management

You can save media resources (images, documents, videos) as KB entries. Files are downloaded and stored in `~/.kbkitt/media/`.

Supported media types: `APNG`, `AVIF`, `CSV`, `GIF`, `JPEG`, `MP4`, `PDF`, `PNG`, `SVG`, `TAR.GZ`, `TXT`, `WEBP`, `YAML`, `ZIP`.

---

## Technology Stack

| Component | Library |
|-----------|---------|
| CLI framework | [Cobra](https://github.com/spf13/cobra) |
| TUI framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| TUI components | [Bubbles](https://github.com/charmbracelet/bubbles) |
| TUI styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| Markdown rendering | [Glamour](https://github.com/charmbracelet/glamour) |
| Local database | SQLite (FTS5 full-text search) |
| Configuration | YAML |
| ID generation | Google UUID |
| Clipboard | golang.design/x/clipboard |

---

## References

* bubbletea examples
https://github.com/charmbracelet/bubbletea/blob/master/examples
