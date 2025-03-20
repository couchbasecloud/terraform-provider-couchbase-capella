GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

BINARY_NAME=terraform-provider-couchbase-capella
DESTINATION=./bin/$(BINARY_NAME)

GOFLAGS=-mod=vendor
GOOPTS="-p 2"
GOFMT_FILES?=$$(find . -name '*.go')

GITTAG=$(shell git describe --tags --abbrev=0)
VERSION=$(GITTAG:v%=%)
LINKER_FLAGS=-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=${VERSION}'

GOLANGCI_VERSION=v1.55.2

export PATH := $(shell go env GOPATH)/bin:$(PATH)
export SHELL := env PATH=$(PATH) /bin/bash

default: build

# Tries to configure GOPATH if it isn't set.
ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
endif

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: fmt
	go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

.PHONY: fmt
fmt:
	@echo "==> Fixing source code with gofmt and goimports..."
	gofmt -s -w $(GOFMT_FILES)
	find . -name "*.go" -exec goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella {} \;

.PHONY: vet
vet:
	@echo "==> Running go vet ."
	@go vet ./... ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: lint-fix
lint-fix:
	@echo "==> Fixing linters errors..."
	fieldalignment -json -fix ./...
	golangci-lint run --fix

.PHONY: setup
setup:  ## Install dev tools
	@echo "==> Installing dependencies..."
	go install github.com/client9/misspell/cmd/misspell@latest
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: check
check: setup tffmt tfcheck fmt docs-lint lint-fix test

.PHONY: docs-lint
docs-lint:
	@echo "==> Checking docs against linters..."
	@misspell -error -source=text docs/

.PHONY: docs
docs:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

.PHONT: build-docs
build-docs:
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: terraform-check tfcheck
tfcheck: terraform-check
terraform-check:
	@terraform fmt -check=true -diff -recursive .

.PHONY: terraform-fmt tffmt
tffmt: terraform-fmt
terraform-fmt:
	@terraform fmt -write -recursive -diff .
.PHONY: terraform-check tfcheck

TEST_FILES ?= $$(go list ./... | grep -v acceptance_tests)
TEST_FLAGS ?= -short -cover -race -coverprofile .testCoverage.txt

# this is for unit tests
.PHONY: test
test:
	go test $(TEST_FILES) $(TEST_FLAGS)

.PHONY: test-acceptance testacc
testacc: test-acceptance
test-acceptance:
	@[ "${TF_VAR_auth_token}" ] || ( echo "export TF_VAR_auth_token before running the acceptance tests"; exit 1 )
	@[ "${TF_VAR_host}" ] || ( echo "export TF_VAR_host before running the acceptance tests"; exit 1 )
	@[ "${TF_VAR_organization_id}" ] || ( echo "export TF_VAR_organization_id before running the acceptance tests"; exit 1 )
	TF_ACC=1 go test -timeout=60m -v ./acceptance_tests/


