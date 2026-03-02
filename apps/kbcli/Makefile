include vars.mk
include build.mk
include run.mk

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help.
	@echo "Available targets:"
	@echo "  docker-build        Build docker image"
	@echo "  build-macos-amd-64  Build binary for MacOS amd64"
	@echo "  test                Run unit tests"
	@echo "  run                 Run run app"
	@echo "  run-add             Run add app"
	@echo "  run-get             Run get app"
	@echo "  run-import          Run import app"
	@echo "  run-export          Run export app"
	@echo "  run-sync            Run sync app"
	@echo "  (and more... see run.mk and build.mk)"
