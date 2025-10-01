# List Existing Cloud Snapshot Backups

This example shows how to list cloud snapshot backups that already exist in Capella for a given cluster. It uses the organization ID, project ID, and cluster ID to do so. 

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
data.couchbase-capella_cloud_snapshot_backups.existing_cloud_snapshot_backups: Reading...
data.couchbase-capella_cloud_snapshot_backups.existing_cloud_snapshot_backups: Read complete after 0s

Changes to Outputs:
  + backups_list = {
      + cluster_id      = "33207f51-4e73-448b-8634-1230a65c76e9"
      + data            = [
          + {
              + cmek                = []
              + created_at          = "2025-10-01T10:50:52.51832872Z"
              + cross_region_copies = [
                  + {
                      + region_code = "ap-southeast-1"
                      + status      = "deleting"
                      + time        = "2025-10-01T15:06:12.294832427Z"
                    },
                  + {
                      + region_code = "eu-west-1"
                      + status      = "deleting"
                      + time        = "2025-10-01T15:06:12.294832635Z"
                    },
                ]
              + expiration          = "2025-10-08T10:50:52.51832872Z"
              + id                  = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
              + progress            = {
                  + status = "deleting"
                  + time   = "2025-10-01T15:06:12.294831885Z"
                }
              + retention           = 168
              + server              = {
                  + version = "7.6.7"
                }
              + size                = 0
              + type                = "on_demand"
            },
          + {
              + cmek                = []
              + created_at          = "2025-10-01T10:53:06.540292046Z"
              + cross_region_copies = []
              + expiration          = "2025-10-08T10:53:06.540292046Z"
              + id                  = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
              + progress            = {
                  + status = "complete"
                  + time   = "2025-10-01T10:54:25.049507013Z"
                }
              + retention           = 168
              + server              = {
                  + version = "7.6.7"
                }
              + size                = 0
              + type                = "on_demand"
            },
          + {
              + cmek                = []
              + created_at          = "2025-10-01T14:22:31.580812042Z"
              + cross_region_copies = [
                  + {
                      + region_code = "ap-southeast-1"
                      + status      = "complete"
                      + time        = "2025-10-01T14:41:47.217926638Z"
                    },
                  + {
                      + region_code = "eu-west-1"
                      + status      = "complete"
                      + time        = "2025-10-01T14:40:21.843839959Z"
                    },
                ]
              + expiration          = "2025-10-08T14:22:31.580812042Z"
              + id                  = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
              + progress            = {
                  + status = "complete"
                  + time   = "2025-10-01T14:23:51.074934843Z"
                }
              + retention           = 168
              + server              = {
                  + version = "7.6.7"
                }
              + size                = 0
              + type                = "on_demand"
            },
        ]
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "a4f8ada6-32c5-4986-b3f1-1936b243e80b"
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
data.couchbase-capella_cloud_snapshot_backups.existing_cloud_snapshot_backups: Reading...
data.couchbase-capella_cloud_snapshot_backups.existing_cloud_snapshot_backups: Read complete after 0s

Changes to Outputs:
  + backups_list = {
      + cluster_id      = "33207f51-4e73-448b-8634-1230a65c76e9"
      + data            = [
          + {
              + cmek                = []
              + created_at          = "2025-10-01T10:50:52.51832872Z"
              + cross_region_copies = [
                  + {
                      + region_code = "ap-southeast-1"
                      + status      = "deleting"
                      + time        = "2025-10-01T15:06:12.294832427Z"
                    },
                  + {
                      + region_code = "eu-west-1"
                      + status      = "deleting"
                      + time        = "2025-10-01T15:06:12.294832635Z"
                    },
                ]
              + expiration          = "2025-10-08T10:50:52.51832872Z"
              + id                  = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
              + progress            = {
                  + status = "deleting"
                  + time   = "2025-10-01T15:06:12.294831885Z"
                }
              + retention           = 168
              + server              = {
                  + version = "7.6.7"
                }
              + size                = 0
              + type                = "on_demand"
            },
          + {
              + cmek                = []
              + created_at          = "2025-10-01T10:53:06.540292046Z"
              + cross_region_copies = []
              + expiration          = "2025-10-08T10:53:06.540292046Z"
              + id                  = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
              + progress            = {
                  + status = "complete"
                  + time   = "2025-10-01T10:54:25.049507013Z"
                }
              + retention           = 168
              + server              = {
                  + version = "7.6.7"
                }
              + size                = 0
              + type                = "on_demand"
            },
          + {
              + cmek                = []
              + created_at          = "2025-10-01T14:22:31.580812042Z"
              + cross_region_copies = [
                  + {
                      + region_code = "ap-southeast-1"
                      + status      = "complete"
                      + time        = "2025-10-01T14:41:47.217926638Z"
                    },
                  + {
                      + region_code = "eu-west-1"
                      + status      = "complete"
                      + time        = "2025-10-01T14:40:21.843839959Z"
                    },
                ]
              + expiration          = "2025-10-08T14:22:31.580812042Z"
              + id                  = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
              + progress            = {
                  + status = "complete"
                  + time   = "2025-10-01T14:23:51.074934843Z"
                }
              + retention           = 168
              + server              = {
                  + version = "7.6.7"
                }
              + size                = 0
              + type                = "on_demand"
            },
        ]
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "a4f8ada6-32c5-4986-b3f1-1936b243e80b"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

backups_list = {
  "cluster_id" = "33207f51-4e73-448b-8634-1230a65c76e9"
  "data" = tolist([
    {
      "cmek" = toset([])
      "created_at" = "2025-10-01T10:50:52.51832872Z"
      "cross_region_copies" = toset([
        {
          "region_code" = "ap-southeast-1"
          "status" = "deleting"
          "time" = "2025-10-01T15:06:12.294832427Z"
        },
        {
          "region_code" = "eu-west-1"
          "status" = "deleting"
          "time" = "2025-10-01T15:06:12.294832635Z"
        },
      ])
      "expiration" = "2025-10-08T10:50:52.51832872Z"
      "id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      "progress" = {
        "status" = "deleting"
        "time" = "2025-10-01T15:06:12.294831885Z"
      }
      "retention" = 168
      "server" = {
        "version" = "7.6.7"
      }
      "size" = 0
      "type" = "on_demand"
    },
    {
      "cmek" = toset([])
      "created_at" = "2025-10-01T10:53:06.540292046Z"
      "cross_region_copies" = toset([])
      "expiration" = "2025-10-08T10:53:06.540292046Z"
      "id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      "progress" = {
        "status" = "complete"
        "time" = "2025-10-01T10:54:25.049507013Z"
      }
      "retention" = 168
      "server" = {
        "version" = "7.6.7"
      }
      "size" = 0
      "type" = "on_demand"
    },
    {
      "cmek" = toset([])
      "created_at" = "2025-10-01T14:22:31.580812042Z"
      "cross_region_copies" = toset([
        {
          "region_code" = "ap-southeast-1"
          "status" = "complete"
          "time" = "2025-10-01T14:41:47.217926638Z"
        },
        {
          "region_code" = "eu-west-1"
          "status" = "complete"
          "time" = "2025-10-01T14:40:21.843839959Z"
        },
      ])
      "expiration" = "2025-10-08T14:22:31.580812042Z"
      "id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      "progress" = {
        "status" = "complete"
        "time" = "2025-10-01T14:23:51.074934843Z"
      }
      "retention" = 168
      "server" = {
        "version" = "7.6.7"
      }
      "size" = 0
      "type" = "on_demand"
    },
  ])
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "a4f8ada6-32c5-4986-b3f1-1936b243e80b"
}
```

