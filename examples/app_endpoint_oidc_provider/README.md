# Capella App Endpoint OIDC Provider Example

This example shows how to create and manage an OpenID Connect (OIDC) provider for a Capella App Endpoint.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we will do the following:

1. CREATE: Create a new OIDC provider for an App Endpoint as defined in `create_app_endpoint_oidc_provider.tf`.
2. UPDATE: Update the OIDC provider configuration using Terraform.
3. IMPORT: Import an existing OIDC provider into Terraform state.
4. DELETE: Delete the OIDC provider.

If you check the `terraform.template.tfvars` file, copy it to `terraform.tfvars` and update the values as per your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

### Apply the Plan, in order to create an OIDC provider

Command: `terraform apply`

## IMPORT
### Remove the resource from the Terraform State file

Command: `terraform state rm couchbase-capella_app_endpoint_oidc_provider.example_oidc_provider`

### Now, import the resource into Terraform

Command: `terraform import couchbase-capella_app_endpoint_oidc_provider.example_oidc_provider organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,app_endpoint_name=<app_endpoint_name>,provider_id=<provider_id>`

Here, we pass the identifiers as a single comma-separated string, with the last value being the `provider_id` of the OIDC provider.

## UPDATE
### Edit `terraform.tfvars` to update provider settings and apply

Command: `terraform apply`

## DELETE
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`
