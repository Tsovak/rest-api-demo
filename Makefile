export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE=on

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

