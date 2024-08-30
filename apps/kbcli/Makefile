GO_TOOL?=go1.23.0
ADD_ARGS?=
GET_ARGS?=
ROOT_ARGS?=
COMMIT_HASH?=$(shell git describe --dirty --tags --always)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION?=0.1.0
LDFLAGS?="-X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds.Version=${VERSION} -X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds.CommitHash=${COMMIT_HASH} -X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds.BuildDate=${BUILD_DATE} -s -w"

.PHONY: test
test: ## Run unit tests.
	$(GO_TOOL) test -race -count=1 ./...

.PHONY: run
run: ## Run run app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go $(ROOT_ARGS)

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

.PHONY: run-get
run-get: ## Run get app.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go get $(GET_ARGS)

.PHONY: run-add-with-args
run-add-with-args: ## Run add app with predefined arguments.
	$(GO_TOOL) run -ldflags ${LDFLAGS} cmd/kbcli/main.go add -k btc -v crypto -n currencies -c crypto -t btc,crypto,currencies,blockchain
                                                