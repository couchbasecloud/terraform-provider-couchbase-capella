
# Capella App Endpoint Default OIDC Provider Example

This example shows how to set and manage the default OIDC provider for an App Endpoint in Capella.

It uses the organization ID, project ID, cluster ID, app service ID, app endpoint name, and the `provider_id` to configure the default provider.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following:

1. CREATE: Set the default OIDC provider as defined in `create_app_endpoint_default_oidc_provider.tf`.
2. READ: Fetch current default OIDC provider.
3. UPDATE: Change the default OIDC provider.
4. IMPORT: Import an existing default OIDC provider configuration into Terraform state.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the variable values for your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

### Apply the plan to set the default OIDC provider

Command: `terraform apply`

## READ
### Show current state

Command: `terraform show`

## UPDATE
### Change the default provider ID

Edit `default_provider_id` in `terraform.tfvars` or update the resource block, then run:

Command: `terraform apply`

## IMPORT
### Import an existing default OIDC provider configuration into Terraform state

Ensure the resource block has the correct IDs configured, then run:

Command: `terraform import couchbase-capella_app_endpoint_default_oidc_provider.example organization_id=<org_id>,project_id=<proj_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,app_endpoint_name=<app_endpoint_name>`



