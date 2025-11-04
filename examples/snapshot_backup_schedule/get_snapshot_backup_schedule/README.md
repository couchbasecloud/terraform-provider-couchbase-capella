# Get Existing Cloud Snapshot Backup Schedule

This example shows how to get a cloud snapshot backup schedule that already exists in Capella for a given cluster. It uses the organization ID, project ID, and cluster ID to do so. 

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Get
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule: Reading...
data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule: Read complete after 0s

Changes to Outputs:
  + existing_cloud_snapshot_backup_schedule = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + copy_to_regions = [
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      + interval        = 12
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention       = 24
      + start_time      = "2025-10-23T03:30:00Z"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

### Apply the Plan

```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule: Reading...
data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule: Read complete after 0s

Changes to Outputs:
  + existing_cloud_snapshot_backup_schedule = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + copy_to_regions = [
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      + interval        = 12
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention       = 24
      + start_time      = "2025-10-23T03:30:00Z"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_cloud_snapshot_backup_schedule = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "copy_to_regions" = toset([
    "ap-southeast-1",
    "eu-west-1",
  ])
  "interval" = 12
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "retention" = 24
  "start_time" = "2025-10-23T03:30:00Z"
}
```