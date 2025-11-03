# List Existing Project Level Cloud Snapshot Backups

This example shows how to list cloud snapshot backups that already exist in Capella for a given cluster. It uses the organization ID and project ID to do so. 

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## List
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```

 terraform plan

data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Reading...
data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Read complete after 0s

Changes to Outputs:
  + project_backups_list = {
      + cursor          = {
          + hrefs = {
              + first    = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=2"
              + last     = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=2"
              + next     = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=0&perPage=2"
              + previous = ""
            }
          + pages = {
              + last        = 1
              + next        = 0
              + page        = 1
              + per_page    = 2
              + previous    = 0
              + total_items = 2
            }
        }
      + data            = [
          + {
              + cloud_provider       = "hostedAWS"
              + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name         = "mintgordonmoore"
              + created_by           = "user@ccouchbase.com"
              + creation_date_time   = "2025-10-30T10:38:00.991148385Z"
              + current_status       = "healthy"
              + most_recent_snapshot = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T11:02:16.291579761Z"
                  + cross_region_copies = []
                  + database_size       = 78068736
                  + expiration          = "2025-11-01T11:02:16.291579761Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T11:03:06.67205009Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 48
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + oldest_snapshot      = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T10:58:29.994815086Z"
                  + cross_region_copies = [
                      + {
                          + region_code = "ap-south-2"
                          + status      = "complete"
                          + time        = "2025-10-30T11:02:27.723455335Z"
                        },
                    ]
                  + database_size       = 78068736
                  + expiration          = "2025-11-06T10:58:29.994815086Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T10:59:20.63490193Z"
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
          + {
              + cloud_provider       = "hostedAWS"
              + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name         = "cheerfuldavidwagner"
              + created_by           = "user@ccouchbase.com"
              + creation_date_time   = "2025-10-30T10:38:13.697622752Z"
              + current_status       = "healthy"
              + most_recent_snapshot = {
                  + app_service         = "3.3.0"
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T10:58:07.582313965Z"
                  + cross_region_copies = [
                      + {
                          + region_code = "ap-south-1"
                          + status      = "complete"
                          + time        = "2025-10-30T11:00:17.706695303Z"
                        },
                    ]
                  + database_size       = 24709758
                  + expiration          = "2025-11-07T10:58:07.582313965Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T10:58:58.635103461Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 192
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + oldest_snapshot      = {
                  + app_service         = "3.3.0"
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T10:58:07.582313965Z"
                  + cross_region_copies = [
                      + {
                          + region_code = "ap-south-1"
                          + status      = "complete"
                          + time        = "2025-10-30T11:00:17.706695303Z"
                        },
                    ]
                  + database_size       = 24709758
                  + expiration          = "2025-11-07T10:58:07.582313965Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T10:58:58.635103461Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 192
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + region               = "us-east-1"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = 1
      + per_page        = 2
      + project_id      = "ffffffff-AAAA-1414-eeee-000000000000"
      + sort_by         = "region"
      + sort_direction  = "asc"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

```

### Apply the Plan

