# Capella Backup Example

This example shows how to create and manage Backups in Capella.

This creates a new backup in the selected Capella cluster. It uses the organization ID, project ID, cluster ID and bucket ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new backup entry in an existing Capella cluster as stated in the `create_backup.tf` file.
2. UPDATE: Triggers restore for the backup.
3. LIST: List existing backups in Capella as stated in the `list_backups.tf` file.
4. IMPORT: Import a backup that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created backup from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan    
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_backups.existing_backups: Reading...
data.capella_backups.existing_backups: Still reading... [10s elapsed]
data.capella_backups.existing_backups: Read complete after 13s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_backup.new_backup will be created
  + resource "capella_backup" "new_backup" {
      + backup_stats            = (known after apply)
      + bucket_id               = "YjE="
      + bucket_name             = (known after apply)
      + cloud_provider          = (known after apply)
      + cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
      + cycle_id                = (known after apply)
      + date                    = (known after apply)
      + elapsed_time_in_seconds = (known after apply)
      + id                      = (known after apply)
      + method                  = (known after apply)
      + organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + restore_before          = (known after apply)
      + schedule_info           = (known after apply)
      + source                  = (known after apply)
      + status                  = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + backups_list = {
      + bucket_id       = "YjE="
      + cluster_id      = "1f6bad22-602f-407b-a567-7a8f672db227"
      + data            = [
          + {
              + backup_stats            = {
                  + cbas       = 0
                  + event      = 0
                  + fts        = 0
                  + gsi        = 0
                  + items      = 0
                  + mutations  = 0
                  + size_in_mb = 0.000527
                  + tombstones = 0
                }
              + bucket_id               = "YjE="
              + bucket_name             = "test-bucket"
              + cloud_provider          = "hostedAWS"
              + cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
              + cycle_id                = "635c196c-f91c-4c30-a33e-66fd1fa86b51"
              + date                    = "2023-11-13T17:37:26.24907334Z"
              + elapsed_time_in_seconds = 6
              + id                      = "0a04a68d-7a05-4189-9d2b-9ae4b5e3e230"
              + method                  = "full"
              + organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
              + restore_before          = "2023-12-14T00:00:00Z"
              + schedule_info           = {
                  + backup_time = "2023-11-13 17:37:26.24907334 +0000 UTC"
                  + backup_type = "Manual"
                  + increment   = 1
                  + retention   = "30days"
                }
              + source                  = "manual"
              + status                  = "ready"
            },
        ]
      + organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
    }
  + new_backup   = {
      + backup_stats            = (known after apply)
      + bucket_id               = "YjE="
      + bucket_name             = (known after apply)
      + cloud_provider          = (known after apply)
      + cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
      + cycle_id                = (known after apply)
      + date                    = (known after apply)
      + elapsed_time_in_seconds = (known after apply)
      + id                      = (known after apply)
      + method                  = (known after apply)
      + organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + restore_before          = (known after apply)
      + schedule_info           = (known after apply)
      + source                  = (known after apply)
      + status                  = (known after apply)
    }

```

### Apply the Plan, in order to create a new Backup entry

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_backups.existing_backups: Reading...
data.capella_backups.existing_backups: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_backup.new_backup will be created
  + resource "capella_backup" "new_backup" {
      + backup_stats            = (known after apply)
      + bucket_id               = "YjE="
      + bucket_name             = (known after apply)
      + cloud_provider          = (known after apply)
      + cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
      + cycle_id                = (known after apply)
      + date                    = (known after apply)
      + elapsed_time_in_seconds = (known after apply)
      + id                      = (known after apply)
      + method                  = (known after apply)
      + organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + restore_before          = (known after apply)
      + schedule_info           = (known after apply)
      + source                  = (known after apply)
      + status                  = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  ~ backups_list = {
      ~ data            = [
          + {
              + backup_stats            = {
                  + cbas       = 0
                  + event      = 0
                  + fts        = 0
                  + gsi        = 0
                  + items      = 0
                  + mutations  = 0
                  + size_in_mb = 0.000527
                  + tombstones = 0
                }
              + bucket_id               = "YjE="
              + bucket_name             = "test-bucket"
              + cloud_provider          = "hostedAWS"
              + cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
              + cycle_id                = "99c971b6-451c-480a-93cf-5313ff13005d"
              + date                    = "2023-11-13T19:05:33.760979469Z"
              + elapsed_time_in_seconds = 6
              + id                      = "92ca9452-e5cb-4c16-923b-985858448e09"
              + method                  = "full"
              + organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
              + project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
              + restore_before          = "2023-12-14T00:00:00Z"
              + schedule_info           = {
                  + backup_time = "2023-11-13 19:05:33.760979469 +0000 UTC"
                  + backup_type = "Manual"
                  + increment   = 1
                  + retention   = "30days"
                }
              + source                  = "manual"
              + status                  = "ready"
            },
            {
                backup_stats            = {
                    cbas       = 0
                    event      = 0
                    fts        = 0
                    gsi        = 0
                    items      = 0
                    mutations  = 0
                    size_in_mb = 0.000527
                    tombstones = 0
                }
                bucket_id               = "YjE="
                bucket_name             = "test-bucket"
                cloud_provider          = "hostedAWS"
                cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
                cycle_id                = "635c196c-f91c-4c30-a33e-66fd1fa86b51"
                date                    = "2023-11-13T17:37:26.24907334Z"
                elapsed_time_in_seconds = 6
                id                      = "0a04a68d-7a05-4189-9d2b-9ae4b5e3e230"
                method                  = "full"
                organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
                project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
                restore_before          = "2023-12-14T00:00:00Z"
                schedule_info           = {
                    backup_time = "2023-11-13 17:37:26.24907334 +0000 UTC"
                    backup_type = "Manual"
                    increment   = 1
                    retention   = "30days"
                }
                source                  = "manual"
                status                  = "ready"
            },
        ]
        # (4 unchanged attributes hidden)
    }
  + new_backup   = {
      + backup_stats            = (known after apply)
      + bucket_id               = "YjE="
      + bucket_name             = (known after apply)
      + cloud_provider          = (known after apply)
      + cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
      + cycle_id                = (known after apply)
      + date                    = (known after apply)
      + elapsed_time_in_seconds = (known after apply)
      + id                      = (known after apply)
      + method                  = (known after apply)
      + organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
      + project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + restore_before          = (known after apply)
      + schedule_info           = (known after apply)
      + source                  = (known after apply)
      + status                  = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_backup.new_backup: Creating...
capella_backup.new_backup: Still creating... [10s elapsed]
capella_backup.new_backup: Still creating... [20s elapsed]
...
...
capella_backup.new_backup: Still creating... [2m0s elapsed]
capella_backup.new_backup: Creation complete after 2m6s [id=58dd0f30-323b-461c-83a8-1d2719f4bcee]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

backups_list = {
  "bucket_id" = "YjE="
  "cluster_id" = "1f6bad22-602f-407b-a567-7a8f672db227"
  "data" = tolist([
    {
      "backup_stats" = {
        "cbas" = 0
        "event" = 0
        "fts" = 0
        "gsi" = 0
        "items" = 0
        "mutations" = 0
        "size_in_mb" = 0.000527
        "tombstones" = 0
      }
      "bucket_id" = "YjE="
      "bucket_name" = "test-bucket"
      "cloud_provider" = "hostedAWS"
      "cluster_id" = "1f6bad22-602f-407b-a567-7a8f672db227"
      "cycle_id" = "99c971b6-451c-480a-93cf-5313ff13005d"
      "date" = "2023-11-13T19:05:33.760979469Z"
      "elapsed_time_in_seconds" = 6
      "id" = "92ca9452-e5cb-4c16-923b-985858448e09"
      "method" = "full"
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      "restore_before" = "2023-12-14T00:00:00Z"
      "schedule_info" = {
        "backup_time" = "2023-11-13 19:05:33.760979469 +0000 UTC"
        "backup_type" = "Manual"
        "increment" = 1
        "retention" = "30days"
      }
      "source" = "manual"
      "status" = "ready"
    },
    {
      "backup_stats" = {
        "cbas" = 0
        "event" = 0
        "fts" = 0
        "gsi" = 0
        "items" = 0
        "mutations" = 0
        "size_in_mb" = 0.000527
        "tombstones" = 0
      }
      "bucket_id" = "YjE="
      "bucket_name" = "test-bucket"
      "cloud_provider" = "hostedAWS"
      "cluster_id" = "1f6bad22-602f-407b-a567-7a8f672db227"
      "cycle_id" = "635c196c-f91c-4c30-a33e-66fd1fa86b51"
      "date" = "2023-11-13T17:37:26.24907334Z"
      "elapsed_time_in_seconds" = 6
      "id" = "0a04a68d-7a05-4189-9d2b-9ae4b5e3e230"
      "method" = "full"
      "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
      "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      "restore_before" = "2023-12-14T00:00:00Z"
      "schedule_info" = {
        "backup_time" = "2023-11-13 17:37:26.24907334 +0000 UTC"
        "backup_type" = "Manual"
        "increment" = 1
        "retention" = "30days"
      }
      "source" = "manual"
      "status" = "ready"
    },
  ])
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
}
new_backup = {
  "backup_stats" = {
    "cbas" = 0
    "event" = 0
    "fts" = 0
    "gsi" = 0
    "items" = 0
    "mutations" = 0
    "size_in_mb" = 0.000527
    "tombstones" = 0
  }
  "bucket_id" = "YjE="
  "bucket_name" = "test-bucket"
  "cloud_provider" = "hostedAWS"
  "cluster_id" = "1f6bad22-602f-407b-a567-7a8f672db227"
  "cycle_id" = "f37575d4-6531-4732-9ad8-734f1831e32e"
  "date" = "2023-11-13T19:15:00.152667728Z"
  "elapsed_time_in_seconds" = 8
  "id" = "58dd0f30-323b-461c-83a8-1d2719f4bcee"
  "method" = "full"
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "restore_before" = "2023-12-14T00:00:00Z"
  "schedule_info" = {
    "backup_time" = "2023-11-13 19:15:00.152667728 +0000 UTC"
    "backup_type" = "Manual"
    "increment" = 1
    "retention" = "30days"
  }
  "source" = "manual"
  "status" = "ready"
}

```

### Note the Backup ID of the new Backup
Command: `terraform output new_backup`

Sample Output:
```
{
  "backup_stats" = {
    "cbas" = 0
    "event" = 0
    "fts" = 0
    "gsi" = 0
    "items" = 0
    "mutations" = 0
    "size_in_mb" = 0.000527
    "tombstones" = 0
  }
  "bucket_id" = "YjE="
  "bucket_name" = "test-bucket"
  "cloud_provider" = "hostedAWS"
  "cluster_id" = "1f6bad22-602f-407b-a567-7a8f672db227"
  "cycle_id" = "f37575d4-6531-4732-9ad8-734f1831e32e"
  "date" = "2023-11-13T19:15:00.152667728Z"
  "elapsed_time_in_seconds" = 8
  "id" = "58dd0f30-323b-461c-83a8-1d2719f4bcee"
  "method" = "full"
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "restore_before" = "2023-12-14T00:00:00Z"
  "schedule_info" = {
    "backup_time" = "2023-11-13 19:15:00.152667728 +0000 UTC"
    "backup_type" = "Manual"
    "increment" = 1
    "retention" = "30days"
  }
  "source" = "manual"
  "status" = "ready"
}

```
In this case, the backup ID for my new backup is `58dd0f30-323b-461c-83a8-1d2719f4bcee`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
terraform state list                                  
data.couchbase-capella_backups.existing_backups
couchbase-capella_backup.new_backup
```

## IMPORT
### Remove the resource `new_backup` from the Terraform State file

Command: `terraform state rm couchbase-capella_backup.new_backup`

Sample Output:
``` 
terraform state rm couchbase-capella_backup.new_backup
Removed capella_backup.new_backup
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_backup.new_backup id=<backup_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_backup.new_backup id=58dd0f30-323b-461c-83a8-1d2719f4bcee,cluster_id=1f6bad22-602f-407b-a567-7a8f672db227,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d`

Sample Output:
``` 
terraform import couchbase-capella_backup.new_backup id=58dd0f30-323b-461c-83a8-1d2719f4bcee,cluster_id=1f6bad22-602f-407b-a567-7a8f672db227,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d
capella_backup.new_backup: Importing from ID "id=58dd0f30-323b-461c-83a8-1d2719f4bcee,cluster_id=1f6bad22-602f-407b-a567-7a8f672db227,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d"...
data.capella_backups.existing_backups: Reading...
capella_backup.new_backup: Import prepared!
  Prepared capella_backup for import
capella_backup.new_backup: Refreshing state... [id=id=58dd0f30-323b-461c-83a8-1d2719f4bcee,cluster_id=1f6bad22-602f-407b-a567-7a8f672db227,project_id=f14134f2-7943-4e7b-b2c5-fc2071728b6e,organization_id=c2e9ccf6-4293-4635-9205-1204d074447d]
data.capella_backups.existing_backups: Read complete after 2s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the backup ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which backup belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit terraform.tfvars file to restore the backup.
```
resource "capella_backup" "new_backup" {
  organization_id            = var.organization_id
  project_id                 = var.project_id
  cluster_id                 = var.cluster_id
  bucket_id                  = var.bucket_id
  restore = {
    target_cluster_id = var.cluster_id
    source_cluster_id = var.cluster_id
    "services": [
      "data",
      "query"
    ],
  }
  restore_times = 1
}
```

Command: `terraform plan`

Sample Output:

``` 
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
capella_backup.new_backup: Refreshing state... [id=b1e946dc-e72b-4547-97e2-4337eabf06af]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_backup.new_backup will be updated in-place
  ~ resource "capella_backup" "new_backup" {
        id                      = "b1e946dc-e72b-4547-97e2-4337eabf06af"
      + restore                 = {
          + services          = [
              + "data",
              + "query",
            ]
          + source_cluster_id = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
          + status            = (known after apply)
          + target_cluster_id = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
        }
      + restore_times           = 1
        # (15 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_backup = {
        id                      = "b1e946dc-e72b-4547-97e2-4337eabf06af"
      ~ restore                 = null -> {
          + auto_remove_collections = null
          + exclude_data            = null
          + filter_keys             = null
          + filter_values           = null
          + force_updates           = null
          + include_data            = null
          + map_data                = null
          + replace_ttl             = null
          + replace_ttl_with        = null
          + services                = [
              + "data",
              + "query",
            ]
          + source_cluster_id       = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
          + target_cluster_id       = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
        }
      ~ restore_times           = null -> 1
        # (15 unchanged attributes hidden)
    }
```

command: `terrafom apply`

Sample Output:

```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
capella_backup.new_backup: Refreshing state... [id=b1e946dc-e72b-4547-97e2-4337eabf06af]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_backup.new_backup will be updated in-place
  ~ resource "capella_backup" "new_backup" {
        id                      = "b1e946dc-e72b-4547-97e2-4337eabf06af"
      + restore                 = {
          + services          = [
              + "data",
              + "query",
            ]
          + source_cluster_id = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
          + status            = (known after apply)
          + target_cluster_id = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
        }
      + restore_times           = 1
        # (15 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_backup = {
        id                      = "b1e946dc-e72b-4547-97e2-4337eabf06af"
      ~ restore                 = null -> {
          + auto_remove_collections = null
          + exclude_data            = null
          + filter_keys             = null
          + filter_values           = null
          + force_updates           = null
          + include_data            = null
          + map_data                = null
          + replace_ttl             = null
          + replace_ttl_with        = null
          + services                = [
              + "data",
              + "query",
            ]
          + source_cluster_id       = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
          + target_cluster_id       = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
        }
      ~ restore_times           = null -> 1
        # (15 unchanged attributes hidden)
    }
capella_backup.new_backup: Modifying... [id=b1e946dc-e72b-4547-97e2-4337eabf06af]
capella_backup.new_backup: Modifications complete after 1s [id=b1e946dc-e72b-4547-97e2-4337eabf06af]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_backup = {
  "backup_stats" = {
    "cbas" = 0
    "event" = 0
    "fts" = 0
    "gsi" = 0
    "items" = 0
    "mutations" = 0
    "size_in_mb" = 0.000527
    "tombstones" = 0
  }
  "bucket_id" = "Q0JFeGFtcGxlMQ=="
  "bucket_name" = "CBExample1"
  "cloud_provider" = "hostedAWS"
  "cluster_id" = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
  "cycle_id" = "36735b2b-f2c1-4a8b-af5b-992b0a326eab"
  "date" = "2023-11-15T17:26:30.816620426Z"
  "elapsed_time_in_seconds" = 6
  "id" = "b1e946dc-e72b-4547-97e2-4337eabf06af"
  "method" = "full"
  "organization_id" = "c2e9ccf6-4293-4635-9205-1204d074447d"
  "project_id" = "ad1b3554-0fc6-45f6-aec9-f994b4ad0729"
  "restore" = {
    "auto_remove_collections" = tobool(null)
    "exclude_data" = tostring(null)
    "filter_keys" = tostring(null)
    "filter_values" = tostring(null)
    "force_updates" = tobool(null)
    "include_data" = tostring(null)
    "map_data" = tostring(null)
    "replace_ttl" = tostring(null)
    "replace_ttl_with" = tostring(null)
    "services" = tolist([
      "data",
      "query",
    ])
    "source_cluster_id" = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
    "status" = "RESTORE INITIATED"
    "target_cluster_id" = "98964fdd-c45f-448c-8cbc-3c77e2e700a5"
  }
  "restore_before" = "2023-12-16T00:00:00Z"
  "restore_times" = 1
  "schedule_info" = {
    "backup_time" = "2023-11-15 17:26:30.816620426 +0000 UTC"
    "backup_type" = "Manual"
    "increment" = 1
    "retention" = "30days"
  }
  "source" = "manual"
  "status" = "ready"
}


```
Note:
```
The 'restore_times' field is incremental in nature. Therefore, whenever we need to trigger a restore, the value should be greater than the previous value.
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
data.capella_backups.existing_backups: Reading...
capella_backup.new_backup: Refreshing state... [id=58dd0f30-323b-461c-83a8-1d2719f4bcee]
data.capella_backups.existing_backups: Read complete after 2s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_backup.new_backup will be destroyed
  - resource "capella_backup" "new_backup" {
      - backup_stats            = {
          - cbas       = 0 -> null
          - event      = 0 -> null
          - fts        = 0 -> null
          - gsi        = 0 -> null
          - items      = 0 -> null
          - mutations  = 0 -> null
          - size_in_mb = 0.000527 -> null
          - tombstones = 0 -> null
        } -> null
      - bucket_id               = "YjE=" -> null
      - bucket_name             = "test-bucket" -> null
      - cloud_provider          = "hostedAWS" -> null
      - cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227" -> null
      - cycle_id                = "f37575d4-6531-4732-9ad8-734f1831e32e" -> null
      - date                    = "2023-11-13T19:15:00.152667728Z" -> null
      - elapsed_time_in_seconds = 8 -> null
      - id                      = "58dd0f30-323b-461c-83a8-1d2719f4bcee" -> null
      - method                  = "full" -> null
      - organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d" -> null
      - project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e" -> null
      - restore_before          = "2023-12-14T00:00:00Z" -> null
      - schedule_info           = {
          - backup_time = "2023-11-13 19:15:00.152667728 +0000 UTC" -> null
          - backup_type = "Manual" -> null
          - increment   = 1 -> null
          - retention   = "30days" -> null
        } -> null
      - source                  = "manual" -> null
      - status                  = "ready" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - backups_list = {
      - bucket_id       = "YjE="
      - cluster_id      = "1f6bad22-602f-407b-a567-7a8f672db227"
      - data            = [
          - {
              - backup_stats            = {
                  - cbas       = 0
                  - event      = 0
                  - fts        = 0
                  - gsi        = 0
                  - items      = 0
                  - mutations  = 0
                  - size_in_mb = 0.000527
                  - tombstones = 0
                }
              - bucket_id               = "YjE="
              - bucket_name             = "test-bucket"
              - cloud_provider          = "hostedAWS"
              - cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
              - cycle_id                = "f37575d4-6531-4732-9ad8-734f1831e32e"
              - date                    = "2023-11-13T19:15:00.152667728Z"
              - elapsed_time_in_seconds = 8
              - id                      = "58dd0f30-323b-461c-83a8-1d2719f4bcee"
              - method                  = "full"
              - organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
              - restore_before          = "2023-12-14T00:00:00Z"
              - schedule_info           = {
                  - backup_time = "2023-11-13 19:15:00.152667728 +0000 UTC"
                  - backup_type = "Manual"
                  - increment   = 1
                  - retention   = "30days"
                }
              - source                  = "manual"
              - status                  = "ready"
            },
          - {
              - backup_stats            = {
                  - cbas       = 0
                  - event      = 0
                  - fts        = 0
                  - gsi        = 0
                  - items      = 0
                  - mutations  = 0
                  - size_in_mb = 0.000527
                  - tombstones = 0
                }
              - bucket_id               = "YjE="
              - bucket_name             = "test-bucket"
              - cloud_provider          = "hostedAWS"
              - cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
              - cycle_id                = "99c971b6-451c-480a-93cf-5313ff13005d"
              - date                    = "2023-11-13T19:05:33.760979469Z"
              - elapsed_time_in_seconds = 6
              - id                      = "92ca9452-e5cb-4c16-923b-985858448e09"
              - method                  = "full"
              - organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
              - restore_before          = "2023-12-14T00:00:00Z"
              - schedule_info           = {
                  - backup_time = "2023-11-13 19:05:33.760979469 +0000 UTC"
                  - backup_type = "Manual"
                  - increment   = 1
                  - retention   = "30days"
                }
              - source                  = "manual"
              - status                  = "ready"
            },
          - {
              - backup_stats            = {
                  - cbas       = 0
                  - event      = 0
                  - fts        = 0
                  - gsi        = 0
                  - items      = 0
                  - mutations  = 0
                  - size_in_mb = 0.000527
                  - tombstones = 0
                }
              - bucket_id               = "YjE="
              - bucket_name             = "test-bucket"
              - cloud_provider          = "hostedAWS"
              - cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
              - cycle_id                = "635c196c-f91c-4c30-a33e-66fd1fa86b51"
              - date                    = "2023-11-13T17:37:26.24907334Z"
              - elapsed_time_in_seconds = 6
              - id                      = "0a04a68d-7a05-4189-9d2b-9ae4b5e3e230"
              - method                  = "full"
              - organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
              - project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
              - restore_before          = "2023-12-14T00:00:00Z"
              - schedule_info           = {
                  - backup_time = "2023-11-13 17:37:26.24907334 +0000 UTC"
                  - backup_type = "Manual"
                  - increment   = 1
                  - retention   = "30days"
                }
              - source                  = "manual"
              - status                  = "ready"
            },
        ]
      - organization_id = "c2e9ccf6-4293-4635-9205-1204d074447d"
      - project_id      = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
    } -> null
  - new_backup   = {
      - backup_stats            = {
          - cbas       = 0
          - event      = 0
          - fts        = 0
          - gsi        = 0
          - items      = 0
          - mutations  = 0
          - size_in_mb = 0.000527
          - tombstones = 0
        }
      - bucket_id               = "YjE="
      - bucket_name             = "test-bucket"
      - cloud_provider          = "hostedAWS"
      - cluster_id              = "1f6bad22-602f-407b-a567-7a8f672db227"
      - cycle_id                = "f37575d4-6531-4732-9ad8-734f1831e32e"
      - date                    = "2023-11-13T19:15:00.152667728Z"
      - elapsed_time_in_seconds = 8
      - id                      = "58dd0f30-323b-461c-83a8-1d2719f4bcee"
      - method                  = "full"
      - organization_id         = "c2e9ccf6-4293-4635-9205-1204d074447d"
      - project_id              = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      - restore_before          = "2023-12-14T00:00:00Z"
      - schedule_info           = {
          - backup_time = "2023-11-13 19:15:00.152667728 +0000 UTC"
          - backup_type = "Manual"
          - increment   = 1
          - retention   = "30days"
        }
      - source                  = "manual"
      - status                  = "ready"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_backup.new_backup: Destroying... [id=58dd0f30-323b-461c-83a8-1d2719f4bcee]
capella_backup.new_backup: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
