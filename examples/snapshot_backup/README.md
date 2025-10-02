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
│  - couchbasecloud/couchbase-capella in /Users/sophie.wegmann/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup.new_snapshot_backup will be created
  + resource "couchbase-capella_snapshot_backup" "new_snapshot_backup" {
      + app_service         = (known after apply)
      + cluster_id          = "caf9f40b-0078-4e37-8700-7bdb2e42ae8f"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress            = (known after apply)
      + project_id          = "c5b4fcc2-3158-4e23-addf-52157a6e6ae0"
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
  + new_snapshot_backup = {
      + app_service                     = (known after apply)
      + cluster_id                      = "caf9f40b-0078-4e37-8700-7bdb2e42ae8f"
      + cmek                            = (known after apply)
      + created_at                      = (known after apply)
      + cross_region_copies             = (known after apply)
      + cross_region_restore_preference = null
      + expiration                      = (known after apply)
      + id                              = (known after apply)
      + organization_id                 = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress                        = (known after apply)
      + project_id                      = "c5b4fcc2-3158-4e23-addf-52157a6e6ae0"
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
│  - couchbasecloud/couchbase-capella in /Users/sophie.wegmann/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup.new_snapshot_backup will be created
  + resource "couchbase-capella_snapshot_backup" "new_snapshot_backup" {
      + app_service         = (known after apply)
      + cluster_id          = "caf9f40b-0078-4e37-8700-7bdb2e42ae8f"
      + cmek                = (known after apply)
      + created_at          = (known after apply)
      + cross_region_copies = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress            = (known after apply)
      + project_id          = "c5b4fcc2-3158-4e23-addf-52157a6e6ae0"
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
  + new_snapshot_backup = {
      + app_service                     = (known after apply)
      + cluster_id                      = "caf9f40b-0078-4e37-8700-7bdb2e42ae8f"
      + cmek                            = (known after apply)
      + created_at                      = (known after apply)
      + cross_region_copies             = (known after apply)
      + cross_region_restore_preference = null
      + expiration                      = (known after apply)
      + id                              = (known after apply)
      + organization_id                 = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + progress                        = (known after apply)
      + project_id                      = "c5b4fcc2-3158-4e23-addf-52157a6e6ae0"
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
couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup: Creation complete after 0s [id=cd7206eb-1914-46b0-832d-1257936b149b]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_snapshot_backup = {
  "app_service" = ""
  "cluster_id" = "caf9f40b-0078-4e37-8700-7bdb2e42ae8f"
  "cmek" = toset([])
  "created_at" = "2025-09-26T16:37:32.127711554Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-09-26T16:39:48.237639173Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-09-26T16:40:22.145533091Z"
    },
  ])
  "cross_region_restore_preference" = tolist(null) /* of string */
  "expiration" = "2025-10-02T16:37:32.127711554Z"
  "id" = "d4faa99f-8967-42ce-8769-3d27ba5a30f4"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "progress" = {
    "status" = "complete"
    "time" = "2025-09-26T16:38:24.102226342Z"
  }
  "project_id" = "c5b4fcc2-3158-4e23-addf-52157a6e6ae0"
  "regions_to_copy" = tolist([
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

### Note the Snapshot Backup ID of the new Snapshot Backup
Command: `terraform output new_cloud_snapshot_backup`

Sample Output:
```
{
  "app_service" = ""
  "cluster_id" = "caf9f40b-0078-4e37-8700-7bdb2e42ae8f"
  "cmek" = toset([])
  "created_at" = "2025-09-26T16:37:32.127711554Z"
  "cross_region_copies" = toset([
    {
      "region_code" = "ap-southeast-1"
      "status" = "complete"
      "time" = "2025-09-26T16:39:48.237639173Z"
    },
    {
      "region_code" = "eu-west-1"
      "status" = "complete"
      "time" = "2025-09-26T16:40:22.145533091Z"
    },
  ])
  "cross_region_restore_preference" = tolist(null) /* of string */
  "expiration" = "2025-10-02T16:37:32.127711554Z"
  "id" = "d4faa99f-8967-42ce-8769-3d27ba5a30f4"
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "progress" = {
    "status" = "complete"
    "time" = "2025-09-26T16:38:24.102226342Z"
  }
  "project_id" = "c5b4fcc2-3158-4e23-addf-52157a6e6ae0"
  "regions_to_copy" = tolist([
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
In this case, the snapshot backup ID for my new snapshot backup is `d4faa99f-8967-42ce-8769-3d27ba5a30f4`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_snapshot_backup.new_snapshot_backup
```

## IMPORT
### Remove the resource `new_snapshot_backup` from the Terraform State file

Command: `terraform state rm couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup`

Sample Output:
``` 
Removed couchbase-capella_snapshot_backup.new_snapshot_backup
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_cloud_snapshot_backup.new_cloud_snapshot_backup id=<snapshot_backup_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_snapshot_backup.new_snapshot_backup id=d4faa99f-8967-42ce-8769-3d27ba5a30f4,cluster_id=caf9f40b-0078-4e37-8700-7bdb2e42ae8f,project_id=c5b4fcc2-3158-4e23-addf-52157a6e6ae0,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76`

Sample Output:
``` 
couchbase-capella_snapshot_backup.new_snapshot_backup: Importing from ID "id=d4faa99f-8967-42ce-8769-3d27ba5a30f4,cluster_id=caf9f40b-0078-4e37-8700-7bdb2e42ae8f,project_id=c5b4fcc2-3158-4e23-addf-52157a6e6ae0,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76"...
couchbase-capella_snapshot_backup.new_snapshot_backup: Import prepared!
  Prepared couchbase-capella_snapshot_backup for import
couchbase-capella_snapshot_backup.new_snapshot_backup: Refreshing state... [id=id=d4faa99f-8967-42ce-8769-3d27ba5a30f4,cluster_id=caf9f40b-0078-4e37-8700-7bdb2e42ae8f,project_id=c5b4fcc2-3158-4e23-addf-52157a6e6ae0,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76]

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
```

command: `terrafom apply`

Sample Output:

```
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
```
