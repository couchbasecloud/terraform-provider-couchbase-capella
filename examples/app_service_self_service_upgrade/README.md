# Capella App Service Self-Service Upgrade Example

This example shows how to upgrade an App Service in Capella, either now or at a scheduled time.

It lists the versions the App Service can be upgraded to with the `couchbase-capella_app_service_self_service_upgrade_available_versions` data source, creates an upgrade with the `couchbase-capella_app_service_self_service_upgrade` resource, and lists the App Service's self-service upgrades with the `couchbase-capella_app_service_self_service_upgrades` data source. Self-service upgrades are only available to organizations enabled for them.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. LIST VERSIONS: List the available versions as stated in the `list_app_service_self_service_upgrade_available_versions.tf` file.
2. CREATE: Schedule an upgrade as stated in the `create_app_service_self_service_upgrade.tf` file.
3. UPDATE: Change the version or queue time using Terraform.
4. LIST UPGRADES: List the self-service upgrades as stated in the `list_app_service_self_service_upgrades.tf` file.
5. IMPORT: Import an upgrade that is not in the terraform state file.
6. DELETE: Cancel the upgrade.

Copy `terraform.template.tfvars` to `terraform.tfvars` and update the values of the variables for your organization.

## LIST VERSIONS

Command: `terraform apply -target=data.couchbase-capella_app_service_self_service_upgrade_available_versions.available_versions`

The `available_version_ids` output lists the versions the App Service can be upgraded to. If it is empty, there is nothing to upgrade to.

## CREATE

Set `version_id` in `terraform.tfvars` to one of the `available_version_ids`, and `queue_at` to a time in the future and no more than 7 days ahead. Then run `terraform plan` and `terraform apply`.

Leaving `queue_at` unset upgrades the App Service now, and `queue_at` is then read back from Capella.

## UPDATE

Change `version_id` or `queue_at` in `terraform.tfvars`, then run `terraform apply`. The upgrade is changed in place. Capella rejects the change once the upgrade is no longer pending. Changing any of the other IDs replaces the resource.

## LIST UPGRADES

Command: `terraform apply -target=data.couchbase-capella_app_service_self_service_upgrades.existing_upgrades`

The `app_service_self_service_upgrades` output lists every self-service upgrade of the App Service with its `id`, `version_id`, `current_state` and `queue_at`. Maintenance initiated by Couchbase is not included.

## IMPORT

Remove the resource from the state file, then import it:

Command: `terraform state rm couchbase-capella_app_service_self_service_upgrade.new_app_service_self_service_upgrade`

Command: `terraform import couchbase-capella_app_service_self_service_upgrade.new_app_service_self_service_upgrade id=<upgrade_id>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Running `terraform plan` afterwards shows no changes.

## DELETE

Command: `terraform destroy`

A pending upgrade is cancelled. An upgrade that has completed, failed or been cancelled is only removed from the state file. An upgrade that has already started cannot be cancelled, so destroy fails and keeps it in state until it finishes.
