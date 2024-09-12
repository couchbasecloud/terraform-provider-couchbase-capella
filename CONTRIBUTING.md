# Contributing to the provider

Thank you for your interest in contributing to the Couchbase Capella Terraform Provider. Before you start, please take a moment to read through our contribution guidelines to ensure a smooth collaboration.

## Requirements

- [Git](https://git-scm.com/)
- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.2
- [Go](https://golang.org/doc/install) >= 1.21

## Environment

- Fork the repository.
- Clone your fork. Use `git clone` to create a local copy on your machine.
- We use Go Modules to manage dependencies, so you can develop outside your `$GOPATH`.
- We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make setup`.

## Local Development Setup

Quickly set up and run Terraform Provider Couchbase Capella locally with these steps:

- Enter the provider directory.
- Install necessary tools with `make setup`.
- Build the local provider binary using `make build`.
- Set a development override to use the local provider binary.
- Define environment variables for the Terraform configuration file.
- Initialize, plan, and apply your Terraform configuration.

#### Enter the Provider Directory:

Navigate to the directory containing the provider's source code. For example: 

```bash
cd $HOME/terraform-provider-couchbase-capella
```

#### Install Dependencies:

Run the following command to install the required tools for building the provider:

```bash
make setup
```

#### Build the Provider Binary:

Generate the binary by executing:

```bash
make build
```

This will create the binary in the `./bin` directory.

#### Use the Provider Binary:

Normally, Terraform installs providers and verifies their versions and checksums when you run `terraform init`.
Terraform will download your providers from either the provider registry or a local registry. However, during
development, we want to test using a a local development build of the provider. The development build will not
have an associated version number or an official set of checksums listed in a provider registry.

After successfully building the local provider binary for Couchbase Capella, the  next step is to instruct
Terraform to use local provider builds by setting a dev_overrides block in a configuration file
with ext .terraformrc or .tfrc. This block overrides all other configured installation methods.
When Terraform runs, it searches for any .terraformrc or .tfrc file in your home directory and applies any
configuration settings that have been set.

#### Create the configuration file if it doesn't already exist

```bash
touch $HOME/dev.terraformrc
```

#### Open and edit the configuration file

**LINUX SPECIFIC STEPS:**

Navigate to `dev.terraformrc` either using console or your preferred text editor and paste in
the following contents:

```terraform
  provider_installation {

  dev_overrides {
    "couchbasecloud/couchbase-capella" = "$HOME/terraform-provider-couchbase-capella/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
  ```
**WINDOWS SPECIFIC STEPS:**

Step 1: Find your %APPDATA% path from powershell
(reference: https://developer.hashicorp.com/terraform/cli/config/config-file#locations)

Step 2: Create a terraform.rc file in the path above

```terraform
  provider_installation {

  dev_overrides {
    "couchbasecloud/couchbase-capella" = "C:\\Users\\Administrator\\couchbasecapella"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}

  ```
**NOTE: Please make sure you are escaping the backslashes in the path.**

#### Define the env var `TF_CLI_CONFIG_FILE` in your console session

  ```bash
  export TF_CLI_CONFIG_FILE=$HOME/dev.terraformrc
  ```

#### Change to an example directory

This repository contains a number of example terraform configurations to enable you to get started quickly.
To get started with an example configuration provided by Couchbase, execute the following command:

```bash
cd $HOME/terraform-provider-couchbase-capella/examples/getting_started
```

If you are not using an example and instead using your own custom configuration, you will instead need
to change into that directory. For a full list of examples, see the appendinx section on examples.

#### Terraform Environment Variables

Each example configuration contains a terraform.templates.tfvars file. This file contains the values
environment variables defined in the example configurations. To use these values, first copy the values
to a new file:

```bash
cp terraform.template.tfvars terraform.tfvars
```

Then, open the file up and replace the placeholders for auth_token and organization_id with their
actual values.

```terraform
auth_token = "<v4-api-key-secret>"
organization_id = "<organization-uuid>"
```

#### Initializing and Running Terraform

- Run `terraform init` to initialize terraform
- Run `terraform plan` to use terraform with the local binary and preview the infrastructure that will be created
- Run `terraform apply` to execute the plan generated by terraform plan. This can also be ran by itself.

For more information about development overrides see [Development Overrides for Provider Developers](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers)

#### Discovering New API features

Most of the new features of the provider are using [capella-public-apis](https://docs.couchbase.com/cloud/management-api-guide/management-api-intro.html)
Public APIs are updated automatically, tracking all new Capella features.

## Appendix A - Examples

Note: You will need to provide both the url of the capella host as well as your V4 API secret for authentication.

```terraform
terraform {
  required_providers {
    couchbase-capella = {
      source = "couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
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

#### Terraform Init

Ordinarily, terraform will download the requested providers on running the command:
```bash
$ terraform init
```
If you are working with a local install of `Terraform-Provider-Couchbase-Capella` provider, this step is not needed and considered optional.
However, if you plan to use any other providers at the same time it may need to be run.

**1\. Review the Terraform plan**

Execute the following command to automatically review and update the formatting of .tf files.
```bash
$ terraform fmt
```

The terraform.template.tfvars defines a template on how the variable values should be added.
Copy this file as `terraform.tfvars` in the same directory and add the value for each variable.

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

### Running the acceptance test

~> **Notice:** Acceptance tests create real resources, and often cost money to run. Please note in any PRs made if you are unable to pay to run acceptance tests for your contribution. We will accept "best effort" implementations of acceptance tests in this case and run them for you on our side. This may delay the contribution but we do not want your contribution blocked by funding.
- Run `make testacc`

## Appendix B - Creating your own environment variables

Environment variables can be set by terraform by creating and adding a terraform.tfvars to your
configuration.

```terraform
auth_token = "<v4-api-key-secret>"
organization_id = "<organization-uuid>"
```

A variables.tf should also be added to define the variables for terraform.
```terraform
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

Alternatively, if you would like to set environment variables locally on your system (as opposed to using terraform.tfvars),
preface them with `TF_VAR_`. Terraform will then apply them your .terraformrc file on running
`terraform apply`. For example:
```bash
export TF_VAR_auth_token=<v4_api_secret_key>
export TF_VAR_organization_id=<organization_id>
```

## Appendix C - Authentication
```
In order to set up authentication with the Couchbase Capella provider a V4 API key must be generated.

To find out how to generate a V4 API Key, please see the following document:
https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html

Once you have generated your api key token, it must be set as an environment variable. 
```
