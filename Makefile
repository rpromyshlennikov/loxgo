GIT_VERSION ?= $(shell git describe --abbrev=4 --dirty --always --tags)

APP_NAME = loxgo
BINARY ?= ./${APP_NAME}

.PHONY: tools
tools: ## Install all needed tools, e.g. for static checks
	@echo Installing tools from tools-versions.txt
	@xargs -tI % go install % < tools-versions.txt

# Main targets

all: test build
.DEFAULT_GOAL := all


.PHONY: build
build: ## Build the project binary
	go build -ldflags "-X main.version=$(GIT_VERSION)" ./cmd/$(APP_NAME)/

.PHONY: test
test: ## Run unit (short) tests
	go test -short ./...

.PHONY: bench
bench: ## Run benchmarks
	go test ./... -short -bench=. -run="Benchmark*"

.PHONY: lint
lint: tools ## Check the project with lint
	staticcheck ./...

.PHONY: vet
vet: ## Check the project with vet
	go vet ./...

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do gofmt -e -l -w $$d/*.go; done)

.PHONY: imports
imports: tools ## Check and fix import section by import rules
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do goimports -e -l -local $$(go list) -w $$d/*.go; done)

.PHONY: cyclomatic
cyclomatic: tools ## Check the project with gocyclo for cyclomatic complexity
	gocyclo -over 10 `find . -type f -iname '*.go' -not -iname '*_test.go' -not -path '*/\.*'`

.PHONY: static_check
static_check: fmt imports vet lint cyclomatic ## Run static checks (fmt, lint, imports, vet, ...) all over the project

.PHONY: check
check: static_check test ## Check project with static checks and tests

.PHONY: dependencies
dependencies: ## Manage go mod dependencies, beautify go.mod and go.sum files
	go mod tidy

.PHONY: run
run: build ## Start the project
	$(BINARY)

.PHONY: astgen
astgen: ## Build the AST generator
	go build -ldflags "-X main.version=$(GIT_VERSION)" ./cmd/astgenerator/


.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean:
	@rm -f ${BINARY}

