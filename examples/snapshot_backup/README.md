# Capella Snapshot Backup Example

This example shows how to create and manage cloud snapshot backups in Capella.

This creates a new cloud snapshot backup of the selected Capella cluster. It uses the organization ID, project ID, and cluster ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new snapshot backup entry of an existing Capella cluster as stated in the `create_snapshot_backup.tf` file.
2. UPDATE: Edits the `retention` of the snapshot backup.
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

  # couchbase-capella_snapshot_backup.new_snapshot_backup will be created
  + resource "couchbase-capella_snapshot_backup" "new_snapshot_backup" {
      + app_service         = (known after apply)
      + cluster_id          = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress            = (known after apply)
      + project_id          = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
      + regions_to_copy     = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + retention           = 72
      + server              = (known after apply)
      + size                = (known after apply)
      + type                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_snapshot_backup = {
      + app_service         = (known after apply)
      + cluster_id          = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress            = (known after apply)
      + project_id          = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
      + regions_to_copy     = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + retention           = 72
      + server              = (known after apply)
      + size                = (known after apply)
      + type                = (known after apply)
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

  # couchbase-capella_snapshot_backup.new_snapshot_backup will be created
  + resource "couchbase-capella_snapshot_backup" "new_snapshot_backup" {
      + app_service         = (known after apply)
      + cluster_id          = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress            = (known after apply)
      + project_id          = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
      + regions_to_copy     = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + retention           = 72
      + server              = (known after apply)
      + size                = (known after apply)
      + type                = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_snapshot_backup = {
      + app_service         = (known after apply)
      + cluster_id          = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress            = (known after apply)
      + project_id          = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
      + regions_to_copy     = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + retention           = 72
      + server              = (known after apply)
      + size                = (known after apply)
      + type                = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_snapshot_backup.new_snapshot_backup: Creating...
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [00m10s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [00m20s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [00m30s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [00m40s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [00m50s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [01m00s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [01m10s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [01m20s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [01m30s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [01m40s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [01m50s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [02m00s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [02m10s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [02m20s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [02m30s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [02m40s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Still creating... [02m50s elapsed]
couchbase-capella_snapshot_backup.new_snapshot_backup: Creation complete after 2m53s [id=a52a2f82-538b-4999-b479-00d6a6f796e7]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_snapshot_backup = {
  "app_service" = ""
  "cluster_id" = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
  "cmek" = toset([])
  "created_at" = "2025-09-15T09:58:49.015571053Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-09-15T10:01:41.105485091Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-09-15T10:00:16.856243927Z"
    },
  ])
  "expiration" = "2025-09-18T09:58:49.015571053Z"
  "id" = "a52a2f82-538b-4999-b479-00d6a6f796e7"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "progress" = {
    "status" = "complete"
    "time" = "2025-09-15T09:59:42.945665925Z"
  }
  "project_id" = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
  "regions_to_copy" = tolist([
    "eu-west-1",
    "ap-southeast-1",
  ])
  "retention" = 72
  "server" = {
    "version" = "7.6.7"
  }
  "size" = 0
  "type" = "on_demand"
}
```

### Note the Snapshot Backup ID of the new Snapshot Backup
Command: `terraform output new_snapshot_backup`

Sample Output:
```
{
  "app_service" = ""
  "cluster_id" = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
  "cmek" = toset([])
  "created_at" = "2025-09-15T09:58:49.015571053Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-09-15T10:01:41.105485091Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-09-15T10:00:16.856243927Z"
    },
  ])
  "expiration" = "2025-09-18T09:58:49.015571053Z"
  "id" = "a52a2f82-538b-4999-b479-00d6a6f796e7"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "progress" = {
    "status" = "complete"
    "time" = "2025-09-15T09:59:42.945665925Z"
  }
  "project_id" = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
  "regions_to_copy" = tolist([
    "eu-west-1",
    "ap-southeast-1",
  ])
  "retention" = 72
  "server" = {
    "version" = "7.6.7"
  }
  "size" = 0
  "type" = "on_demand"
}
```
In this case, the snapshot backup ID for my new snapshot backup is `a52a2f82-538b-4999-b479-00d6a6f796e7`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_snapshot_backup.new_snapshot_backup
```

## IMPORT
### Remove the resource `new_snapshot_backup` from the Terraform State file

Command: `terraform state rm couchbase-capella_snapshot_backup.new_snapshot_backup`

Sample Output:
``` 
Removed couchbase-capella_snapshot_backup.new_snapshot_backup
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_snapshot_backup.new_snapshot_backup id=<snapshot_backup_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_snapshot_backup.new_snapshot_backup "id="a52a2f82-538b-4999-b479-00d6a6f796e7",cluster_id=98c6b0c0-3d95-4a74-a119-3ccef0aedc43,project_id=eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76"`

Sample Output:
``` 
couchbase-capella_snapshot_backup.new_snapshot_backup: Importing from ID "id=a52a2f82-538b-4999-b479-00d6a6f796e7,cluster_id=98c6b0c0-3d95-4a74-a119-3ccef0aedc43,project_id=eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76"...
couchbase-capella_snapshot_backup.new_snapshot_backup: Import prepared!
  Prepared couchbase-capella_snapshot_backup for import
couchbase-capella_snapshot_backup.new_snapshot_backup: Refreshing state... [id=id=a52a2f82-538b-4999-b479-00d6a6f796e7,cluster_id=98c6b0c0-3d95-4a74-a119-3ccef0aedc43,project_id=eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76]

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
    retention = 48
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
couchbase-capella_snapshot_backup.new_snapshot_backup: Refreshing state... [id=a52a2f82-538b-4999-b479-00d6a6f796e7]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup.new_snapshot_backup will be updated in-place
  ~ resource "couchbase-capella_snapshot_backup" "new_snapshot_backup" {
      ~ cmek                = [] -> (known after apply)
      ~ cross_region_copies = [
          - {
              - region_code = "ap-southeast-1" -> null
              - status      = "complete" -> null
              - time        = "2025-09-15T10:01:41.105485091Z" -> null
            },
          - {
              - region_code = "eu-west-1" -> null
              - status      = "complete" -> null
              - time        = "2025-09-15T10:00:16.856243927Z" -> null
            },
        ] -> (known after apply)
      ~ expiration          = "2025-09-18T09:58:49.015571053Z" -> (known after apply)
        id                  = "a52a2f82-538b-4999-b479-00d6a6f796e7"
      ~ progress            = {
          ~ status = "complete" -> (known after apply)
          ~ time   = "2025-09-15T09:59:42.945665925Z" -> (known after apply)
        } -> (known after apply)
      ~ retention           = 72 -> 48
      ~ server              = {
          ~ version = "7.6.7" -> (known after apply)
        } -> (known after apply)
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_snapshot_backup = {
      ~ cmek                = [] -> (known after apply)
      ~ cross_region_copies = [
          - {
              - region_code = "ap-southeast-1"
              - status      = "complete"
              - time        = "2025-09-15T10:01:41.105485091Z"
            },
          - {
              - region_code = "eu-west-1"
              - status      = "complete"
              - time        = "2025-09-15T10:00:16.856243927Z"
            },
        ] -> (known after apply)
      ~ expiration          = "2025-09-18T09:58:49.015571053Z" -> (known after apply)
        id                  = "a52a2f82-538b-4999-b479-00d6a6f796e7"
      ~ progress            = {
          - status = "complete"
          - time   = "2025-09-15T09:59:42.945665925Z"
        } -> (known after apply)
      ~ retention           = 72 -> 48
      ~ server              = {
          - version = "7.6.7"
        } -> (known after apply)
        # (8 unchanged attributes hidden)
    }
```

command: `terrafom apply`

Sample Output:

```
new_snapshot_backup = {
  "app_service" = ""
  "cluster_id" = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
  "cmek" = toset([])
  "created_at" = "2025-09-15T09:58:49.015571053Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-09-15T10:01:41.105485091Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-09-15T10:00:16.856243927Z"
    },
  ])
  "expiration" = "2025-09-17T09:58:49.015571053Z"
  "id" = "a52a2f82-538b-4999-b479-00d6a6f796e7"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "progress" = {
    "status" = "complete"
    "time" = "2025-09-15T09:59:42.945665925Z"
  }
  "project_id" = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
  "regions_to_copy" = tolist(null) /* of string */
  "retention" = 48
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
couchbase-capella_snapshot_backup.new_snapshot_backup: Refreshing state... [id=a52a2f82-538b-4999-b479-00d6a6f796e7]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup.new_snapshot_backup will be destroyed
  - resource "couchbase-capella_snapshot_backup" "new_snapshot_backup" {
      - app_service         = "" -> null
      - cluster_id          = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43" -> null
      - cmek                = [] -> null
      - created_at          = "2025-09-15T09:58:49.015571053Z" -> null
      - cross_region_copies = [
          - {
              - region_code = "ap-southeast-1" -> null
              - status      = "complete" -> null
              - time        = "2025-09-15T10:01:41.105485091Z" -> null
            },
          - {
              - region_code = "eu-west-1" -> null
              - status      = "complete" -> null
              - time        = "2025-09-15T10:00:16.856243927Z" -> null
            },
        ] -> null
      - expiration          = "2025-09-17T09:58:49.015571053Z" -> null
      - id                  = "a52a2f82-538b-4999-b479-00d6a6f796e7" -> null
      - organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76" -> null
      - progress            = {
          - status = "complete" -> null
          - time   = "2025-09-15T09:59:42.945665925Z" -> null
        } -> null
      - project_id          = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8" -> null
      - retention           = 48 -> null
      - server              = {
          - version = "7.6.7" -> null
        } -> null
      - size                = 0 -> null
      - type                = "on_demand" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_snapshot_backup = {
      - app_service         = ""
      - cluster_id          = "98c6b0c0-3d95-4a74-a119-3ccef0aedc43"
      - cmek                = []
      - created_at          = "2025-09-15T09:58:49.015571053Z"
      - cross_region_copies = [
          - {
              - region_code = "ap-southeast-1"
              - status      = "complete"
              - time        = "2025-09-15T10:01:41.105485091Z"
            },
          - {
              - region_code = "eu-west-1"
              - status      = "complete"
              - time        = "2025-09-15T10:00:16.856243927Z"
            },
        ]
      - expiration          = "2025-09-17T09:58:49.015571053Z"
      - id                  = "a52a2f82-538b-4999-b479-00d6a6f796e7"
      - organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      - progress            = {
          - status = "complete"
          - time   = "2025-09-15T09:59:42.945665925Z"
        }
      - project_id          = "eb16e9a7-6b5b-44ac-84c1-2a9a57b18cd8"
      - regions_to_copy     = null
      - retention           = 48
      - server              = {
          - version = "7.6.7"
        }
      - size                = 0
      - type                = "on_demand"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_snapshot_backup.new_snapshot_backup: Destroying... [id=a52a2f82-538b-4999-b479-00d6a6f796e7]
couchbase-capella_snapshot_backup.new_snapshot_backup: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
