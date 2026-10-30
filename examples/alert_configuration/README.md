# Capella Alert Configuration Example

This example shows how to manage the alert configuration of a Capella cluster, and how to read the alert configuration of an organization, project and cluster.

Alert configuration is inherited. Capella has a default for every alert, an organization can override it, a project can override its organization, and a cluster can override its project. For each level, the data sources return:

- what the level sets itself: `configs`, `recipients` and `muted_alerts`;
- what is in effect at the level after inheritance: `resolved_configs` and `resolved_muted_alerts`.

`configs` and `muted_alerts` are maps keyed by alert identifier, for example `cpu_usage_high`.

A resolved config can combine fields set at different levels, so `resolved_configs` entries only carry the values, without `updated_at` or `updated_by`. Use `configs` at each level to see who changed what.

The `couchbase-capella_cluster_alert_configuration` resource manages only what the cluster sets for itself. Manage each cluster with only one resource. Each attribute follows its own rule about what Terraform owns:

- `configs`: only the alerts listed, and only the fields set on each. Alerts configured outside Terraform are left alone.
- `muted_alerts`: every mute at the cluster, once the attribute is set. A mute added in the Capella UI shows up in the next plan, and `terraform apply` clears it.
- `recipients`: the whole list, if set. Removing the attribute leaves the list as it is.

Capella cannot remove a threshold or frequency override. Removing an alert or a field from `configs` leaves its value in place, and `terraform apply` warns which values were left. On destroy, Terraform clears the mutes and recipients it manages and warns which thresholds and frequencies stay.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Set the alert configuration of a cluster, as stated in the `create_cluster_alert_configuration.tf` file.
2. UPDATE: Remove an alert and a field from the cluster's `configs`.
3. GET: Read and display the alert configuration of an organization, a project and a cluster, as stated in the `get_organization_alert_configuration.tf`, `get_project_alert_configuration.tf` and `get_cluster_alert_configuration.tf` files.
4. IMPORT: Import the alert configuration of a cluster.
5. DELETE: Clear the mutes and recipients the resource manages, and remove it and the data outputs from terraform state.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
  # couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration will be created
  + resource "couchbase-capella_cluster_alert_configuration" "new_cluster_alert_configuration" {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000003"
      + configs         = {
          + "CBHardOutOfMemoryErrors" = {
              + frequency_hours = 1
            },
          + "cpu_usage_high" = {
              + critical_threshold = 90
              + frequency_hours    = 3
              + warning_threshold  = 80
            },
        }
      + muted_alerts    = {
          + "disk_usage_high" = {
              + action      = "mute"
              + muted_until = "2027-01-01T00:00:00Z"
            },
        }
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
      + recipients      = [
          + "oncall@example.com",
        ]
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the Plan, in order to set the alert configuration

Command: `terraform apply`

```
couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration: Creating...
couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

If the cluster already had mutes set in the Capella UI, the apply clears them, because the config sets `muted_alerts`. Alerts configured in the UI that are not in `configs` stay as they are.

## UPDATE
### Remove an alert and a field from configs

Remove `CBHardOutOfMemoryErrors` and `critical_threshold` from `create_cluster_alert_configuration.tf`, then run `terraform apply`. Capella keeps both values, so nothing is sent, and Terraform warns:

```
couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration: Modifying...
couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration: Modifications complete after 1s

Warning: Alert Configuration Left In Capella

Capella cannot remove a threshold or frequency override, so the cluster keeps
these values and Terraform no longer manages them:
  - CBHardOutOfMemoryErrors: frequency_hours = 1
  - cpu_usage_high: critical_threshold = 90

To change a value, set the value you want in configs instead of removing it.

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## GET
### View the alert configuration

Command: `terraform plan`

Sample Output for the cluster, where the cluster sets its own CPU frequency and inherits the thresholds and a mute:
```
data.couchbase-capella_cluster_alert_configuration.existing_cluster_alert_configuration: Reading...
data.couchbase-capella_cluster_alert_configuration.existing_cluster_alert_configuration: Read complete after 1s

Changes to Outputs:
  + existing_cluster_alert_configuration = {
      + cluster_id            = "ffffffff-aaaa-1414-eeee-000000000003"
      + configs               = {
          + cpu_usage_high = {
              + critical_threshold = null
              + frequency_hours    = 3
              + threshold          = null
              + updated_at         = "2026-09-16T10:50:57Z"
              + updated_by         = "ffffffff-aaaa-1414-eeee-000000000001"
              + warning_threshold  = null
            }
        }
      + muted_alerts          = {}
      + organization_id       = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id            = "ffffffff-aaaa-1414-eeee-000000000002"
      + recipients            = []
      + resolved_configs      = {
          + cpu_usage_high = {
              + critical_threshold = 90
              + frequency_hours    = 3
              + threshold          = null
              + warning_threshold  = 80
            }
        }
      + resolved_muted_alerts = {
          + disk_usage_high = {
              + action      = "mute"
              + muted_from  = null
              + muted_until = "2026-11-01T00:00:00Z"
              + updated_at  = "2026-09-20T09:58:12.104Z"
              + updated_by  = "Jane Doe"
            }
        }
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

The organization and project outputs have the same shape, without `cluster_id` (organization and project) and `project_id` (organization).

### Apply the Plan, in order to get the alert configuration.

Command: `terraform apply`

Running `terraform apply` shows the same plan, asks for confirmation, and then prints the three outputs.

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

## IMPORT
### Import the alert configuration of a cluster

Command: `terraform import couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Import only sets the IDs. The first `terraform apply` after import sends the values in your config.

## DELETE
### Destroy the resources

Command: `terraform destroy`

Destroying the resource clears the cluster's mutes and recipients, so the cluster goes back to inheriting them, and warns which thresholds and frequencies stay. The data sources create nothing in Capella, so destroying them only removes the outputs from the Terraform state.

```
couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration: Destroying...
couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration: Destruction complete after 1s

Warning: Alert Configuration Left In Capella

Capella cannot remove a threshold or frequency override, so the cluster keeps
these values after destroy:
  - cpu_usage_high: frequency_hours = 3, warning_threshold = 80

Destroy complete! Resources: 1 destroyed.
```
