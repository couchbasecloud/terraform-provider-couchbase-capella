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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule will be created
  + resource "couchbase-capella_cloud_snapshot_backup_schedule" "new_cloud_snapshot_backup_schedule" {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + copy_to_regions = [
          + "af-south-1",
        ]
      + interval        = 12
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention       = 24
      + start_time      = "2025-12-15T10:00:00Z"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_cloud_snapshot_backup_schedule = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + copy_to_regions = [
          + "af-south-1",
        ]
      + interval        = 12
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention       = 24
      + start_time      = "2025-12-15T10:00:00Z"
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule will be created
  + resource "couchbase-capella_cloud_snapshot_backup_schedule" "new_cloud_snapshot_backup_schedule" {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + copy_to_regions = [
          + "af-south-1",
        ]
      + interval        = 12
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention       = 24
      + start_time      = "2025-12-15T10:00:00Z"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_cloud_snapshot_backup_schedule = {
      + cluster_id      = "36c03330-94c5-44a6-b812-265b10098f33"
      + copy_to_regions = [
          + "af-south-1",
        ]
      + interval        = 12
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + retention       = 24
      + start_time      = "2025-12-15T10:00:00Z"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Creating...
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_cloud_snapshot_backup_schedule = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "copy_to_regions" = toset([
    "af-south-1",
  ])
  "interval" = 12
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "retention" = 24
  "start_time" = "2025-12-15T10:00:00Z"
}
```

### Note the Cluster ID of the new Snapshot Backup Schedule required for import.
Command: `terraform output new_cloud_snapshot_backup_schedule`

Sample Output:
```
{
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "copy_to_regions" = toset([
    "af-south-1",
  ])
  "interval" = 12
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "retention" = 24
  "start_time" = "2025-12-15T10:00:00Z"
}
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule
```

## IMPORT
### Remove the resource `new_cloud_snapshot_backup_schedule` from the Terraform State file

Command: `terraform state rm couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule`

Sample Output:
``` 
Removed couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule "cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>"`

In this case, the complete command is:
`terraform import couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule "cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"`

Sample Output:
``` 
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Importing from ID "cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Import prepared!
  Prepared couchbase-capella_cloud_snapshot_backup_schedule for import
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Refreshing state...

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule will be updated in-place
  ~ resource "couchbase-capella_cloud_snapshot_backup_schedule" "new_cloud_snapshot_backup_schedule" {
      ~ copy_to_regions = [
          - "af-south-1",
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      ~ interval        = 12 -> 8
      ~ retention       = 24 -> 48
      ~ start_time      = "2025-12-15T10:00:00Z" -> "2025-12-16T10:00:00Z"
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cloud_snapshot_backup_schedule = {
      ~ copy_to_regions = [
          - "af-south-1",
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      ~ interval        = 12 -> 8
      ~ retention       = 24 -> 48
      ~ start_time      = "2025-12-15T10:00:00Z" -> "2025-12-16T10:00:00Z"
        # (3 unchanged attributes hidden)
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule will be updated in-place
  ~ resource "couchbase-capella_cloud_snapshot_backup_schedule" "new_cloud_snapshot_backup_schedule" {
      ~ copy_to_regions = [
          - "af-south-1",
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      ~ interval        = 12 -> 8
      ~ retention       = 24 -> 48
      ~ start_time      = "2025-12-15T10:00:00Z" -> "2025-12-16T10:00:00Z"
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cloud_snapshot_backup_schedule = {
      ~ copy_to_regions = [
          - "af-south-1",
          + "ap-southeast-1",
          + "eu-west-1",
        ]
      ~ interval        = 12 -> 8
      ~ retention       = 24 -> 48
      ~ start_time      = "2025-12-15T10:00:00Z" -> "2025-12-16T10:00:00Z"
        # (3 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Modifying...
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Modifications complete after 0s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_cloud_snapshot_backup_schedule = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "copy_to_regions" = toset([
    "ap-southeast-1",
    "eu-west-1",
  ])
  "interval" = 8
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "retention" = 48
  "start_time" = "2025-12-16T10:00:00Z"
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule will be destroyed
  - resource "couchbase-capella_cloud_snapshot_backup_schedule" "new_cloud_snapshot_backup_schedule" {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - copy_to_regions = [
          - "ap-southeast-1",
          - "eu-west-1",
        ] -> null
      - interval        = 8 -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - retention       = 48 -> null
      - start_time      = "2025-12-16T10:00:00Z" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_cloud_snapshot_backup_schedule = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - copy_to_regions = [
          - "ap-southeast-1",
          - "eu-west-1",
        ]
      - interval        = 8
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - retention       = 48
      - start_time      = "2025-12-16T10:00:00Z"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Destroying...
couchbase-capella_cloud_snapshot_backup_schedule.new_cloud_snapshot_backup_schedule: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
