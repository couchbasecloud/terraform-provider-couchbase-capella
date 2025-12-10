# Capella Snapshot Backup Example

This example shows how to create and manage cloud snapshot backups in Capella.

This creates a new cloud snapshot backup of the selected Capella cluster. It uses the organization ID, project ID, and cluster ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new snapshot backup entry of an existing Capella cluster as stated in the `create_snapshot_backup.tf` file.
2. UPDATE: Edits the `retention` of the snapshot backup and/or restores the backup.
3. IMPORT: Import a snapshot backup that exists in Capella but not in the terraform state file.
4. DELETE: Delete the newly created snapshot backup from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
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

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup will be created
  + resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + progress            = (known after apply)
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + regions_to_copy     = [
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      + retention           = 144
      + server              = (known after apply)
      + size                = (known after apply)
      + type                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_cloud_snapshot_backup = {
      + cluster_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      + cmek                            = (known after apply)
      + created_at                      = (known after apply)
      + cross_region_copies             = (known after apply)
      + cross_region_restore_preference = null
      + expiration                      = (known after apply)
      + id                              = (known after apply)
      + organization_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + progress                        = (known after apply)
      + project_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      + regions_to_copy                 = [
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      + restore_times                   = null
      + retention                       = 144
      + server                          = (known after apply)
      + size                            = (known after apply)
      + type                            = (known after apply)
    }
```

### Apply the Plan, in order to create a new Snapshot Backup entry

Command: `terraform apply`

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

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup will be created
  + resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + progress            = (known after apply)
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + regions_to_copy     = [
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      + retention           = 144
      + server              = (known after apply)
      + size                = (known after apply)
      + type                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_cloud_snapshot_backup = {
      + cluster_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      + cmek                            = (known after apply)
      + created_at                      = (known after apply)
      + cross_region_copies             = (known after apply)
      + cross_region_restore_preference = null
      + expiration                      = (known after apply)
      + id                              = (known after apply)
      + organization_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + progress                        = (known after apply)
      + project_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      + regions_to_copy                 = [
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      + restore_times                   = null
      + retention                       = 144
      + server                          = (known after apply)
      + size                            = (known after apply)
      + type                            = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Creating...
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Creation complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_cloud_snapshot_backup = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cmek" = toset([])
  "created_at" = "2025-10-14T11:58:56.331819513Z"
  "cross_region_copies" = toset([])
  "cross_region_restore_preference" = tolist(null) /* of string */
  "expiration" = "2025-10-20T11:58:56.331819513Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "progress" = {
    "status" = "processing"
    "time" = "2025-10-14T11:58:56.374332305Z"
  }
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "regions_to_copy" = toset([
    "ap-southeast-1",
    "eu-west-1",
  ])
  "restore_times" = tonumber(null)
  "retention" = 144
  "server" = {
    "version" = ""
  }
  "size" = 0
  "type" = "on_demand"
}
```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup
```

## IMPORT
### Remove the resource `new_cloud_snapshot_backup` from the Terraform State file

Command: `terraform state rm couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup`

Sample Output:
``` 
Removed couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup id=<snapshot_backup_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Importing from ID "id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Import prepared!
  Prepared couchbase-capella_cloud_snapshot_backup for import
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Refreshing state... [id=id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the snapshot backup ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the snapshot backup belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit the snapshot backup retention in the terraform.tfvars file.
```
snapshot_backup = {
    retention = 168
    regions_to_copy = ["ap-southeast-1", "eu-west-1"]
    restore_times = 1
    cross_region_restore_preference = ["eu-west-1"]
}
```

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
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup has changed
  ~ resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      ~ cross_region_copies = [
          + {
              + region_code = "ap-southeast-1"
              + status      = "complete"
              + time        = "2025-10-14T12:00:47.10315555Z"
            },
          + {
              + region_code = "eu-west-1"
              + status      = "complete"
              + time        = "2025-10-14T12:01:20.18712051Z"
            },
        ]
        id                  = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ progress            = {
          ~ status = "processing" -> "complete"
          ~ time   = "2025-10-14T11:58:56.374332305Z" -> "2025-10-14T11:59:49.049001884Z"
        }
      ~ server              = {
          ~ version = "" -> "7.6.7"
        }
        # (10 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or
respond to these changes.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup will be updated in-place
  ~ resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      ~ cmek                            = [] -> (known after apply)
      ~ cross_region_copies             = [
          - {
              - region_code = "ap-southeast-1" -> null
              - status      = "complete" -> null
              - time        = "2025-10-14T12:00:47.10315555Z" -> null
            },
          - {
              - region_code = "eu-west-1" -> null
              - status      = "complete" -> null
              - time        = "2025-10-14T12:01:20.18712051Z" -> null
            },
        ] -> (known after apply)
      + cross_region_restore_preference = [
          + "eu-west-1",
        ]
      ~ expiration                      = "2025-10-20T11:58:56.331819513Z" -> (known after apply)
        id                              = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ progress                        = {
          ~ status = "complete" -> (known after apply)
          ~ time   = "2025-10-14T11:59:49.049001884Z" -> (known after apply)
        } -> (known after apply)
      + restore_times                   = 1
      ~ retention                       = 144 -> 168
      ~ server                          = {
          ~ version = "7.6.7" -> (known after apply)
        } -> (known after apply)
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cloud_snapshot_backup = {
      ~ cmek                            = [] -> (known after apply)
      ~ cross_region_copies             = [] -> (known after apply)
      ~ cross_region_restore_preference = null -> [
          + "eu-west-1",
        ]
      ~ expiration                      = "2025-10-20T11:58:56.331819513Z" -> (known after apply)
        id                              = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ progress                        = {
          - status = "processing"
          - time   = "2025-10-14T11:58:56.374332305Z"
        } -> (known after apply)
      ~ restore_times                   = null -> 1
      ~ retention                       = 144 -> 168
      ~ server                          = {
          - version = ""
        } -> (known after apply)
        # (7 unchanged attributes hidden)
    }
```

command: `terrafom apply`

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
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup has changed
  ~ resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      ~ cross_region_copies = [
          + {
              + region_code = "ap-southeast-1"
              + status      = "complete"
              + time        = "2025-10-14T12:00:47.10315555Z"
            },
          + {
              + region_code = "eu-west-1"
              + status      = "complete"
              + time        = "2025-10-14T12:01:20.18712051Z"
            },
        ]
        id                  = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ progress            = {
          ~ status = "processing" -> "complete"
          ~ time   = "2025-10-14T11:58:56.374332305Z" -> "2025-10-14T11:59:49.049001884Z"
        }
      ~ server              = {
          ~ version = "" -> "7.6.7"
        }
        # (10 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or
respond to these changes.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup will be updated in-place
  ~ resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      ~ cmek                            = [] -> (known after apply)
      ~ cross_region_copies             = [
          - {
              - region_code = "ap-southeast-1" -> null
              - status      = "complete" -> null
              - time        = "2025-10-14T12:00:47.10315555Z" -> null
            },
          - {
              - region_code = "eu-west-1" -> null
              - status      = "complete" -> null
              - time        = "2025-10-14T12:01:20.18712051Z" -> null
            },
        ] -> (known after apply)
      + cross_region_restore_preference = [
          + "eu-west-1",
        ]
      ~ expiration                      = "2025-10-20T11:58:56.331819513Z" -> (known after apply)
        id                              = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ progress                        = {
          ~ status = "complete" -> (known after apply)
          ~ time   = "2025-10-14T11:59:49.049001884Z" -> (known after apply)
        } -> (known after apply)
      + restore_times                   = 1
      ~ retention                       = 144 -> 168
      ~ server                          = {
          ~ version = "7.6.7" -> (known after apply)
        } -> (known after apply)
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cloud_snapshot_backup = {
      ~ cmek                            = [] -> (known after apply)
      ~ cross_region_copies             = [] -> (known after apply)
      ~ cross_region_restore_preference = null -> [
          + "eu-west-1",
        ]
      ~ expiration                      = "2025-10-20T11:58:56.331819513Z" -> (known after apply)
        id                              = "ffffffff-aaaa-1414-eeee-000000000000"
      ~ progress                        = {
          - status = "processing"
          - time   = "2025-10-14T11:58:56.374332305Z"
        } -> (known after apply)
      ~ restore_times                   = null -> 1
      ~ retention                       = 144 -> 168
      ~ server                          = {
          - version = ""
        } -> (known after apply)
        # (7 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Modifying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Modifications complete after 0s [id=ffffffff-aaaa-1414-eeee-000000000000]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_cloud_snapshot_backup = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cmek" = toset([])
  "created_at" = "2025-10-14T11:58:56.331819513Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-10-14T12:00:47.10315555Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-10-14T12:01:20.18712051Z"
    },
  ])
  "cross_region_restore_preference" = tolist([
    "eu-west-1",
  ])
  "expiration" = "2025-10-21T11:58:56.331819513Z"
  "id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "progress" = {
    "status" = "complete"
    "time" = "2025-10-14T11:59:49.049001884Z"
  }
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "regions_to_copy" = toset([
    "ap-southeast-1",
    "eu-west-1",
  ])
  "restore_times" = 1
  "retention" = 168
  "server" = {
    "version" = "7.6.7"
  }
  "size" = 0
  "type" = "on_demand"
}
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

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
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Refreshing state... [id=ffffffff-aaaa-1414-eeee-000000000000]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup will be destroyed
  - resource "couchbase-capella_cloud_snapshot_backup" "new_cloud_snapshot_backup" {
      - cluster_id                      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cmek                            = [] -> null
      - created_at                      = "2025-10-14T11:58:56.331819513Z" -> null
      - cross_region_copies             = [
          - {
              - region_code = "ap-southeast-1" -> null
              - status      = "complete" -> null
              - time        = "2025-10-14T12:00:47.10315555Z" -> null
            },
          - {
              - region_code = "eu-west-1" -> null
              - status      = "complete" -> null
              - time        = "2025-10-14T12:01:20.18712051Z" -> null
            },
        ] -> null
      - cross_region_restore_preference = [
          - "eu-west-1",
        ] -> null
      - expiration                      = "2025-10-21T11:58:56.331819513Z" -> null
      - id                              = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - progress                        = {
          - status = "complete" -> null
          - time   = "2025-10-14T11:59:49.049001884Z" -> null
        } -> null
      - project_id                      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - regions_to_copy                 = [
          - "ap-southeast-1",
          - "eu-west-1",
        ] -> null
      - restore_times                   = 1 -> null
      - retention                       = 168 -> null
      - server                          = {
          - version = "7.6.7" -> null
        } -> null
      - size                            = 0 -> null
      - type                            = "on_demand" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_cloud_snapshot_backup = {
      - cluster_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      - cmek                            = []
      - created_at                      = "2025-10-14T11:58:56.331819513Z"
      - cross_region_copies             = [
          - {
              - region_code = "ap-southeast-1"
              - status      = "complete"
              - time        = "2025-10-14T12:00:47.10315555Z"
            },
          - {
              - region_code = "eu-west-1"
              - status      = "complete"
              - time        = "2025-10-14T12:01:20.18712051Z"
            },
        ]
      - cross_region_restore_preference = [
          - "eu-west-1",
        ]
      - expiration                      = "2025-10-21T11:58:56.331819513Z"
      - id                              = "ffffffff-aaaa-1414-eeee-000000000000"
      - organization_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      - progress                        = {
          - status = "complete"
          - time   = "2025-10-14T11:59:49.049001884Z"
        }
      - project_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
      - regions_to_copy                 = [
          - "ap-southeast-1",
          - "eu-west-1",
        ]
      - restore_times                   = 1
      - retention                       = 168
      - server                          = {
          - version = "7.6.7"
        }
      - size                            = 0
      - type                            = "on_demand"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Destroying... [id=ffffffff-aaaa-1414-eeee-000000000000]
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
