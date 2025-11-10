GOFMT_FILES?=$$(find . -name '*.go')
BINARY_NAME=terraform-provider-couchbase-capella
DESTINATION=./bin/$(BINARY_NAME)

GOFLAGS=-mod=vendor
GOOPTS="-p 2"

GITTAG=$(shell git describe --tags --abbrev=0)
VERSION=$(GITTAG:v%=%)
LINKER_FLAGS=-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=${VERSION}'

GOLANGCI_VERSION=v1.64.8

export PATH := $(shell go env GOPATH)/bin:$(PATH)
export SHELL := env PATH=$(PATH) /bin/bash

default: build

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
endif

.PHONY: help
help: ## Show this help message
	@echo "Terraform Provider Couchbase Capella - Available Commands:"
	@echo ""
	@grep -h -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick Start:"
	@echo "  make setup    - Install all dev tools"
	@echo "  make check    - Run all quality checks"
	@echo "  make test     - Run unit tests"
	@echo "  make testacc  - Run acceptance tests"

# ============================================================================
# Build
# ============================================================================

.PHONY: build
build: ## Build the provider binary
	@$(MAKE) fmt
	@go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: setup
setup: ## Install all dev tools
	@echo "==> Installing dependencies..."
	@go install github.com/client9/misspell/cmd/misspell@latest
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: check
check: ## Run all quality checks (fmt, lint, test)
	@$(MAKE) setup tffmt fmt docs-lint lint-fix test

.PHONY: fmt
fmt: ## Format Go code
	@echo "==> Fixing source code with gofmt and goimports..."
	@gofmt -s -w $(GOFMT_FILES)
	@find . -name "*.go" -exec goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella {} \;

.PHONY: tffmt
tffmt: ## Format Terraform files
	@terraform fmt -write -recursive -diff .

.PHONY: lint-fix
lint-fix: ## Fix linter errors automatically
	@echo "==> Fixing linter errors..."
	@golangci-lint run --fix

.PHONY: vet
vet: ## Run go vet
	@echo "==> Running go vet..."
	@go vet ./... ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: tfcheck
tfcheck: ## Check Terraform formatting
	@terraform fmt -check=true -diff -recursive .

.PHONY: docs-lint
docs-lint: ## Lint documentation
	@echo "==> Checking docs against linters..."
	@misspell -error -source=text docs/

# ============================================================================
# Testing
# ============================================================================

TEST_FILES ?= $$(go list ./... | grep -v acceptance_tests)
TEST_FLAGS ?= -short -cover -race -coverprofile .testCoverage.txt

.PHONY: test
test: ## Run unit tests
	@CAPELLA_OPENAPI_SPEC_PATH=$(PWD)/openapi.generated.yaml go test $(TEST_FILES) $(TEST_FLAGS)

.PHONY: testacc
testacc: ## Run acceptance tests (requires TF_VAR_auth_token, TF_VAR_host, TF_VAR_organization_id)
	@[ "${TF_VAR_auth_token}" ] || ( echo "ERROR: export TF_VAR_auth_token before running acceptance tests"; exit 1 )
	@[ "${TF_VAR_host}" ] || ( echo "ERROR: export TF_VAR_host before running acceptance tests"; exit 1 )
	@[ "${TF_VAR_organization_id}" ] || ( echo "ERROR: export TF_VAR_organization_id before running acceptance tests"; exit 1 )
	@CAPELLA_OPENAPI_SPEC_PATH=$(PWD)/openapi.generated.yaml TF_ACC=1 go test -timeout=120m -v ./acceptance_tests/

# ============================================================================
# Documentation
# ============================================================================

.PHONY: build-docs
build-docs: ## Generate provider documentation
	@go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	@CAPELLA_OPENAPI_SPEC_PATH="$(shell pwd)/openapi.generated.yaml" go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --examples-dir ./examples

# ============================================================================
# Code Generation
# ============================================================================

.PHONY: gen-api
gen-api: ## Generate OpenAPI client code
	@echo "==> Generating OpenAPI client (internal/generated/api)"
	@PATH="$(shell go env GOPATH)/bin:$(PATH)" go generate ./internal/generated/api
	@echo "==> Done"

# ============================================================================
# Release Management
# ============================================================================

.PHONY: release-prep
release-prep: ## Prepare release artifacts (usage: make release-prep VERSION=1.5.4)
	@echo "==> Preparing release..."
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION not specified"; \
		echo "Usage: make release-prep VERSION=1.5.4"; \
		exit 1; \
	fi
	@echo "Detecting previous version..."
	$(eval PREV_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v1.5.3"))
	@echo "Previous version: $(PREV_VERSION)"
	@echo "New version: v$(VERSION)"
	@echo ""
	@echo "[1/4] Generating changelog..."
	@bash scripts/generate-changelog.sh v$(VERSION) $(PREV_VERSION)
	@echo ""
	@echo "[2/4] Validating PR quality..."
	@if command -v python3 &> /dev/null && python3 -c "import github" 2>&1 >/dev/null; then \
		python3 scripts/validate-prs.py v$(VERSION) $(PREV_VERSION) || echo "   WARNING: Validation warnings found (continuing anyway)"; \
	else \
		echo "   Skipped (install PyGithub for enhanced validation)"; \
	fi
	@echo ""
	@echo "[3/4] Generating upgrade guide..."
	@if command -v python3 &> /dev/null && python3 -c "import github" 2>&1 >/dev/null; then \
		python3 scripts/generate-upgrade-guide.py $(VERSION) $(PREV_VERSION); \
	else \
		bash scripts/generate-upgrade-guide.sh $(VERSION) $(PREV_VERSION); \
		echo "   Tip: Install PyGithub for enhanced content extraction"; \
	fi
	@echo ""
	@echo "[4/4] Building documentation..."
	@$(MAKE) build-docs
	@echo ""
	@echo "SUCCESS: Release preparation complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Review templates/guides/$(VERSION)-upgrade-guide.md"
	@echo "  2. Review CHANGELOG.md"
	@echo "  3. Commit: git add . && git commit -m 'Prepare release v$(VERSION)'"
	@echo "  4. Tag: git tag v$(VERSION) && git push origin v$(VERSION)"
	@echo "  5. GitHub Actions will create the release automatically"
	@echo ""
	@echo "For full checklist: make release-checklist"

.PHONY: release-checklist
release-checklist: ## Show complete release checklist
	@echo "Terraform Provider Release Checklist"
	@echo ""
	@echo "Pre-Release:"
	@echo "  [ ] All PRs merged and labeled correctly"
	@echo "  [ ] make test - Unit tests pass"
	@echo "  [ ] make testacc - Acceptance tests pass"
	@echo "  [ ] make check - Code quality checks pass"
	@echo ""
	@echo "Generate Artifacts:"
	@echo "  [ ] make release-prep VERSION=X.Y.Z"
	@echo "  [ ] Review/edit templates/guides/X.Y.Z-upgrade-guide.md"
	@echo "  [ ] Review CHANGELOG.md"
	@echo "  [ ] Test examples from upgrade guide"
	@echo ""
	@echo "Release:"
	@echo "  [ ] Commit all changes"
	@echo "  [ ] git tag vX.Y.Z && git push origin vX.Y.Z"
	@echo "  [ ] Verify GitHub Actions workflow succeeds"
	@echo "  [ ] Verify release on GitHub releases page"
	@echo "  [ ] Verify published to Terraform Registry"
	@echo ""
	@echo "Post-Release:"
	@echo "  [ ] Announce in team channels"
	@echo "  [ ] Update related documentation"
	@echo ""
	@echo "For details: See RELEASE_PROCESS.md"


