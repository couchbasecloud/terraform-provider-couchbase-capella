# The Couchbase Capella Terraform Provider

The Capella Terraform Provider is a powerful way of programmatically managing Capella API keys, users, organizations, projects, clusters, buckets, and other resources.


## Example Usage

```terraform
# Configure the Couchbase Capella Provider 
provider "couchbase-capella" {
  authentication_token = var.couchbasecapella_auth_token
}
# Create the resources
```

## Configure Capella Access

For authentication with the Couchbase Capella provider a [V4 REST API Key](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#understand-management-api-keys) must be [generated](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#generate-management-api-keys).
An [Organization Owner](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html) role should be sufficient to get you going --- although you may wish to review this level of access, and use a different key in production.
This API key is then used for authenticating the Terraform Provider against Couchbase Capella.


**IMPORTANT**  
Although in examples below, the API key secret is specified as an environmental variable and hardcoded in a config file, it is recommended that the API secret credentials be stored in a remote secrets manager such as Hashicorp Vault or AWS Secrets Manager that the Terraform Provider can then retrieve and use for authentication.


### Terraform Environment Variables

The provider configuration block accepts the following arguments. In most cases it is recommended to set them via the indicated environment variables in order to keep credential information out of the configuration.

* `host` - Allows you to specify the host for the capella API. This can be useful if you are behind a reverse proxy. May be set vai the `CAPELLA_HOST` environment variable. If not provided will default to "https://cloudapi.cloud.couchbase.com"
* `authentication_token` - A valid [V4 REST API Key](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#understand-management-api-keys) for authenticating with the Couchbase Capella API. May be set via the `CAPELLA_AUTHENTICATION_TOKEN` environment variable.

## Create and manage resources using terraform

### Example Usage

Note: You will need to provide the V4 API secret for authentication.

```terraform
terraform {
  required_providers {
    couchbase-capella = {
      source = "couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  authentication_token = "capella authentication token"
}


resource "couchbase-capella_project" "project" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  name = "example-name"
  description = "example-description"
}

output "example_project" {
  value = couchbase-capella_project.project
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
$ terraform destroy -target=couchbase-capella_project.project
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

To get started, see the [Provider Example Configs](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/getting_started):

* [Retrieve Organization Details in Capella](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/organization):

  Couchbase Capella uses an ordered hierarchy to help you keep all of your data organized and securely accessible.
  The entity at the top of the hierarchy is called an organization.
  Everything you do in Capella, from creating a cluster to managing your billing, happens inside an [organization](https://docs.couchbase.com/cloud/organizations/organizations.html).

* [Create and Manage Users](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/user):

  Users have roles within an [organization](https://docs.couchbase.com/cloud/organizations/manage-organization-users.html), and within [individual projects](https://docs.couchbase.com/cloud/projects/manage-project-users.html).

* [Create and Manage API Keys](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/apikey):

  Every API key is associated with an allowed IP Address list, and one or more organization roles, which determine the [privileges that the API key has](https://docs.couchbase.com/cloud/management-api-guide/management-api-start.html#understand-management-api-keys) within the organization.

* [Create & Manage Projects](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/project):

  Within organizations, [projects](https://docs.couchbase.com/cloud/projects/projects.html) are used to organize and manage groups of Couchbase databases.
  An organization can contain any number of projects, and a project can contain any number of databases.

* [Create & Manage Capella Clusters (databases)](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/cluster):

  The Cluster is the individual instance of a [Couchbase Database](https://docs.couchbase.com/cloud/clusters/databases.html), spanning one or more nodes on your Cloud Service Provider, and containing the Data Service, and any other services which you choose to deploy.
  Within this sits the hierarchy of bucket, scope, collection, and document.

* [Retrieve Cluster Certificate Details](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/certificate):

  Retrieve the certificate details for a Capella cluster;
  list the certificate details based on the cluster ID and authentication access token.

* [Manage Database Credentials](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/database_credential):

  Database credentials are separate from organization roles and project roles.
  A [database credential](https://docs.couchbase.com/cloud/clusters/manage-database-users.html#about-database-credentials) is specific to a database and consists of a database access name, secret, and a set of bucket and scope access levels.
  It's required for applications to remotely authenticate on a database and access bucket data.

* [Create & Manage Allowlists](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/allowlist):

  More than one [allowlist](https://docs.couchbase.com/cloud/security/security.html#access-management) gives extra security across testing, development, and deployment infrastructure, and different projects.

* [Create & Manage App Services Allowed CIDRs](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/app_services_cidr):

  App Services only allow trusted IP addresses to connect and use its REST APIs.
  Each App Service has a list of allowed IP addresses that can include up to 75 entries. 
  Each entry can be a single IP address or an IP address space. The addresses are written in CIDR notation.
  Any IP address you add to this list can have a user-specified expiration time for temporary access, or be permanent. 
  Capella automatically denies any connection attempts to and from an IP not in the allowed IP list.

* [Configure App Services](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/appservice):

  App Services synchronizes data between the Couchbase Capella database and your apps running on mobile applications.

* [Configure App Endpoints](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_endpoint):

  App Endpoints are logical entities that define how Couchbase Lite clients connect to and sync data with your Couchbase cluster through App Services. 
  Each App Endpoint represents a distinct configuration that controls authentication, data access, CORS policies, and sync behavior for a specific set of clients.

* [Manage App Endpoint Activation Status](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_endpoint_activation_status):

  Control the activation state of App Endpoints. Endpoints can be activated or deactivated to manage client connectivity.

* [Configure App Endpoint CORS](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_endpoint_cors):

  Configure Cross-Origin Resource Sharing (CORS) settings for App Endpoints to control which web applications can access your endpoint.

* [Configure App Endpoint OIDC Authentication](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_endpoint_oidc_provider):

  Set up OpenID Connect (OIDC) authentication providers for App Endpoints to enable secure user authentication.

* [Set App Endpoint Default OIDC Provider](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_endpoint_default_oidc_provider):

  Configure the default OIDC provider for an App Endpoint when multiple providers are configured.

* [Manage App Endpoint Access Control Functions](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_endpoint_access_control_function):

  Define custom JavaScript access control functions to implement fine-grained data access policies for App Endpoints.

* [Configure App Endpoint Import Filters](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/import_filter):

  Set up import filters to control which documents are synced from the Couchbase Server to mobile clients.

* [Create & Manage Buckets](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/main/examples/bucket):

  The [buckets](https://docs.couchbase.com/cloud/clusters/data-service/about-buckets-scopes-collections.html#buckets) is the top-level storage container for data in a Capella database.

* [Configure Bucket Backup & Restore](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/backup):

  A bucket is the fundamental space for storing data in Couchbase Capella.
  Create and Manage Backups in Capella.

* [Create & Manage Scopes](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/bucket/scopes):

  Create and Manage Scopes within a bucket. A scope is a mechanism for the grouping of multiple collections.

* [Create & Manage Collections](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/bucket/collections):

  Create and Manage Collections within a bucket. A collection is a data container for related documents.

* [Configure Cluster On/Off Schedule](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/cluster_onoff_schedule):

  The On/Off Schedule endpoint enables you to schedule when your provisioned database should turn on or off to save costs.

* [Turn Clusters On/Off On Demand](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/cluster_onoff_ondemand):

  Turning off your database only turns off the compute. All of your data, schema (buckets, scopes, and collections), and indexes remain, as well as your cluster configuration, including users and allow lists.
  When you turn your provisioned database off, you will be charged the OFF amount for the database.

* [Turn App Service On/Off On Demand](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/app_service_onoff_ondemand):

  You can turn the cluster and any linked app services on or off on demand using the [cluster API](https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/clusters).

* [Import Sample Dataset Buckets](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/sample_bucket):

  The sampleBucket endpoint lets users easily create a bucket filled with sample data. This is a quick way for users to try out features and learn how things work with ready-to-use data.

* [Manage Storage Auto Expansion for Azure Clusters](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/cluster/azure):

  Manage Storage Auto-expansion for Azure Clusters.

* [Manage Audit Log Settings in Capella](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/audit_log_settings):

  Users can easily configure audit log support on Capella database.

* [Export Capella Audit Logs](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/audit_log_export):

  Users can export audit logs from cloud blob storage to an AWS S3 bucket.

* [Retrieve Audit Log Events for a Cluster](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/audit_logs_event_ids):

  Users can retrieve audit logs from a pre-signed download URL. Logs are retained for 30 days.

* [Manage Network Peer](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/network_peer):

  Network Peering enables you to configure a secure private network connection between the Virtual Private Cloud (VPC) hosting your applications and the VPC of your Couchbase Capella database.

* [Retrieve a Specific Event](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/event):

  Fetch the details of any specific event using the eventId and organizationId. The results are always limited by the role and scope of the caller's privileges.

* [Retrieve All Events](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/events):

  List the information of all the events within an organization. The list can be customized using filters.

* [Retrieve Private Endpoint Command for AWS](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/private_endpoint_command/AWS):

  Retrieve the AWS command used to configure a VPC endpoint.

* [Retrieve Private Endpoint Command for GCP](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/private_endpoint_command/GCP):

  Retrieve the GCP command used to configure a VPC endpoint.

* [Retrieve Private Endpoint Command for Azure](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/private_endpoint_command/Azure):

  Retrieve the Azure command used to configure a private endpoint.

* [Manage Private Endpoint Service](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/private_endpoint_service):

  Enable or disable the Private Endpoint Service in Capella.

* [Manage Private Endpoints](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/private_endpoints):

  Accept or reject a private endpoint in Capella.

* [Manage Azure Network Peer](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/network_peer):

  Configure a secure private network connection between the Azure VNet hosting your applications and the VNet of your Couchbase Capella operational cluster.

* [Manage Flush Bucket](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/flush_bucket):

  Flush a bucket to have the system delete all its data at the earliest opportunity available.

* [Manage GSI](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/gsi):

  Manage secondary indexes on a Capella operational cluster. It's recommended to use deferred index builds, especially for larger indexes.

* [Create & Manage Free Tier Clusters](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/free_tier_cluster):

  Create and manage a free tier operational cluster in Capella. Free tier operational clusters provide a cost-effective way to get started with Couchbase Capella. This is a Single Node cluster that runs only the Data, Query, Index, and Search services. Only 1 free tier operational cluster is available per organization.

* [Create & Manage Free Tier Buckets](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/free_tier_bucket):

  Create and manage buckets within a free tier operational cluster. Free tier buckets have specific memory allocation limits. This is a Couchbase bucket where only the name and memory quota is configurable. Other bucket properties use default values.

* [Create & Manage Free Tier App Services](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/free_tier_appservice):

  Set up and manage App Services for your free tier operational cluster. This is a Single Node App Service that can only be linked to a free tier operational cluster. The App Service can only be turned off and on when the linked free tier operational cluster is turned off and on.

* [Manage Free Tier Cluster On/Off](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/free_tier_cluster_on_off):

  Turn your free tier operational cluster on or off on demand to optimize resource usage. Turning the free tier cluster on or off will also turn on or off any linked App Services. 
