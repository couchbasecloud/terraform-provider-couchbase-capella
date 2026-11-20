# Capella Alert Configuration Overrides Example

This example shows how to list the alert configuration overrides in a Capella organization.

An override is an alert that the organization, a project or a cluster configures or disables itself instead of inheriting from the level above. Each entry is one alert at one level, and carries the configuration in effect now and what the most recent change replaced.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. GET: Read and display the alert configuration overrides as stated in the `list_alert_configuration_overrides.tf` file.
2. DELETE: Delete the alert configuration overrides data output from terraform state.

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
data.couchbase-capella_alert_configuration_overrides.existing_alert_configuration_overrides: Reading...
data.couchbase-capella_alert_configuration_overrides.existing_alert_configuration_overrides: Read complete after 1s

Changes to Outputs:
  + existing_alert_configuration_overrides = {
      + data            = [
          + {
              + alert_kind          = "cpu_usage_high"
              + config              = {
                  + critical_threshold = 92
                  + frequency_hours    = 1
                  + threshold          = null
                  + updated_at         = "2026-09-16T10:50:57.793147089Z"
                  + updated_by         = "ffffffff-aaaa-1414-eeee-000000000001"
                  + warning_threshold  = 85
                }
              + muted_alert         = null
              + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
              + previous_config     = {
                  + critical_threshold = 90
                  + frequency_hours    = 3
                  + threshold          = null
                  + updated_at         = null
                  + updated_by         = null
                  + warning_threshold  = 80
                }
              + previous_recipients = null
              + project_id          = "ffffffff-aaaa-1414-eeee-000000000002"
              + recipients          = null
              + resource_id         = "ffffffff-aaaa-1414-eeee-000000000003"
              + scope               = "cluster"
            },
          + {
              + alert_kind          = "disk_usage_high"
              + config              = null
              + muted_alert         = {
                  + action      = "mute"
                  + muted_from  = null
                  + muted_until = null
                  + updated_at  = "2026-09-20T09:58:12.104Z"
                  + updated_by  = "Jane Doe"
                }
              + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
              + previous_config     = null
              + previous_recipients = null
              + project_id          = "ffffffff-aaaa-1414-eeee-000000000002"
              + recipients          = null
              + resource_id         = "ffffffff-aaaa-1414-eeee-000000000003"
              + scope               = "cluster"
            },
          + {
              + alert_kind          = "node_memory_usage"
              + config              = {
                  + critical_threshold = 90
                  + frequency_hours    = 6
                  + threshold          = null
                  + updated_at         = "2026-09-15T14:02:11.004Z"
                  + updated_by         = "ffffffff-aaaa-1414-eeee-000000000001"
                  + warning_threshold  = 80
                }
              + muted_alert         = null
              + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
              + previous_config     = null
              + previous_recipients = null
              + project_id          = "ffffffff-aaaa-1414-eeee-000000000002"
              + recipients          = [
                  + "project-oncall@example.com",
                ]
              + resource_id         = null
              + scope               = "project"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

### Apply the Plan, in order to get the alert configuration overrides.

Command: `terraform apply`

Running `terraform apply` shows the same plan, asks for confirmation, and then prints the `existing_alert_configuration_overrides` output with the values above.

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
