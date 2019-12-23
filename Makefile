export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

BIN_DIR = bin
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
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APPNAME) cmd/app/*.go
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/migrate    cmd/migration/*.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: unit
unit:
	go test -v ./... -count 10 -race

.PHONY: test
test: unit
	go test -v ./... -tags integration -count 10 -race --failfast

.PHONY: migrate
migrate: ## migration
	go run ./cmd/migration/main.go -dir scripts/migrations -init

.PHONY: lint
lint: ## linter
	@golangci-lint --color=always run ./... -v --timeout 5m
