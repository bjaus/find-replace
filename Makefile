.PHONY: help
## help| Show this help dialogue
help:
	@sed -n 's/^##//p' Makefile | column -t -c 2 -s '|'

.PHONY: build
## build| Create binary
build:
	@go build -o /usr/local/bin/replace .

.PHONY: lint
## lint| Run linters
lint: mod
	@golangci-lint run

## tidy| Run go mod tidy
.PHONY: tidy
tidy:
	@go mod tidy

## vendor| Run go mod vendor
.PHONY: vendor
vendor:
	@go mod vendor

## mod| Run go mod tidy & vendor
.PHONY: mod
mod: tidy vendor
