# Capella App Service Maintenance Window Example

This example shows how to manage the preferred weekly maintenance window of an App Service in Capella.

It sets the maintenance window of an existing App Service and reads it back with the data source. All times are in UTC. The maintenance window endpoints are only available to organizations enabled for self-service upgrades.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Set the maintenance window of the App Service as stated in the `create_app_service_maintenance_window.tf` file.
2. GET: Read the maintenance window with the data source as stated in the `get_app_service_maintenance_window.tf` file.
3. UPDATE: Change the preferred day, hour or minute using Terraform.
4. IMPORT: Import a maintenance window that exists in Capella but not in the terraform state file.
5. DELETE: Remove the maintenance window from the App Service.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values of the variables for your organization.

An App Service has at most one maintenance window. If one is already set (for example through the Capella UI), `terraform apply` fails with an error that tells you to import it instead of overwriting it.

## CREATE & GET

Command: `terraform plan`, then `terraform apply`

The resource and the data source outputs show `preferred_day`, `preferred_hour` and `preferred_minute`. Until the window exists, the data source fields are null.

## UPDATE

Change `maintenance_window` in `terraform.tfvars`, for example to `preferred_day = "wednesday"`, `preferred_hour = 23`, `preferred_minute = 59`, then run `terraform apply`. The window is updated in place. Changing any of the IDs replaces the resource.

## IMPORT

Remove the resource from the state file, then import it:

Command: `terraform state rm couchbase-capella_app_service_maintenance_window.new_app_service_maintenance_window`

Command: `terraform import couchbase-capella_app_service_maintenance_window.new_app_service_maintenance_window organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>`

Running `terraform plan` afterwards shows no changes.

## DELETE

Command: `terraform destroy`

This removes the maintenance window from the App Service. The App Service itself is not changed.
