.DEFAULT_GOAL := help

.PHONY: help build check ci clean codex-links fmt test

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-14s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build Go and TypeScript packages
	go build ./...
	pnpm build

check: ## Run static checks and schema validation
	test -z "$$(gofmt -l ./cmd ./internal)"
	go vet ./...
	golangci-lint run ./...
	GOTOOLCHAIN=auto go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	pnpm check
	pnpm lint
	pnpm verify:schemas

ci: check test build ## Run the complete local CI suite

codex-links: ## Create or refresh local AI Central links
	pnpm codex:links

fmt: ## Format Go and TypeScript sources
	gofmt -w cmd internal
	pnpm format

test: ## Run all tests
	go test ./...
	pnpm test

clean: ## Remove generated build output
	pnpm clean
	go clean
