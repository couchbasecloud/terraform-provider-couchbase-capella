# Get Existing Cloud Snapshot Backup

This example shows how to get a cloud snapshot backup that already exists in Capella for a given cluster. It uses the organization ID, project ID, cluster ID, and backup ID to do so. 

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
data.couchbase-capella_cloud_snapshot_backup.cloud_snapshot_backup: Reading...
data.couchbase-capella_cloud_snapshot_backup.cloud_snapshot_backup: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Changes to Outputs:
  + cloud_snapshot_backup = {
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + cmek                = []
      + created_at          = "2025-10-27T11:08:27.566341709Z"
      + cross_region_copies = [
          + {
              + region_code = "ap-southeast-1"
              + status      = "complete"
              + time        = "2025-10-27T11:10:55.137154625Z"
            },
          + {
              + region_code = "eu-west-1"
              + status      = "complete"
              + time        = "2025-10-27T11:09:53.14435443Z"
            },
        ]
      + expiration          = "2025-11-03T11:08:27.566341709Z"
      + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + progress            = {
          + status = "complete"
          + time   = "2025-10-27T11:09:20.067713345Z"
        }
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention           = 168
      + server              = {
          + version = "7.6.7"
        }
      + size                = 0
      + type                = "on_demand"
    }
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
data.couchbase-capella_cloud_snapshot_backup.cloud_snapshot_backup: Reading...
data.couchbase-capella_cloud_snapshot_backup.cloud_snapshot_backup: Read complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Changes to Outputs:
  + cloud_snapshot_backup = {
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + cmek                = []
      + created_at          = "2025-10-27T11:08:27.566341709Z"
      + cross_region_copies = [
          + {
              + region_code = "ap-southeast-1"
              + status      = "complete"
              + time        = "2025-10-27T11:10:55.137154625Z"
            },
          + {
              + region_code = "eu-west-1"
              + status      = "complete"
              + time        = "2025-10-27T11:09:53.14435443Z"
            },
        ]
      + expiration          = "2025-11-03T11:08:27.566341709Z"
      + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + progress            = {
          + status = "complete"
          + time   = "2025-10-27T11:09:20.067713345Z"
        }
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention           = 168
      + server              = {
          + version = "7.6.7"
        }
      + size                = 0
      + type                = "on_demand"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

cloud_snapshot_backup = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cmek" = toset([])
  "created_at" = "2025-10-27T11:08:27.566341709Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-10-27T11:10:55.137154625Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-10-27T11:09:53.14435443Z"
    },
  ])
  "expiration" = "2025-11-03T11:08:27.566341709Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "progress" = {
    "status" = "complete"
    "time" = "2025-10-27T11:09:20.067713345Z"
  }
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "retention" = 168
  "server" = {
    "version" = "7.6.7"
  }
  "size" = 0
  "type" = "on_demand"
}
```