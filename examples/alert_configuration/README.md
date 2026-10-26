# Capella Alert Configuration Example

This example shows how to read the alert configuration of a Capella organization, project and cluster.

Alert configuration is inherited. Capella has a default for every alert, an organization can override it, a project can override its organization, and a cluster can override its project. For each level, the data sources return:

- what the level sets itself: `configs`, `recipients` and `muted_alerts`;
- what is in effect at the level after inheritance: `resolved_configs`, `resolved_recipients` and `resolved_muted_alerts`.

`configs` and `muted_alerts` are maps keyed by alert identifier, for example `cpu_usage_high`.

A resolved config can combine fields set at different levels, so `resolved_configs` entries have no `updated_at` or `updated_by`. Use `configs` at each level to see who changed what.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the alert configuration of an organization, a project and a cluster, as stated in the `get_organization_alert_configuration.tf`, `get_project_alert_configuration.tf` and `get_cluster_alert_configuration.tf` files.
2. DELETE: Delete the alert configuration data output from terraform state.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## GET
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output for the cluster, where the cluster sets its own CPU frequency and inherits the thresholds, a recipient and a mute:
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
              + updated_at         = null
              + updated_by         = null
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
      + resolved_recipients   = [
          + "org-oncall@example.com",
        ]
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

## DELETE
### Destroy the resources

Command: `terraform destroy`

The data sources create nothing in Capella, so destroying them only removes the outputs from the Terraform state.

```
Destroy complete! Resources: 0 destroyed.
```
