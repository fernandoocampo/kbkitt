GO_TOOL?=go
GO_BUILD=$(GO_TOOL) build
ADD_ARGS?=
GET_ARGS?=
ROOT_ARGS?=
KBNAMESPACE?=default
COMMIT_HASH?=$(shell git describe --dirty --tags --always)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
IMPORT_FILE?=/tmp/sample.yaml
BINARY_NAME=kbkitt
BINARY_UNIX=$(BINARY_NAME)-linux
BINARY_DARWIN=$(BINARY_NAME)-darwin
VERSION?=0.1.0
LDFLAGS?="-X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions.Version=${VERSION} -X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions.CommitHash=${COMMIT_HASH} -X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions.BuildDate=${BUILD_DATE} -s -w"

.PHONY: help
help: ## Display this help.
		@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/{ printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run unit tests.
	$(GO_TOOL) test -race -count=1 ./...

.PHONY: run
run: ## Run run app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go $(ROOT_ARGS)

.PHONY: build-macos-amd-64
build-macos-amd-64: ## Build binary for Linux amd64
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${GO_BUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN}-amd64 ./cmd/kbcli/main.go

.PHONY: run-help
run-help: ## Run help app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go help

.PHONY: run-version
run-version: ## Run version app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go version

.PHONY: run-conf
run-conf: ## Run configure app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go configure

.PHONY: run-add
run-add: ## Run add app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go add $(ADD_ARGS)

.PHONY: run-add-ux
run-add-ux: ## Run add app in an interactive mode.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go add -u -n default

.PHONY: run-add-ux-with-ns
run-add-ux-with-ns: ## Run add app in an interactive mode with predefined namespace.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go add -u -n $(KBNAMESPACE)

.PHONY: run-import
run-import: ## Run import app .
	$(GO_TOOL) run cmd/kbcli/main.go import

.PHONY: run-import-sample
run-import-sample: ## Run import app to load sample kbs.
	$(GO_TOOL) run cmd/kbcli/main.go import -f ../../docs/samples/import-sample.yaml

.PHONY: run-get
run-get: ## Run get app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go get $(GET_ARGS)

.PHONY: run-get-ux
run-get-ux: ## Run get app with ux.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go get -u

.PHONY: run-update
run-update: ## Run update kb app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go update

.PHONY: run-add-with-args
run-add-with-args: ## Run add app with predefined arguments.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go add \
	-k btc -v crypto -o currencies -c crypto \
	-t btc,crypto,currencies,blockchain \
	-r dementor -n personal

.PHONY: run-add-media
run-add-media: ## Run add app to save kb with media.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go add \
	-k btc -o currencies -c media \
	-t btc,crypto,currencies,blockchain \
	-r dementor -n personal \
	-v 'https://pbs.twimg.com/media/GZL9kSeXgAAXf3B?format=jpg&name=4096x4096'

.PHONY: run-import-with-args
run-import-with-args: ## Run import app with predefined arguments.
	$(GO_TOOL) run cmd/kbcli/main.go import -f ../../docs/samples/import-sample.yaml --show-added-kbs --show-failed-kbs

.PHONY: run-export-with-args
run-export-with-args: ## Run export app with predefined arguments.
	@$(GO_TOOL) run cmd/kbcli/main.go export -c quote

.PHONY: run-export-all
run-export-all: ## Run export app to get all kbs.
	@$(GO_TOOL) run cmd/kbcli/main.go export

.PHONY: run-export-with-ns-cat
run-export-with-ns-cat: ## Run export app to get all kbs that match specific category and namespace.
	@$(GO_TOOL) run cmd/kbcli/main.go export -c quote -n Default

.PHONY: run-sync-with-args
run-sync-with-args: ## Run sync app with predefined arguments.
	$(GO_TOOL) run cmd/kbcli/main.go sync --show-added-kbs --show-failed-kbs
