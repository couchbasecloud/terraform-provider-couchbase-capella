# The Couchbase Capella Terraform Provider

The Capella Terraform Provider is a powerful way of programmatically managing Capella API keys, users, organizations, projects, clusters, buckets, and other resources.


## Example Usage

```terraform
# Configure the Couchbase Capella Provider 
provider "couchbase-capella" {
  authentication_token = var.couchbasecapella_auth_token
}
# Create the resources
````

## Configure Capella Access

For authentication with the Couchbase Capella provider a [V4 REST API Key](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#understand-management-api-keys) must be [generated](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#generate-management-api-keys).
An [Organization Owner](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html) role should be sufficient to get you going --- although you may wish to review this level of access, and use a different key in production.
This API key is then used for authenticating the Terraform Provider against Couchbase Capella.


[!IMPORTANT]  
Although in examples below, the API key secret is specified as a environmental variable and hardcoded in a config file, it is recommend that the API secret credentials be stored in a remote secrets manager such as Hashicorp Vault or AWS Secrets Manager that the Terraform Provider can then retrieve and use for authentication.


### Terraform Environment Variables

Environment variables can be set by terraform by creating and adding terraform.template.tfvars
```terraform
auth_token = "<v4-api-key-secret>"
organization_id = "<replace with organization id>"
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
resource "couchbase-capella_project" "example" {
  organization_id = var.organization_id
  name = var.project_name
  description = "A Capella Project that will host many Capella clusters."
}
```

## Create and manage resources using terraform

### Example Usage

Note: You will need to provide the V4 API secret for authentication.

```terraform
terraform {
  required_providers {
    couchbase-capella = {
      source = "hashicorp.com/couchabasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  authentication_token = "capella authentication token"
}


resource "couchbase-capella_project" "example" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  name = "example-name"
  description = "example-description"
}

output "example_project" {
  value = couchbase-capella_project.example
}
```

This repository contains a number of example directories containing examples of Hashicorp Configuration Language (HCL) code being used to create and manage Capella resources.
To try these examples out for yourself, change into one of them and run the below commands.

### Commands

#### Optional Terraform Init

Ordinarily, terraform will download the requested providers on running the command:
```bash
$ terraform init
```

**1\. Review the Terraform plan**

Execute the following command to automatically review and update the formatting of .tf files.
```bash
$ terraform fmt
```

Execute the following command to review the resources that will be deployed.

```bash
$ terraform plan -var-file=terraform.template.tfvars
```
NOTE: If using a terraform.template.tfvars file to specify variables, then the -var-file flag will need to be used.
If instead, variables are set either using a terraform.tfvars file or by using TF_VAR_ prefaced environment variables, then the -var-file flag can be omitted.
This also applies for `terraform apply`.

**2\. Execute the Terraform apply**

Execute the plan to deploy the Couchbase Capella resources.

```bash
$ terraform apply -var-file=terraform.template.tfvars
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
$ terraform destroy -target=couchbase-capella_project.example
```

**4\. To refresh the state file to sync with the remote**

```bash
$ terraform apply --refresh-only
```

**5\. To import remote resource**

```bash
$ terraform import RESOURCE_TYPE.NAME RESOURCE_IDENTIFIER
```

## Version Compatibility

For Couchbase Capella provider version 1.0.0 and above:

* Terraform Version Requirement: Terraform >= 1.5.2
* Go >= 1.20

## Supported OS and Architecture

As per HashiCorp's recommendations, we fully support the following operating system / architecture combinations:

* Darwin / AMD64
* Darwin / ARMv8
* Linux / AMD64
* Linux / ARMv8 (sometimes referred to as AArch64 or ARM64)
* Linux / ARMv6
* Windows / AMD64

## More Information

* [Report bugs](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues)
* [Open an issue on Git Hub](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/issues)
* [Support Plan information](https://docs.couchbase.com/cloud/support/support.html)
* [Github repo](https://github.com/couchbasecloud/terraform-provider-couchbase-capella)

## Example Usage

To get started, see the [Provider Example Configs](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples):

* [Retrieve organization details in Capella](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/organization):

  Couchbase Capella uses an ordered hierarchy to help you keep all of your data organized and securely accessible.
  The entity at the top of the hierarchy is called an organization.
  Everything you do in Capella  --  whether it’s creating a cluster or managing billing  --  happens within the scope of an [organization](https://docs.couchbase.com/cloud/organizations/organizations.html).

* [Create and manage users](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/user):

  Users have roles within an [organization](https://docs.couchbase.com/cloud/organizations/manage-organization-users.html), and within [individual projects](https://docs.couchbase.com/cloud/projects/manage-project-users.html).

* [Create and manage API Keys](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/apikey):

  Every API key is associated with an allowed IP Address list, and one or more organization roles, which determine the [privileges that the API key has](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#understand-management-api-keys) within the organization.

* [Create & manage projects](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/project):

  Within organizations, [projects](https://docs.couchbase.com/cloud/projects/projects.html) are used to organize and manage groups of Couchbase databases.
  An organization can contain any number of projects, and a project can contain any number of databases.

* [Create & manage Capella clusters (databases)](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/cluster):

  The Cluster is the indivdual instance of a [Couchbase Database](https://docs.couchbase.com/cloud/clusters/databases.html), spanning one or more nodes on your Cloud Service Provider, and containing the Data Service, and any other services which you choose to deploy.
  Within this sits the heirarchy of bucket, scope, collection, and document.

* [Retrieve cluster certificate details](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/certificate):

  Retrieve the certificate details for a Capella cluster;
  list the certificate details based on the cluster ID and authentication access token.

* [Manage database credential](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/database_credential):

  Database credentials are separate from organization roles and project roles.
  A [database credential](https://docs.couchbase.com/cloud/clusters/manage-database-users.html#about-database-credentials) is specific to a database and consists of a database access name, secret, and a set of bucket and scope access levels.
  It’s required for applications to remotely authenticate on a database and access bucket data.

* [Create & manage allowlists](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/allowlist):

  More than one [allowlist](https://docs.couchbase.com/cloud/security/security.html#access-management) gives extra security across testing, development, and deployment infrastructure, and different projects.

* [Create & manage buckets](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/bucket):

  The [buckets](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html#buckets) is the top-level storage container for data in a Capella database.

* [Configure App Services](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/registry-docs/examples/appservice)

  Create and manage App Services in Capella.

* [Configure Bucket Backup & Restore](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/registry-docs/examples/backup)

  Create and manage Backups in Capella.