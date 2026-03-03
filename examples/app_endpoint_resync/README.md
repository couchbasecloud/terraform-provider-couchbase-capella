# Capella App Endpoint Resync Example

This example shows how to start and stop an App Endpoint Resync in Capella. It uses the organization ID, project ID, cluster ID, App Service ID, and App Endpoint name to do so.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Start a resync for a Capella app endpoint as stated in the `create_app_endpoint_resync.tf` file.
2. IMPORT: Read the current resync status that exists in Capella but not in the Terraform state file.
3. DELETE: Stop the resync, and remove the App Endpoint Resync resource from the Terraform state file.

The resource cannot be updated because an app endpoint resync operation cannot be updated.

If you wish to use the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & READ
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job will be created
  + resource "couchbase-capella_app_endpoint_resync_job" "new_app_endpoint_resync_job" {
      + app_endpoint_name      = "test_app_endpoint"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = (known after apply)
      + docs_changed           = (known after apply)
      + docs_processed         = (known after apply)
      + last_error             = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + scopes                 = {
          + "inventory" = [
              + "airline",
              + "hotel",
            ]
        }
      + start_time             = (known after apply)
      + state                  = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_endpoint_resync_job = {
      + app_endpoint_name      = "test_app_endpoint"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = (known after apply)
      + docs_changed           = (known after apply)
      + docs_processed         = (known after apply)
      + last_error             = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + scopes                 = {
          + inventory = [
              + "airline",
              + "hotel",
            ]
        }
      + start_time             = (known after apply)
      + state                  = (known after apply)
    }
```

### Apply the Plan, in order to create the App Endpoint Resync Resource

Command: `terraform apply`

Sample Output:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job will be created
  + resource "couchbase-capella_app_endpoint_resync_job" "new_app_endpoint_resync_job" {
      + app_endpoint_name      = "test_app_endpoint"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = (known after apply)
      + docs_changed           = (known after apply)
      + docs_processed         = (known after apply)
      + last_error             = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + scopes                 = {
          + "inventory" = [
              + "airline",
              + "hotel",
            ]
        }
      + start_time             = (known after apply)
      + state                  = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_endpoint_resync_job = {
      + app_endpoint_name      = "test_app_endpoint"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = (known after apply)
      + docs_changed           = (known after apply)
      + docs_processed         = (known after apply)
      + last_error             = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + scopes                 = {
          + inventory = [
              + "airline",
              + "hotel",
            ]
        }
      + start_time             = (known after apply)
      + state                  = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Creating...
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_endpoint_resync_job = {
  "app_endpoint_name" = "test_app_endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collections_processing" = tomap({
    "inventory" = toset([
      "airline",
      "hotel",
    ])
  })
  "docs_changed" = 0
  "docs_processed" = 0
  "last_error" = ""
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scopes" = tomap({
    "inventory" = toset([
      "airline",
      "hotel",
    ])
  })
  "start_time" = "2026-03-03T15:40:56Z"
  "state" = "running"
}
```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job
```

## IMPORT
### Remove the resource `new_app_endpoint_resync_job` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job`

Sample Output:
``` 
Removed couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job
Successfully removed 1 resource instance(s).
```

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job app_endpoint_name=<app_endpoint_name>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job app_endpoint_name=test_app_endpoint,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Importing from ID "app_endpoint_name=test_app_endpoint,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-0000000000006"...
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Import prepared!
  Prepared couchbase-capella_app_endpoint_resync_job for import
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the app endpoint name i.e. the name of the app endpoint to which the resync belongs.
The second ID is the app service ID i.e. the ID of the app service to which the app endpoint belongs.
The third ID is the cluster ID i.e. the ID of the cluster to which the app service belongs.
The fourth ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fifth ID is the organization ID i.e. the ID of the organization to which the project belongs.


## DESTROY

Note that destroying this resource stops the current app endpoint resync.

### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job will be destroyed
  - resource "couchbase-capella_app_endpoint_resync_job" "new_app_endpoint_resync_job" {
      - app_endpoint_name      = "test_app_endpoint" -> null
      - app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collections_processing = {
          - "inventory" = [
              - "airline",
              - "hotel",
            ]
        } -> null
      - docs_changed           = 0 -> null
      - docs_processed         = 0 -> null
      - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - start_time             = "2026-03-03T17:23:54Z" -> null
      - state                  = "completed" -> null
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_app_endpoint_resync_job = {
      - app_endpoint_name      = "test_app_endpoint"
      - app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      - collections_processing = {
          - inventory = [
              - "airline",
              - "hotel",
            ]
        }
      - docs_changed           = 0
      - docs_processed         = 0
      - last_error             = ""
      - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      - scopes                 = null
      - start_time             = "2026-03-03T17:23:54Z"
      - state                  = "completed"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Destroying...
couchbase-capella_app_endpoint_resync_job.new_app_endpoint_resync_job: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```