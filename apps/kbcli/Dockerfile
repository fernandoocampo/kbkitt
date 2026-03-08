# Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Copy the go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
ARG VERSION=0.1.0
ARG COMMIT_HASH=unknown
ARG BUILD_DATE=unknown
ARG LDFLAGS="-X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions.Version=${VERSION} -X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions.CommitHash=${COMMIT_HASH} -X github.com/fernandoocampo/kbkitt/apps/kbcli/internal/cmds/versions.BuildDate=${BUILD_DATE} -s -w"

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -o /kbcli ./cmd/kbcli/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage to a common path
COPY --from=builder /kbcli /usr/local/bin/kbcli

# Define the entrypoint
ENTRYPOINT ["kbcli"]
