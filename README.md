# Terraform Provider Capella 

This is the repository for Couchbase's Terraform-Provider-Capella which forms a Terraform plugin for use with Couchbase Capella.

## Requirements

- [Git](https://git-scm.com/)
- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.2
- [Go](https://golang.org/doc/install) >= 1.21

### Environment

- We use Go Modules to manage dependencies, so you can develop outside your `$GOPATH`.
- We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make setup`.

## Using the Provider

To use a released provider in your Terraform environment, run `terraform init` and Terraform will automatically install the provider.
Documentation about the provider specific configuration options can be found on the [provider's website](https://developer.hashicorp.com/terraform/language/providers).

## Contributing to the Provider
See [Contributing.md](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/CONTRIBUTING.md)

## Discovering New API features

Most of the new features of the provider are using [capella-public-apis](https://docs.couchbase.com/cloud/management-api-guide/management-api-intro.html)
Public APIs are updated automatically, tracking all new Capella features.

## Generated API client

This repository includes an OpenAPI-generated client in `internal/generated/api` (keeping the existing hand-written client in `internal/api` for backward compatibility).

Generate/update the client before working on a new resource or data source:

1) Ensure the generator is installed:

   `go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`

2) Regenerate from the root `openapi.generated.yaml`:

   `make gen-api`

The command writes the client/types to `internal/generated/api/openapi.gen.go`.

Notes:
- Provider wiring makes both clients available:
  - `providerschema.Data.ClientV1`: legacy HTTP client (`internal/api`)
  - `providerschema.Data.ClientV2`: generated client with typed methods (`internal/generated/api`)
- When adding a new resource/data source, prefer calling `ClientV2` for new endpoints and migrate incrementally.
