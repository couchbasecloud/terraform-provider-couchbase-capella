# Capella App Endpoint Default OIDC Provider Example

This example shows how to set the default OpenID Connect (OIDC) provider for a Capella App Endpoint.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we will do the following:

1. CREATE: Set the default OIDC provider for an App Endpoint as defined in `set_default_oidc_provider.tf`.
2. UPDATE: Update the default provider by changing `provider_id`.
3. IMPORT: Import the default OIDC provider state into Terraform.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values as per your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

### Apply the Plan, in order to set the default OIDC provider

Command: `terraform apply`

## IMPORT
### Remove the resource from the Terraform State file

Command: `terraform state rm couchbase-capella_app_endpoint_default_oidc_provider.example_default`

### Now, import the resource into Terraform

Command: `terraform import couchbase-capella_app_endpoint_default_oidc_provider.example_default organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,app_endpoint_name=<app_endpoint_name>`

Here, we pass the identifiers as a single comma-separated string. The current default provider is derived by the provider during Read.

## UPDATE
### Edit `terraform.tfvars` to change `provider_id` and apply

Command: `terraform apply`

## DELETE
This resource represents a selection; the underlying API does not support unsetting the default provider. Deleting this resource only removes it from Terraform state.
