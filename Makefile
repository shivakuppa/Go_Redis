SHELL := /bin/bash
GO_BIN ?= $(shell pwd)/.bin
GOCI_LINT_VERSION ?= v2.3.1

PATH := $(GO_BIN):$(PATH)

format-go::
	golangci-lint run --fix ./...

tools::
	mkdir -p $(GO_BIN)
	curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b ${GO_BIN} ${GOCI_LINT_VERSION}

run::
	go run ./cmd/server/main.go
