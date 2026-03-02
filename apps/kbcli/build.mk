export MSYS_NO_PATHCONV := 1
GOPATH_FWD=$(subst \,/,$(GOPATH_HOST))
GO_SHELL_TOOL=docker run --rm -v $(CURDIR):/app -v $(GOPATH_FWD):/go -e GOPATH=/go -w /app -e CGO_ENABLED -e GOOS -e GOARCH $(GO_IMAGE) sh -c

.PHONY: test
test: ## Run unit tests.
	$(GO_SHELL_TOOL) "apt-get update -qq && apt-get install -y -qq libx11-dev && go test -race -count=1 ./..."

.PHONY: coverage
coverage: ## Run unit tests with coverage report.
	$(GO_SHELL_TOOL) "apt-get update -qq && apt-get install -y -qq libx11-dev && go test -race -count=1 -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"

.PHONY: mod-tidy
mod-tidy:
	$(GO_TOOL) mod tidy

.PHONY: build-macos-amd-64
build-macos-amd-64: ## Build binary for MacOS amd64
	@mkdir -p bin
	env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 ${GO_BUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN}-amd64 ./cmd/kbcli/main.go

.PHONY: docker-build
docker-build: ## Build docker image
	docker build \
	--build-arg VERSION=${VERSION} \
	--build-arg COMMIT_HASH=${COMMIT_HASH} \
	--build-arg BUILD_DATE=${BUILD_DATE} \
	-t $(IMAGE_NAME) .
