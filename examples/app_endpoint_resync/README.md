# Capella App Endpoint Resync Example

This example shows how to initiate and manage App Endpoint Resync operations in Capella.

It uses the organization ID, project ID, cluster ID, app service ID, and app endpoint name to start a resync. You can optionally scope the resync to specific collections per scope.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following:

1. CREATE: Initiate a resync for an App Endpoint as defined in `create_app_endpoint_resync.tf`.
2. UPDATE: Re-run a resync with updated scope/collection selections.
3. IMPORT: Import an existing resync into Terraform state (for tracking/read-back).
4. DESTROY: Stop an ongoing resync operation.

If you check the `terraform.template.tfvars` file — copy it to `terraform.tfvars` and update variable values for your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

### Apply the plan to initiate a resync

Command: `terraform apply`

Notes:
- Set `scopes` to limit the resync to specific collections per scope. Omit or set to `null` to resync everything.

## UPDATE
### Re-run resync with updated scope/collection selections

Edit `scopes` in `terraform.tfvars` or update the resource block, then:

Command: `terraform apply`

## IMPORT
### Import an existing resync into Terraform state

Ensure the resource block has the correct `organization_id`, `project_id`, `cluster_id`, and `app_service_id` configured. Then run:

Command: `terraform import couchbase-capella_app_endpoint_resync.example app_endpoint=<app_endpoint_name>`

Notes:
- The import ID is the `app_endpoint` (name). Other IDs must be provided in configuration.

## DESTROY
### Stop the resync operation managed by Terraform

Command: `terraform destroy`


