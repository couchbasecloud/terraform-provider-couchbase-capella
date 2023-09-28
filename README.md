# Terraform Provider Capella 

This is the repository for Couchbase's Terraform-Provider-Capella which forms a Terraform plugin for use with Couchbase Capella.

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

#### Create the terraform configuration file
Create a new file called .terraformrc in your home directory (~), then add the dev_overrides block below.

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

`<PATH>` should be replaced with the directory corresponding to the terraform provider binary. 
This could be the default folder where Go installs your binaries, which can be determined by calling: 

```shell
go env GOBIN
```

Alternatively, it may be the default path which is of the form: 
``` shell
/Users/<Username>/go/bin
```

#### Build and install the executable 

Now build the terraform provider. Ensure to specify the build location using `<PATH>` as described above. 

`go build -o <PATH>`

### Authentication

In order to set up authentication with the Couchbase Capella provider a V4 API key must be generated. 

To find out how to generate a V4 API Key, please see the following document: 
https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html

### Terraform Environment Variables

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

Set the environment variables by using the following notation: 
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
export TF_VAR_host= "https://cloudapi.cloud.couchbase.com"
```

### Create and manage resources using terraform

#### Example Usage

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

This repository contains a number of example directories containing examples of Hashicorp Configuration Language (HCL) code 
being used to create and manage Capella resources. To try these examples out for yourself, change into one of them and run
the below commands.

#### Commands

#### n.b. Terraform Init

Ordinarily, terraform will downloaded the requested providers on running the command: 
```bash
$ terraform init
```
As we are working with a local install of `Terraform-Provider-Capella` this command is not needed and considered optional. 
However if you plan to use any other providers at the same time it may need to be ran. 

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
