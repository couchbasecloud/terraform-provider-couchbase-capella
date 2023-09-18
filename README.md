# Terraform Provider Capella 

This is the repository for the Couchbase Terraform Provider Capella which which forms a Terraform plugin for use with Couchbase Capella.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Using the provider

### Prepare Terraform for local provider install
Terraform installs providers and verifies their versions and checksums when you run `terraform init`. Terraform will download your 
providers from either the provider registry or a local registry. However, while building your provider you will want to 
test a Terraform configuration against a local development build of the provider. The development build will not have an associated 
version number or an official set of checksums listed in a provider registry.

Terraform allows you to use local provider builds by setting a dev_overrides block in a configuration file called .terraformrc. 
This block overrides all other configured installation methods.

Terraform searches for the .terraformrc file in your home directory and applies any configuration settings you set. 

#### For Mac


First, find the GOBIN path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

```shell
go env GOBIN
/Users/<Username>/go/bin
```

Create a new file called .terraformrc in your home directory (~), then add the dev_overrides block below. 
Change the <PATH> to the value returned from the go env GOBIN command above. 
If the GOBIN go environment variable is not set, use the default path, /Users/<Username>/go/bin.

```shell
provider_installation {

dev_overrides {
"hashicorp.com/couchabasecloud/capella" = "<PATH>"
}

# For all other providers, install them directly from their origin provider
# registries as normal. If you omit this, Terraform will _only_ use
# the dev_overrides block, and so no other providers will be available.
direct {}
}
```

Now build the terraform provider from this source code

`go build -o <PATH>`


### Configuring Programmatic Access

In order to set up authentication with the Couchbase Capella provider a V4 API key must be generated. We need base 64 encoded api-key

### Authenticating the Provider
You will need to provide host of the capella and your credentials for authentication

### Example Usage

```terraform
terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchabasecloud/capella"
    }
  }
}

provider "capella" {
  host     = "hostname of the capella instance"
  authentication_token = "capella authentication token"
}


resource "capella_project" "example" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  name = "example-name"
  description = "example-description"
}

output "example_project" {
  value = capella_project.example
}
```


You can also provide host of the capella and your credentials for authentication via the environment variables,
`CAPELLA_HOST`, `CAPELLA_AUTHENTICATION_TOKEN` for host and your authentication token.

Usage :

```shell
$  export CAPELLA_HOST="xxxx"
$  export CAPELLA_AUTHENTICATION_TOKEN="xxxx"
```

**1\. Review the Terraform plan**
Execute the following command to automatically review and update the formatting of .tf files.
```bash
$ terraform fmt
```

Execute the following command to review the resources that will be deployed.

```bash
$ terraform plan
```

**2\. Execute the Terraform apply**

Execute the plan to deploy the Couchbase Capella resources.

```bash
$ terraform apply
```

**3\. Destroy the resources**

Execute the following command to destroy all the resources.

```bash
$ terraform destroy
```

To destroy specific resource

```bash
$ terraform destroy -target=RESOURCE_ADDRESS
```
Example

```bash
$ terraform destroy -target=capella_project.example
```

**4\. To refresh the state file to sync with the remote**

```bash
$ terraform apply --refresh-only
```

**5\. To import remote resource**

```bash
$ terraform import RESOURCE_TYPE.NAME RESOURCE_IDENTIFIER
```
