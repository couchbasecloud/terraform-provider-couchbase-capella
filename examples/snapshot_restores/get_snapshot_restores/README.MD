# Get Existing Cloud Snapshot Restore

This example shows how to get a cloud snapshot backup restore that already exists in Capella for a given cluster. It uses the organization ID, project ID, cluster ID, and restore ID to do so. 

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
data.couchbase-capella_cloud_snapshot_restore.cloud_snapshot_restore: Reading...
data.couchbase-capella_cloud_snapshot_restore.cloud_snapshot_restore: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Changes to Outputs:
  + cloud_snapshot_restore = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at      = "2025-10-24T08:38:11.629258221Z"
      + id              = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + restore_to      = "ffffffff-aaaa-1414-eeee-000000000000"
      + snapshot        = "2025-10-24T08:32:13.56530818Z"
      + status          = "complete"
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
data.couchbase-capella_cloud_snapshot_restore.cloud_snapshot_restore: Reading...
data.couchbase-capella_cloud_snapshot_restore.cloud_snapshot_restore: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Changes to Outputs:
  + cloud_snapshot_restore = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at      = "2025-10-24T08:38:11.629258221Z"
      + id              = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + restore_to      = "ffffffff-aaaa-1414-eeee-000000000000"
      + snapshot        = "2025-10-24T08:32:13.56530818Z"
      + status          = "complete"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

cloud_snapshot_restore = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "created_at" = "2025-10-24T08:38:11.629258221Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "restore_to" = "ffffffff-aaaa-1414-eeee-000000000000"
  "snapshot" = "2025-10-24T08:32:13.56530818Z"
  "status" = "complete"
}
```