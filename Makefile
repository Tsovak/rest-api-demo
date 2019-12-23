export GOPATH ?= $(shell go env GOPATH)
export GO111MODULE ?= on

BIN_DIR = bin
APPNAME = app
LDFLAGS ?=
COVERPROFILE ?= coverage.txt

#.DEFAULT_GOAL := all

.PHONY: all
all: build

.PHONY: mod
mod:
	go mod download

.PHONY: clean
clean: ## run all cleanup tasks
	go clean ./...
	rm -f $(COVERPROFILE)
	rm -rf $(BIN_DIR)


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


.PHONY: test-with-coverage
test-with-coverage:
	go test -v ./... -tags integration -count 1 --coverprofile=$(COVERPROFILE) --covermode=count

.PHONY: migrate
migrate: ## migration
	go run ./cmd/migration/main.go -dir scripts/migrations -init

.PHONY: lint
lint: ## linter
	@golangci-lint --color=always run ./... -v --timeout 5m
