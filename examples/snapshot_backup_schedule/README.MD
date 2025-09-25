# Capella Cloud Snapshot Backup Schedule Example

This example shows how to create and manage Snapshot Backup Schedules in Capella.

This creates a new snapshot backup schedule in the selected Capella cluster. It uses the organization ID, project ID, and cluster ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new snapshot backup schedule entry in an existing Capella cluster as stated in the `create_snapshot_backup_schedule.tf` file.
2. UPDATE: Update the snapshot backup schedule configuration using Terraform.
3. IMPORT: Import a snapshot backup schedule that exists in Capella but not in the terraform state file.
4. DELETE: Delete the newly created snapshot backup schedule from Capella.

If you wish to use the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & READ
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

  # couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule will be created
  + resource "couchbase-capella_snapshot_backup_schedule" "new_snapshot_backup_schedule" {
      + copy_to_regions = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      + interval        = 4
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
      + retention       = 24
      + start_time      = "2025-09-15T16:00:00+00:00"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_snapshot_backup_schedule = {
      + copy_to_regions = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      + interval        = 4
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
      + retention       = 24
      + start_time      = "2025-09-15T16:00:00+00:00"
    }
```

### Apply the Plan, in order to create a new Snapshot Backup Schedule entry

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

  # couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule will be created
  + resource "couchbase-capella_snapshot_backup_schedule" "new_snapshot_backup_schedule" {
      + copy_to_regions = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      + interval        = 4
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
      + retention       = 24
      + start_time      = "2025-09-15T16:00:00+00:00"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_snapshot_backup_schedule = {
      + copy_to_regions = [
          + "eu-west-1",
          + "ap-southeast-1",
        ]
      + id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      + interval        = 4
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + project_id      = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
      + retention       = 24
      + start_time      = "2025-09-15T16:00:00+00:00"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Creating...
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Creation complete after 0s [id=eaa5facc-2a30-463a-8ffc-e1df47539028]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_snapshot_backup_schedule = {
  "copy_to_regions" = tolist([
    "eu-west-1",
    "ap-southeast-1",
  ])
  "id" = "eaa5facc-2a30-463a-8ffc-e1df47539028"
  "interval" = 4
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
  "retention" = 24
  "start_time" = "2025-09-15T16:00:00+00:00"
}
```

### Note the Cluster ID of the new Snapshot Backup Schedule required for import.
Command: `terraform output new_snapshot_backup_schedule`

Sample Output:
```
{
  "copy_to_regions" = tolist([
    "eu-west-1",
    "ap-southeast-1",
  ])
  "id" = "eaa5facc-2a30-463a-8ffc-e1df47539028"
  "interval" = 4
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
  "retention" = 24
  "start_time" = "2025-09-15T16:00:00+00:00"
}
```
In this case, the cluster ID is `eaa5facc-2a30-463a-8ffc-e1df47539028`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule
```

## IMPORT
### Remove the resource `new_snapshot_backup_schedule` from the Terraform State file

Command: `terraform state rm couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule`

Sample Output:
``` 
Removed couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule "id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>"`

In this case, the complete command is:
`terraform import couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule "id=eaa5facc-2a30-463a-8ffc-e1df47539028,project_id=d6552b0c-ca66-4412-b1e1-68bd36a7e364,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76"`

Sample Output:
``` 
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Importing from ID "id=eaa5facc-2a30-463a-8ffc-e1df47539028,project_id=d6552b0c-ca66-4412-b1e1-68bd36a7e364,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76"...
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Import prepared!
  Prepared couchbase-capella_snapshot_backup_schedule for import
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Refreshing state... [id=id=eaa5facc-2a30-463a-8ffc-e1df47539028,project_id=d6552b0c-ca66-4412-b1e1-68bd36a7e364,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the cluster ID i.e. the ID of the cluster to which the snapshot backup schedule belongs.
The second ID is the project ID i.e. the ID of the project to which the cluster belongs.
The third ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit terraform.tfvars file to change the snapshot backup schedule configuration settings.

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
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Refreshing state... [id=eaa5facc-2a30-463a-8ffc-e1df47539028]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule will be updated in-place
  ~ resource "couchbase-capella_snapshot_backup_schedule" "new_snapshot_backup_schedule" {
      ~ copy_to_regions = [
          - "eu-west-1",
          - "ap-southeast-1",
          + "af-south-1",
        ]
        id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      ~ interval        = 4 -> 8
      ~ retention       = 24 -> 168
      ~ start_time      = "2025-09-15T16:00:00Z" -> "2026-01-23T03:30:00+00:00"
        # (2 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_snapshot_backup_schedule = {
      ~ copy_to_regions = [
          - "eu-west-1",
          - "ap-southeast-1",
          + "af-south-1",
        ]
        id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      ~ interval        = 4 -> 8
      ~ retention       = 24 -> 168
      ~ start_time      = "2025-09-15T16:00:00Z" -> "2026-01-23T03:30:00+00:00"
        # (2 unchanged attributes hidden)
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
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Refreshing state... [id=eaa5facc-2a30-463a-8ffc-e1df47539028]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule will be updated in-place
  ~ resource "couchbase-capella_snapshot_backup_schedule" "new_snapshot_backup_schedule" {
      ~ copy_to_regions = [
          - "eu-west-1",
          - "ap-southeast-1",
          + "af-south-1",
        ]
        id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      ~ interval        = 4 -> 8
      ~ retention       = 24 -> 168
      ~ start_time      = "2025-09-15T16:00:00Z" -> "2026-01-23T03:30:00+00:00"
        # (2 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_snapshot_backup_schedule = {
      ~ copy_to_regions = [
          - "eu-west-1",
          - "ap-southeast-1",
          + "af-south-1",
        ]
        id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      ~ interval        = 4 -> 8
      ~ retention       = 24 -> 168
      ~ start_time      = "2025-09-15T16:00:00Z" -> "2026-01-23T03:30:00+00:00"
        # (2 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Modifying... [id=eaa5facc-2a30-463a-8ffc-e1df47539028]
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Modifications complete after 0s [id=eaa5facc-2a30-463a-8ffc-e1df47539028]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_snapshot_backup_schedule = {
  "copy_to_regions" = tolist([
    "af-south-1",
  ])
  "id" = "eaa5facc-2a30-463a-8ffc-e1df47539028"
  "interval" = 8
  "organization_id" = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
  "project_id" = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
  "retention" = 168
  "start_time" = "2026-01-23T03:30:00+00:00"
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
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Refreshing state... [id=eaa5facc-2a30-463a-8ffc-e1df47539028]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule will be destroyed
  - resource "couchbase-capella_snapshot_backup_schedule" "new_snapshot_backup_schedule" {
      - copy_to_regions = [
          - "af-south-1",
        ] -> null
      - id              = "eaa5facc-2a30-463a-8ffc-e1df47539028" -> null
      - interval        = 8 -> null
      - organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76" -> null
      - project_id      = "d6552b0c-ca66-4412-b1e1-68bd36a7e364" -> null
      - retention       = 168 -> null
      - start_time      = "2026-01-23T03:30:00+00:00" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_snapshot_backup_schedule = {
      - copy_to_regions = [
          - "af-south-1",
        ]
      - id              = "eaa5facc-2a30-463a-8ffc-e1df47539028"
      - interval        = 8
      - organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      - project_id      = "d6552b0c-ca66-4412-b1e1-68bd36a7e364"
      - retention       = 168
      - start_time      = "2026-01-23T03:30:00+00:00"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Destroying... [id=eaa5facc-2a30-463a-8ffc-e1df47539028]
couchbase-capella_snapshot_backup_schedule.new_snapshot_backup_schedule: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
