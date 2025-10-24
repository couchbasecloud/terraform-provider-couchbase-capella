# List Existing Cloud Snapshot Restores

This example shows how to list cloud snapshot restores that already exist in Capella for a given cluster. It uses the organization ID, project ID, and cluster ID to do so. 

You can filter by the snapshot restore status. In this example, we will only list the snapshot restores that are 'complete'.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## List
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
data.couchbase-capella_cloud_snapshot_restores.existing_cloud_snapshot_restores: Reading...
data.couchbase-capella_cloud_snapshot_restores.existing_cloud_snapshot_restores: Read complete after 0s

Changes to Outputs:
  + cloud_snapshot_restores_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + created_at = "2025-10-24T08:33:48.890635377Z"
              + id         = ffffffff-aaaa-1414-eeee-000000000000"
              + restore_to = "ffffffff-aaaa-1414-eeee-000000000000"
              + snapshot   = "2025-10-24T08:32:13.56530818Z"
              + status     = "complete"
            },
          + {
              + created_at = "2025-10-24T08:38:11.629258221Z"
              + id         = "ffffffff-aaaa-1414-eeee-000000000000"
              + restore_to = "ffffffff-aaaa-1414-eeee-000000000000"
              + snapshot   = "2025-10-24T08:32:13.56530818Z"
              + status     = "complete"
            },
        ]
      + filter          = {
          + name   = "status"
          + values = [
              + "complete",
            ]
        }
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
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
data.couchbase-capella_cloud_snapshot_restores.existing_cloud_snapshot_restores: Reading...
data.couchbase-capella_cloud_snapshot_restores.existing_cloud_snapshot_restores: Read complete after 0s

Changes to Outputs:
  + cloud_snapshot_restores_list = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + created_at = "2025-10-24T08:33:48.890635377Z"
              + id         = "ffffffff-aaaa-1414-eeee-000000000000"
              + restore_to = "ffffffff-aaaa-1414-eeee-000000000000"
              + snapshot   = "2025-10-24T08:32:13.56530818Z"
              + status     = "complete"
            },
          + {
              + created_at = "2025-10-24T08:38:11.629258221Z"
              + id         = "ffffffff-aaaa-1414-eeee-000000000000"
              + restore_to = "ffffffff-aaaa-1414-eeee-000000000000"
              + snapshot   = "2025-10-24T08:32:13.56530818Z"
              + status     = "complete"
            },
        ]
      + filter          = {
          + name   = "status"
          + values = [
              + "complete",
            ]
        }
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

cloud_snapshot_restores_list = {
  "cluster_id" = "34dd72af-ce0f-481d-a197-e0fdd9296966"
  "data" = tolist([
    {
      "created_at" = "2025-10-24T08:33:48.890635377Z"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "restore_to" = "ffffffff-aaaa-1414-eeee-000000000000"
      "snapshot" = "2025-10-24T08:32:13.56530818Z"
      "status" = "complete"
    },
    {
      "created_at" = "2025-10-24T08:38:11.629258221Z"
      "id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "restore_to" = "ffffffff-aaaa-1414-eeee-000000000000"
      "snapshot" = "2025-10-24T08:32:13.56530818Z"
      "status" = "complete"
    },
  ])
  "filter" = {
    "name" = "status"
    "values" = toset([
      "complete",
    ])
  }
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```