```
terraform apply

data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Reading...
data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups: Read complete after 0s

Changes to Outputs:
  + project_backups_list = {
      + cursor          = {
          + hrefs = {
              + first    = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=2"
              + last     = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=2"
              + next     = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=0&perPage=2"
              + previous = ""
            }
          + pages = {
              + last        = 1
              + next        = 0
              + page        = 1
              + per_page    = 2
              + previous    = 0
              + total_items = 2
            }
        }
      + data            = [
          + {
              + cloud_provider       = "hostedAWS"
              + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name         = "mintgordonmoore"
              + created_by           = "user@ccouchbase.com"
              + creation_date_time   = "2025-10-30T10:38:00.991148385Z"
              + current_status       = "healthy"
              + most_recent_snapshot = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T11:02:16.291579761Z"
                  + cross_region_copies = []
                  + database_size       = 78068736
                  + expiration          = "2025-11-01T11:02:16.291579761Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T11:03:06.67205009Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 48
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + oldest_snapshot      = {
                  + app_service         = ""
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T10:58:29.994815086Z"
                  + cross_region_copies = [
                      + {
                          + region_code = "ap-south-2"
                          + status      = "complete"
                          + time        = "2025-10-30T11:02:27.723455335Z"
                        },
                    ]
                  + database_size       = 78068736
                  + expiration          = "2025-11-06T10:58:29.994815086Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T10:59:20.63490193Z"
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
          + {
              + cloud_provider       = "hostedAWS"
              + cluster_id           = "ffffffff-aaaa-1414-eeee-000000000000"
              + cluster_name         = "cheerfuldavidwagner"
              + created_by           = "user@couchbase.com"
              + creation_date_time   = "2025-10-30T10:38:13.697622752Z"
              + current_status       = "healthy"
              + most_recent_snapshot = {
                  + app_service         = "3.3.0"
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T10:58:07.582313965Z"
                  + cross_region_copies = [
                      + {
                          + region_code = "ap-south-1"
                          + status      = "complete"
                          + time        = "2025-10-30T11:00:17.706695303Z"
                        },
                    ]
                  + database_size       = 24709758
                  + expiration          = "2025-11-07T10:58:07.582313965Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T10:58:58.635103461Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 192
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + oldest_snapshot      = {
                  + app_service         = "3.3.0"
                  + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
                  + cmek                = []
                  + created_at          = "2025-10-30T10:58:07.582313965Z"
                  + cross_region_copies = [
                      + {
                          + region_code = "ap-south-1"
                          + status      = "complete"
                          + time        = "2025-10-30T11:00:17.706695303Z"
                        },
                    ]
                  + database_size       = 24709758
                  + expiration          = "2025-11-07T10:58:07.582313965Z"
                  + id                  = "ffffffff-aaaa-1414-eeee-000000000000"
                  + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
                  + progress            = {
                      + status = "complete"
                      + time   = "2025-10-30T10:58:58.635103461Z"
                    }
                  + project_id          = "ffffffff-AAAA-1414-eeee-000000000000"
                  + retention           = 192
                  + server              = {
                      + version = "7.6.7"
                    }
                  + type                = "on_demand"
                }
              + region               = "us-east-1"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + page            = 1
      + per_page        = 2
      + project_id      = "ffffffff-AAAA-1414-eeee-000000000000"
      + sort_by         = "region"
      + sort_direction  = "asc"
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
      "first" = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=2"
      "last" = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=1&perPage=2"
      "next" = "https://cloudapi.cloud.couchbase.com/v4/organizations/ffffffff-aaaa-1414-eeee-000000000000/projects/ffffffff-AAAA-1414-eeee-000000000000/cloudsnapshotbackups?page=0&perPage=2"
      "previous" = ""
    }
    "pages" = {
      "last" = 1
      "next" = 0
      "page" = 1
      "per_page" = 2
      "previous" = 0
      "total_items" = 2
    }
  }
  "data" = tolist([
    {
      "cloud_provider" = "hostedAWS"
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "mintgordonmoore"
      "created_by" = "user@ccouchbase.com"
      "creation_date_time" = "2025-10-30T10:38:00.991148385Z"
      "current_status" = "healthy"
      "most_recent_snapshot" = {
        "app_service" = ""
        "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "cmek" = toset([])
        "created_at" = "2025-10-30T11:02:16.291579761Z"
        "cross_region_copies" = toset([])
        "database_size" = 78068736
        "expiration" = "2025-11-01T11:02:16.291579761Z"
        "id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "progress" = {
          "status" = "complete"
          "time" = "2025-10-30T11:03:06.67205009Z"
        }
        "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "retention" = 48
        "server" = {
          "version" = "7.6.7"
        }
        "type" = "on_demand"
      }
      "oldest_snapshot" = {
        "app_service" = ""
        "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "cmek" = toset([])
        "created_at" = "2025-10-30T10:58:29.994815086Z"
        "cross_region_copies" = toset([
          {
            "region_code" = "ap-south-2"
            "status" = "complete"
            "time" = "2025-10-30T11:02:27.723455335Z"
          },
        ])
        "database_size" = 78068736
        "expiration" = "2025-11-06T10:58:29.994815086Z"
        "id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "progress" = {
          "status" = "complete"
          "time" = "2025-10-30T10:59:20.63490193Z"
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
    {
      "cloud_provider" = "hostedAWS"
      "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
      "cluster_name" = "cheerfuldavidwagner"
      "created_by" = "user@ccouchbase.com"
      "creation_date_time" = "2025-10-30T10:38:13.697622752Z"
      "current_status" = "healthy"
      "most_recent_snapshot" = {
        "app_service" = "3.3.0"
        "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "cmek" = toset([])
        "created_at" = "2025-10-30T10:58:07.582313965Z"
        "cross_region_copies" = toset([
          {
            "region_code" = "ap-south-1"
            "status" = "complete"
            "time" = "2025-10-30T11:00:17.706695303Z"
          },
        ])
        "database_size" = 24709758
        "expiration" = "2025-11-07T10:58:07.582313965Z"
        "id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "progress" = {
          "status" = "complete"
          "time" = "2025-10-30T10:58:58.635103461Z"
        }
        "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "retention" = 192
        "server" = {
          "version" = "7.6.7"
        }
        "type" = "on_demand"
      }
      "oldest_snapshot" = {
        "app_service" = "3.3.0"
        "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "cmek" = toset([])
        "created_at" = "2025-10-30T10:58:07.582313965Z"
        "cross_region_copies" = toset([
          {
            "region_code" = "ap-south-1"
            "status" = "complete"
            "time" = "2025-10-30T11:00:17.706695303Z"
          },
        ])
        "database_size" = 24709758
        "expiration" = "2025-11-07T10:58:07.582313965Z"
        "id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
        "progress" = {
          "status" = "complete"
          "time" = "2025-10-30T10:58:58.635103461Z"
        }
        "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
        "retention" = 192
        "server" = {
          "version" = "7.6.7"
        }
        "type" = "on_demand"
      }
      "region" = "us-east-1"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "page" = 1
  "per_page" = 2
  "project_id" = "ffffffff-AAAA-1414-eeee-000000000000"
  "sort_by" = "region"
  "sort_direction" = "asc"
}

```

