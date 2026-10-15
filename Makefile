.DEFAULT_GOAL := help

.PHONY: help build check ci clean codex-links fmt lint-go-full test

# Override with the exact PR base commit to reproduce its changed-code gates.
CI_BASE ?= origin/main

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-14s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build Go and TypeScript packages
	go build ./...
	pnpm build

check: ## Run local CI static checks (CI_BASE defaults to origin/main)
	test -z "$$(gofmt -l ./cmd ./internal)"
	go vet ./...
	node scripts/lint-go.mjs --base "$(CI_BASE)"
	GOTOOLCHAIN=auto go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
	pnpm audit:dependencies
	pnpm site:build
	pnpm verify:schemas
	pnpm verify:adrs
	pnpm check
	pnpm lint
	pnpm verify:jsdoc --changed-since "$(CI_BASE)" --include-worktree

# One prerequisite finishes before this recipe, including under make -j.
# check emits TypeScript once; test:built reuses those outputs.
ci: check ## Run local CI gates and full tests; hosted-only jobs remain separate
	go test ./...
	go build ./...
	pnpm test:built

lint-go-full: ## Audit all Go lint including existing documentation debt (2.12.2)
	node scripts/lint-go.mjs --full

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
