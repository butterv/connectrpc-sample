export

GOBIN := $(PWD)/bin
PATH := $(GOBIN):$(PATH)

.PHONY: init
init:
	@go mod download

.PHONY: dependencies
dependencies:
	@./scripts/dependencies.sh

.PHONY: lint
lint:
	@$(GOBIN)/staticcheck ./...
	@$(GOBIN)/go-consistent ./...

.PHONY: test
test:
	@go test -v -race ./...

.PHONY: test-coverage
test-coverage:
	@go test -race -cover -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
