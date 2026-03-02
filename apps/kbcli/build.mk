.PHONY: test
test: ## Run unit tests.
	$(GO_TOOL) test -race -count=1 ./...

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
