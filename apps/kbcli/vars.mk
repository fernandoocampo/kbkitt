GOPATH_HOST?=$(GOPATH)
GO_IMAGE?=golang:1.25.7
GO_TOOL=docker run --rm -v $(CURDIR):/app -v $(GOPATH_HOST):/go -e GOPATH=/go -w /app -e CGO_ENABLED -e GOOS -e GOARCH $(GO_IMAGE) go
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

# Docker variables
IMAGE_NAME?=kbcli-image
DOCKER_RUN=docker run --rm -it -v $(CURDIR):/app -w /app $(IMAGE_NAME)
