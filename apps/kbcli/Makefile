include vars.mk
include build.mk
include run.mk

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help.
	@echo "Available targets:"
	@echo "  docker-build        Build docker image"
	@echo "  build-macos-amd-64  Build binary for MacOS amd64"
	@echo "  build-macos-arm-64  Build binary for MacOS arm64 (Apple Silicon)"
	@echo "  test                Run unit tests"
	@echo "  install-hooks       Install git hooks (pre-commit, pre-push, commit-msg)"
	@echo "  run                 Run run app"
	@echo "  run-add             Run add app"
	@echo "  run-get             Run get app"
	@echo "  run-import          Run import app"
	@echo "  run-export          Run export app"
	@echo "  run-sync            Run sync app"
	@echo "  (and more... see run.mk and build.mk)"

.PHONY: install-hooks
install-hooks: ## Install git hooks to enforce unit tests on commit and push, and commitlint on commit-msg.
	git config core.hooksPath "$(CURDIR)/.githooks"
ifneq ($(OS),Windows_NT)
	sed -i.bak 's/\r//' .githooks/pre-commit .githooks/pre-push .githooks/commit-msg
	rm -f .githooks/pre-commit.bak .githooks/pre-push.bak .githooks/commit-msg.bak
	chmod +x .githooks/pre-commit .githooks/pre-push .githooks/commit-msg
endif
