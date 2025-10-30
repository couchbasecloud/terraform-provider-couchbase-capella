# List Existin Project Level Cloud Snapshot Backups

This example shows how to list cloud snapshot backups that already exist in Capella for a given cluster. It uses the organization ID and project ID to do so. 

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## List
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/laura.silaja/code/Lagher0/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Reading...
data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Read complete after 1s

Changes to Outputs:
  + project_backups_list = {
      + cursor          = {
          + hrefs = {
              + first    = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=10"
              + last     = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=10"
              + next     = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=0&perPage=10"
              + previous = ""
            }
          + pages = {
              + last        = 1
              + next        = 0
              + page        = 1
              + per_page    = 10
              + previous    = 0
              + total_items = 1
            }
        }
      + data            = [
          + {
              + cloud_provider       = "hostedAWS"
              + cluster_id           = "ffffffff-AAAA-1414-eeee-000000000000"
              + cluster_name         = "goldavinashkak"
              + created_by           = "laura.silaja@couchbase.com"
              + creation_date_time   = "2025-10-29T14:57:09.838906045Z"
              + current_status       = "healthy"
              + most_recent_snapshot = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-29T15:01:04.912087543Z"
                  + cross_region_copies = []
                  + database_size       = 0
                  + expiration          = "2025-11-05T15:01:04.912087543Z"
                  + id                  = "ffffffff-aaaa-1414-EEEE-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-29T15:02:21.564367676Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 168
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + oldest_snapshot      = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-29T15:01:04.912087543Z"
                  + cross_region_copies = []
                  + database_size       = 0
                  + expiration          = "2025-11-05T15:01:04.912087543Z"
                  + id                  = "ffffffff-aaaa-1414-EEEE-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-29T15:02:21.564367676Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 168
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + region               = "us-east-1"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = null
      + per_page        = null
      + project_id      = "ffffffff-AAAA-1414-eeee-000000000000"
      + sort_by         = null
      + sort_direction  = null
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```

### Apply the Plan

```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/laura.silaja/code/Lagher0/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Reading...
data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Read complete after 0s

Changes to Outputs:
  + project_backups_list = {
      + cursor          = {
          + hrefs = {
              + first    = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=10"
              + last     = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=10"
              + next     = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=0&perPage=10"
              + previous = ""
            }
          + pages = {
              + last        = 1
              + next        = 0
              + page        = 1
              + per_page    = 10
              + previous    = 0
              + total_items = 1
            }
        }
      + data            = [
          + {
              + cloud_provider       = "hostedAWS"
              + cluster_id           = "ffffffff-AAAA-1414-eeee-000000000000"
              + cluster_name         = "goldavinashkak"
              + created_by           = "testUser@couchbase.com"
              + creation_date_time   = "2025-10-29T14:57:09.838906045Z"
              + current_status       = "healthy"
              + most_recent_snapshot = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-29T15:01:04.912087543Z"
                  + cross_region_copies = []
                  + database_size       = 0
                  + expiration          = "2025-11-05T15:01:04.912087543Z"
                  + id                  = "ffffffff-aaaa-1414-EEEE-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-29T15:02:21.564367676Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 168
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + oldest_snapshot      = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-29T15:01:04.912087543Z"
                  + cross_region_copies = []
                  + database_size       = 0
                  + expiration          = "2025-11-05T15:01:04.912087543Z"
                  + id                  = "ffffffff-aaaa-1414-EEEE-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-29T15:02:21.564367676Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 168
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + region               = "us-east-1"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = null
      + per_page        = null
      + project_id      = "ffffffff-AAAA-1414-eeee-000000000000"
      + sort_by         = null
      + sort_direction  = null
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

project_backups_list = {
  "cursor" = {
    "hrefs" = {
      "first" = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=10"
      "last" = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=10"
      "next" = "http://localhost:8084/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=0&perPage=10"
      "previous" = ""
    }
    "pages" = {
      "last" = 1
      "next" = 0
      "page" = 1
      "per_page" = 10
      "previous" = 0
      "total_items" = 1
    }
  }
  "data" = tolist([
    {
      "cloud_provider" = "hostedAWS"
      "cluster_id" = "ffffffff-AAAA-1414-eeee-000000000000"
      "cluster_name" = "goldavinashkak"
      "created_by" = "testUser@couchbase.com"
      "creation_date_time" = "2025-10-29T14:57:09.838906045Z"
      "current_status" = "healthy"
      "most_recent_snapshot" = {
        "app_service" = ""
        "cluster_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "cmek" = toset([])
        "created_at" = "2025-10-29T15:01:04.912087543Z"
        "cross_region_copies" = toset([])
        "database_size" = 0
        "expiration" = "2025-11-05T15:01:04.912087543Z"
        "id" = "ffffffff-aaaa-1414-EEEE-000000000000"
        "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "progress" = {
          "status" = "complete"
          "time" = "2025-10-29T15:02:21.564367676Z"
        }
        "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "retention" = 168
        "server" = {
          "version" = "7.6.7"
        }
        "type" = "on_demand"
      }
      "oldest_snapshot" = {
        "app_service" = ""
        "cluster_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "cmek" = toset([])
        "created_at" = "2025-10-29T15:01:04.912087543Z"
        "cross_region_copies" = toset([])
        "database_size" = 0
        "expiration" = "2025-11-05T15:01:04.912087543Z"
        "id" = "ffffffff-aaaa-1414-EEEE-000000000000"
        "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "progress" = {
          "status" = "complete"
          "time" = "2025-10-29T15:02:21.564367676Z"
        }
        "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "retention" = 168
        "server" = {
          "version" = "7.6.7"
        }
        "type" = "on_demand"
      }
      "region" = "us-east-1"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "page" = tonumber(null)
  "per_page" = tonumber(null)
  "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
  "sort_by" = tostring(null)
  "sort_direction" = tostring(null)
}

```

