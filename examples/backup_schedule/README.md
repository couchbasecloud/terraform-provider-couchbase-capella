# Capella Backup Schedule Example

This example shows how to create and manage Backup Schedules in Capella.

This creates a new backup schedule in the selected Capella cluster. It uses the organization ID, project ID, cluster ID and bucket ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new backup schedule entry in an existing Capella cluster as stated in the `create_backup_schedule.tf` file.
2. UPDATE: Update the backup schedule configuration using Terraform.
3. IMPORT: Import a backup schedule that exists in Capella but not in the terraform state file.
4. DELETE: Delete the newly created backup schedule from Capella.

If you wish to use the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & READ
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

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_backup_schedule.new_backup_schedule will be created
  + resource "capella_backup_schedule" "new_backup_schedule" {
      + bucket_id       = "dGVzdC1idWNrZXQ="
      + cluster_id      = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = "53b0e002-eb11-4317-9a53-6a781b29960e"
      + type            = "weekly"
      + weekly_schedule = {
          + cost_optimized_retention = false
          + day_of_week              = "Monday"
          + incremental_every        = 6
          + retention_time           = "30days"
          + start_at                 = 10
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_backup_schedule = {
      + bucket_id       = "dGVzdC1idWNrZXQ="
      + cluster_id      = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = "53b0e002-eb11-4317-9a53-6a781b29960e"
      + type            = "weekly"
      + weekly_schedule = {
          + cost_optimized_retention = false
          + day_of_week              = "Monday"
          + incremental_every        = 6
          + retention_time           = "30days"
          + start_at                 = 10
        }
    }
```

### Apply the Plan, in order to create a new Backup Schedule entry

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

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_backup_schedule.new_backup_schedule will be created
  + resource "capella_backup_schedule" "new_backup_schedule" {
      + bucket_id       = "dGVzdC1idWNrZXQ="
      + cluster_id      = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = "53b0e002-eb11-4317-9a53-6a781b29960e"
      + type            = "weekly"
      + weekly_schedule = {
          + cost_optimized_retention = false
          + day_of_week              = "Monday"
          + incremental_every        = 6
          + retention_time           = "30days"
          + start_at                 = 10
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_backup_schedule = {
      + bucket_id       = "dGVzdC1idWNrZXQ="
      + cluster_id      = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
      + organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      + project_id      = "53b0e002-eb11-4317-9a53-6a781b29960e"
      + type            = "weekly"
      + weekly_schedule = {
          + cost_optimized_retention = false
          + day_of_week              = "Monday"
          + incremental_every        = 6
          + retention_time           = "30days"
          + start_at                 = 10
        }
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_backup_schedule.new_backup_schedule: Creating...
capella_backup_schedule.new_backup_schedule: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_backup_schedule = {
  "bucket_id" = "dGVzdC1idWNrZXQ="
  "cluster_id" = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "53b0e002-eb11-4317-9a53-6a781b29960e"
  "type" = "weekly"
  "weekly_schedule" = {
    "cost_optimized_retention" = false
    "day_of_week" = "Monday"
    "incremental_every" = 6
    "retention_time" = "30days"
    "start_at" = 10
  }
}
```

### Note the Bucket ID of the new Backup Schedule required for import.
Command: `terraform output new_backup_schedule`

Sample Output:
```
terraform output new_backup_schedule           
{
  "bucket_id" = "dGVzdC1idWNrZXQ="
  "cluster_id" = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "53b0e002-eb11-4317-9a53-6a781b29960e"
  "type" = "weekly"
  "weekly_schedule" = {
    "cost_optimized_retention" = false
    "day_of_week" = "Monday"
    "incremental_every" = 6
    "retention_time" = "30days"
    "start_at" = 10
  }
}
```
In this case, the bucket ID is `dGVzdC1idWNrZXQ=`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
terraform state list                                  
couchbase-capella_backup_schedule.new_backup_schedule
```

## IMPORT
### Remove the resource `new_backup_schedule` from the Terraform State file

Command: `terraform state rm couchbase-capella_backup_schedule.new_backup_schedule`

Sample Output:
``` 
terraform state rm couchbase-capella_backup_schedule.new_backup_schedule
Removed capella_backup_schedule.new_backup_schedule
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_backup.new_backup id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_backup_schedule.new_backup_schedule bucket_id=dGVzdC1idWNrZXQ=,cluster_id=4f728ab7-dbbc-45a2-9789-ee172f09851e,project_id=53b0e002-eb11-4317-9a53-6a781b29960e,organization_id=7a99d00c-f55b-4b39-bc72-1b4cc68ba894`

Sample Output:
``` 
terraform import couchbase-capella_backup_schedule.new_backup_schedule bucket_id=dGVzdC1idWNrZXQ=,cluster_id=4f728ab7-dbbc-45a2-9789-ee172f09851e,project_id=53b0e002-eb11-4317-9a53-6a781b29960e,organization_id=7a99d00c-f55b-4b39-bc72-1b4cc68ba894
capella_backup_schedule.new_backup_schedule: Importing from ID "bucket_id=dGVzdC1idWNrZXQ=,cluster_id=4f728ab7-dbbc-45a2-9789-ee172f09851e,project_id=53b0e002-eb11-4317-9a53-6a781b29960e,organization_id=7a99d00c-f55b-4b39-bc72-1b4cc68ba894"...
capella_backup_schedule.new_backup_schedule: Import prepared!
  Prepared capella_backup_schedule for import
capella_backup_schedule.new_backup_schedule: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the bucket ID i.e. the ID of the bucket to which the backup schedule belongs.
The second ID is the cluster ID i.e. the ID of the cluster to which bucket belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit terraform.tfvars file to change the backup schedule configuration settings.

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
capella_backup_schedule.new_backup_schedule: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_backup_schedule.new_backup_schedule will be updated in-place
  ~ resource "capella_backup_schedule" "new_backup_schedule" {
      ~ weekly_schedule = {
          ~ cost_optimized_retention = false -> true
          ~ day_of_week              = "Monday" -> "Friday"
          ~ incremental_every        = 6 -> 4
          ~ retention_time           = "30days" -> "60days"
          ~ start_at                 = 10 -> 7
        }
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_backup_schedule = {
      ~ weekly_schedule = {
          ~ cost_optimized_retention = false -> true
          ~ day_of_week              = "Monday" -> "Friday"
          ~ incremental_every        = 6 -> 4
          ~ retention_time           = "30days" -> "60days"
          ~ start_at                 = 10 -> 7
        }
        # (5 unchanged attributes hidden)
    }
```

command: `terrafom apply`

Sample Output:

```
$ terraform apply                                                                              
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
capella_backup_schedule.new_backup_schedule: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_backup_schedule.new_backup_schedule will be updated in-place
  ~ resource "capella_backup_schedule" "new_backup_schedule" {
      ~ weekly_schedule = {
          ~ cost_optimized_retention = false -> true
          ~ day_of_week              = "Monday" -> "Friday"
          ~ incremental_every        = 6 -> 4
          ~ retention_time           = "30days" -> "60days"
          ~ start_at                 = 10 -> 7
        }
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_backup_schedule = {
      ~ weekly_schedule = {
          ~ cost_optimized_retention = false -> true
          ~ day_of_week              = "Monday" -> "Friday"
          ~ incremental_every        = 6 -> 4
          ~ retention_time           = "30days" -> "60days"
          ~ start_at                 = 10 -> 7
        }
        # (5 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_backup_schedule.new_backup_schedule: Modifying...
capella_backup_schedule.new_backup_schedule: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_backup_schedule = {
  "bucket_id" = "dGVzdC1idWNrZXQ="
  "cluster_id" = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
  "organization_id" = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
  "project_id" = "53b0e002-eb11-4317-9a53-6a781b29960e"
  "type" = "weekly"
  "weekly_schedule" = {
    "cost_optimized_retention" = true
    "day_of_week" = "Friday"
    "incremental_every" = 4
    "retention_time" = "60days"
    "start_at" = 7
  }
}
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
capella_backup_schedule.new_backup_schedule: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_backup_schedule.new_backup_schedule will be destroyed
  - resource "capella_backup_schedule" "new_backup_schedule" {
      - bucket_id       = "dGVzdC1idWNrZXQ=" -> null
      - cluster_id      = "4f728ab7-dbbc-45a2-9789-ee172f09851e" -> null
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894" -> null
      - project_id      = "53b0e002-eb11-4317-9a53-6a781b29960e" -> null
      - type            = "weekly" -> null
      - weekly_schedule = {
          - cost_optimized_retention = true -> null
          - day_of_week              = "Friday" -> null
          - incremental_every        = 4 -> null
          - retention_time           = "60days" -> null
          - start_at                 = 7 -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_backup_schedule = {
      - bucket_id       = "dGVzdC1idWNrZXQ="
      - cluster_id      = "4f728ab7-dbbc-45a2-9789-ee172f09851e"
      - organization_id = "7a99d00c-f55b-4b39-bc72-1b4cc68ba894"
      - project_id      = "53b0e002-eb11-4317-9a53-6a781b29960e"
      - type            = "weekly"
      - weekly_schedule = {
          - cost_optimized_retention = true
          - day_of_week              = "Friday"
          - incremental_every        = 4
          - retention_time           = "60days"
          - start_at                 = 7
        }
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_backup_schedule.new_backup_schedule: Destroying...
capella_backup_schedule.new_backup_schedule: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
