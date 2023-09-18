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


### Authentication

In order to set up authentication with the Couchbase Capella provider a V4 API key must be generated. 

To find out how to generate a V4 API Key, please see the following document: 
https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html

### Example Usage

Note: You will need to provide both the url of the capella host as well as your V4 API secret for authentication. 

```terraform
terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchabasecloud/capella"
    }
  }
}

provider "capella" {
  host     = "the host url of couchbase cloud"
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

# Terraform Environment Variables

Environment variables can be set by terraform by creating and adding terraform.template.tfvars
```terraform
auth_token = "<v4-api-key-secret>"
organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
host = "https://cloudapi.dev.nonprod-project-avengers.com"
```

A variables.tf should also be added to define the variables for terraform. 
```terraform
variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
}
```

The environment variables by uisng the following notation: 
```terraform
resource "capella_project" "example" {
  organization_id = var.organization_id
  name = var.project_name
  description = "A Capella Project that will host many Capella clusters."
}
```

Alternatively, if you would like to set environment variables locally on your system (as opposed to using terraform.template.tfvars),
preface them with `TF_VAR_`. Terraform will then apply them your .terraformrc file on running
`terraform apply`. For example: 
```bash
export TF_VAR_auth_token=<v4_api_secret_key>
export TF_VAR_organization_id="6af08c0a-8cab-4c1c-b257-b521575c16d0"
export TF_VAR_host= "https://cloudapi.dev.nonprod-project-avengers.com"
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
