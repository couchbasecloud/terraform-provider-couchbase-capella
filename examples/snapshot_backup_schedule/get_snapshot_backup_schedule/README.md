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
data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule: Read complete after 0s [id=50b28378-37fd-4a9a-8daf-551345e3b69e]

Changes to Outputs:
  + existing_cloud_snapshot_backup_schedule = {
      + copy_to_regions = null
      + id              = "50b28378-37fd-4a9a-8daf-551345e3b69e"
      + interval        = 4
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "f35fea80-6b5c-4ebd-9c24-d564553bb21d"
      + retention       = 168
      + start_time      = "2025-10-02T13:30:00Z"
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
data.couchbase-capella_cloud_snapshot_backup_schedule.existing_cloud_snapshot_backup_schedule: Read complete after 0s [id=50b28378-37fd-4a9a-8daf-551345e3b69e]

Changes to Outputs:
  + existing_cloud_snapshot_backup_schedule = {
      + copy_to_regions = null
      + id              = "50b28378-37fd-4a9a-8daf-551345e3b69e"
      + interval        = 4
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "f35fea80-6b5c-4ebd-9c24-d564553bb21d"
      + retention       = 168
      + start_time      = "2025-10-02T13:30:00Z"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_cloud_snapshot_backup_schedule = {
  "copy_to_regions" = tostring(null)
  "id" = "50b28378-37fd-4a9a-8daf-551345e3b69e"
  "interval" = 4
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "f35fea80-6b5c-4ebd-9c24-d564553bb21d"
  "retention" = 168
  "start_time" = "2025-10-02T13:30:00Z"
}
```