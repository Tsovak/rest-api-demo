export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

APPNAME = app
LDFLAGS ?=

#.DEFAULT_GOAL := all

.PHONY: all
all: build

.PHONY: mod
mod:
	go mod download

.PHONY: build
build:
	go build -o $(APPNAME) -ldflags "$(LDFLAGS)" *.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: unit
unit:
	go test -v ./... -count 10 -race

.PHONY: test
test: unit
	go test -v ./... -tags integration -count 1 -race