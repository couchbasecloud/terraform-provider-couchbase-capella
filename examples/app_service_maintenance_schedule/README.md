# Capella App Service Maintenance Schedule Example

This example shows how to reschedule an App Service maintenance in Capella.

It lists the maintenance of an existing App Service with the `couchbase-capella_app_service_maintenances` data source, then sets the queue time of the pending maintenance you choose with the `couchbase-capella_app_service_maintenance_schedule` resource. The resource does not create a maintenance; it only changes when an existing one starts.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. LIST: List the maintenance of the App Service as stated in the `list_app_service_maintenances.tf` file.
2. CREATE: Reschedule a maintenance as stated in the `create_app_service_maintenance_schedule.tf` file.
3. UPDATE: Change the queue time using Terraform.
4. IMPORT: Import the schedule of a maintenance that is not in the terraform state file.
5. DELETE: Remove the resource from the terraform state file.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values of the variables for your organization.

## LIST

Command: `terraform apply -target=data.couchbase-capella_app_service_maintenances.existing_maintenances`

The `app_service_maintenances` output lists every maintenance with its `id`, `kind`, `initiated_by`, `current_state`, `queue_at` and `reschedule_by`. The `reschedulable_maintenance_ids` output lists the IDs of the pending ones, which are the only ones that can be rescheduled. If it is empty, the App Service has nothing to reschedule.

## CREATE

Set `maintenance_id` in `terraform.tfvars` to one of the `reschedulable_maintenance_ids`, then run `terraform plan` and `terraform apply`.

The resource adopts that maintenance and, if its `queue_at` differs from `var.queue_at`, reschedules it. The ID is pinned in `terraform.tfvars` so later plans keep managing the same maintenance after it starts or completes.

## UPDATE

Change `queue_at` in `terraform.tfvars` to another time in the future and no later than the maintenance's `reschedule_by`, then run `terraform apply`. The maintenance is rescheduled in place. Capella rejects the change if the maintenance is no longer pending or the time is outside the window. Changing any of the IDs replaces the resource.

## IMPORT

Remove the resource from the state file, then import it:

Command: `terraform state rm couchbase-capella_app_service_maintenance_schedule.new_app_service_maintenance_schedule`

Command: `terraform import couchbase-capella_app_service_maintenance_schedule.new_app_service_maintenance_schedule organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,maintenance_id=<maintenance_id>`

Running `terraform plan` afterwards shows no changes.

## DELETE

Command: `terraform destroy`

This only removes the resource from the state file. The maintenance stays in Capella with its last queue time.
