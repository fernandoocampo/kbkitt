# Read the file content first to understand how build-macos-arm-64 is defined.

export MSYS_NO_PATHCONV := 1
GOPATH_FWD=$(subst \,/,$(GOPATH_HOST))
GO_SHELL_TOOL=docker run --rm -v $(CURDIR):/app -v $(GOPATH_FWD):/go -e GOPATH=/go -w /app -e CGO_ENABLED -e GOOS -e GOARCH $(GO_IMAGE) sh -c
GO_SHELL_TEST_TOOL=docker run --rm -v $(CURDIR):/app -v $(GOPATH_FWD):/go -e GOPATH=/go -w /app -e CGO_ENABLED -e GOOS -e GOARCH $(GO_TEST_IMAGE) sh -c

TEST_IMAGE_SENTINEL=.test-image-built

.PHONY: test-image
test-image: ## Build the CI test image (auto-triggered by test/lint/coverage).
	docker build -f Dockerfile.ci -t $(GO_TEST_IMAGE) .
	@touch $(TEST_IMAGE_SENTINEL)

$(TEST_IMAGE_SENTINEL): Dockerfile.ci
	$(MAKE) test-image

.PHONY: test
test: $(TEST_IMAGE_SENTINEL) ## Run unit tests.
	$(GO_SHELL_TEST_TOOL) "go test -race -count=1 ./..."

.PHONY: test-ci
test-ci: ## Run unit tests in CI (assumes already in container).
	go test -race -count=1 ./...

.PHONY: coverage
coverage: $(TEST_IMAGE_SENTINEL) ## Run unit tests with coverage report.
	$(GO_SHELL_TEST_TOOL) "go test -race -count=1 -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"

.PHONY: mod-tidy
mod-tidy:
	$(GO_TOOL) mod tidy

.PHONY: lint
lint: $(TEST_IMAGE_SENTINEL) ## Run linter.
	$(GO_SHELL_TEST_TOOL) "go run github.com/golangci/golangci-lint/$(LINT_MAJOR_VERSION)/cmd/golangci-lint@$(LINT_VERSION) run --allow-parallel-runners -c $(LINT_PATH).golangci.yml"

.PHONY: lint-ci
lint-ci: ## Run linter in CI (assumes already in container).
	go run github.com/golangci/golangci-lint/$(LINT_MAJOR_VERSION)/cmd/golangci-lint@$(LINT_VERSION) run --allow-parallel-runners -c $(LINT_PATH).golangci.yml

.PHONY: build-macos-amd-64
build-macos-amd-64: ## Build binary for MacOS amd64
	@mkdir -p bin
	env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN}-amd64 ./cmd/kbcli/main.go

.PHONY: build-macos-arm-64
build-macos-arm-64: ## Build binary for MacOS arm64 (Apple Silicon)
	@mkdir -p bin
	env CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN}-arm64 ./cmd/kbcli/main.go

.PHONY: docker-build
docker-build: ## Build docker image
	docker build \
	--build-arg VERSION=${VERSION} \
	--build-arg COMMIT_HASH=${COMMIT_HASH} \
	--build-arg BUILD_DATE=${BUILD_DATE} \
	-t $(IMAGE_NAME) .

.PHONY: clean
clean: ## Remove build artifacts and test image sentinel.
	rm -f $(TEST_IMAGE_SENTINEL)
