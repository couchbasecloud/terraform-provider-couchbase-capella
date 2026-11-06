# Capella Alerts Example

This example shows how to list the firing and resolved alerts in a Capella organization.

There is one entry for each alert on each resource, so an alert that fires for each node or bucket has one entry for each of them. Muted alerts are not included.

By default the list covers the last 48 hours. Use the optional `from` and `to` fields to read a different time range, up to 31 days back, and `status`, `cluster_ids`, `project_ids` and `severity_levels` to narrow the list.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the alerts as stated in the `list_alerts.tf` file.
2. DELETE: Delete the alerts data output from terraform state.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## GET
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/terraform-provider-capella
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_alerts.existing_alerts: Reading...
data.couchbase-capella_alerts.existing_alerts: Read complete after 1s

Changes to Outputs:
  + existing_alerts = {
      + cluster_ids     = null
      + data            = [
          + {
              + alert_id      = "ffffffff-aaaa-1414-eeee-000000000003:CPUUsageHighWarning"
              + bucket        = null
              + ends_at       = null
              + event_id      = "ffffffff-aaaa-1414-eeee-000000000009"
              + event_time    = "2026-10-12T10:41:00Z"
              + kind          = "CPUUsageHighWarning"
              + name          = "High CPU Usage"
              + node_hostname = "svc-d-node-001.example.cloud.couchbase.com"
              + project_id    = "ffffffff-aaaa-1414-eeee-000000000002"
              + project_name  = "Payments"
              + resource_id   = "ffffffff-aaaa-1414-eeee-000000000003"
              + resource_name = "orders-prod"
              + resource_type = "cluster"
              + severity      = "warning"
              + starts_at     = "2026-10-12T09:41:00Z"
              + status        = "firing"
            },
          + {
              + alert_id      = "ffffffff-aaaa-1414-eeee-000000000003:DiskUsageHighWarning"
              + bucket        = null
              + ends_at       = "2026-10-12T08:15:00Z"
              + event_id      = "ffffffff-aaaa-1414-eeee-000000000011"
              + event_time    = "2026-10-12T08:15:00Z"
              + kind          = "DiskUsageHighWarning"
              + name          = "High Disk Usage"
              + node_hostname = "svc-d-node-001.example.cloud.couchbase.com"
              + project_id    = "ffffffff-aaaa-1414-eeee-000000000002"
              + project_name  = "Payments"
              + resource_id   = "ffffffff-aaaa-1414-eeee-000000000003"
              + resource_name = "orders-prod"
              + resource_type = "cluster"
              + severity      = "warning"
              + starts_at     = "2026-10-12T07:30:00Z"
              + status        = "resolved"
            },
        ]
      + from            = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_ids     = null
      + severity_levels = null
      + status          = null
      + to              = null
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

### Apply the Plan, in order to get the alerts.

Command: `terraform apply`

Running `terraform apply` shows the same plan, asks for confirmation, and then prints the `existing_alerts` output with the values above.

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

## DELETE
### Destroy the resources

Command: `terraform destroy`

The data source creates nothing in Capella, so destroying it only removes the output from the Terraform state.

```
Destroy complete! Resources: 0 destroyed.
```
