# Capella App Endpoint Resync Example

This example shows how to initiate and manage App Endpoint Resync operations in Capella.

It uses the organization ID, project ID, cluster ID, app service ID, and app endpoint name to start a resync. 
You can add collections to be resynced within the scopes attribute on the resource. The specified collections will be resynced, along with any collections that require a resync.

To run, configure your Couchbase Capella provider as described in the README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following:

1. CREATE: Initiate a resync operation for an App Endpoint as defined in `create_app_endpoint_resync.tf`.
2. DESTROY: Stop an ongoing resync operation.
3. IMPORT: Import the status an existing resync operation into Terraform state.

If you check the `terraform.template.tfvars` file — copy it to `terraform.tfvars` and update variable values for your environment.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Example output:
```
 $ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_resync.r1 will be created
  + resource "couchbase-capella_app_endpoint_resync" "r1" {
      + app_endpoint_name      = "api"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = (known after apply)
      + docs_changed           = (known after apply)
      + docs_processed         = (known after apply)
      + last_error             = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + start_time             = (known after apply)
      + state                  = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the plan to initiate a resync

Command: `terraform apply`

```

$ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_resync.r1 will be created
  + resource "couchbase-capella_app_endpoint_resync" "r1" {
      + app_endpoint_name      = "api"
      + app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + collections_processing = (known after apply)
      + docs_changed           = (known after apply)
      + docs_processed         = (known after apply)
      + last_error             = (known after apply)
      + organization_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id             = "ffffffff-aaaa-1414-eeee-000000000000"
      + start_time             = (known after apply)
      + state                  = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_resync.r1: Creating...
couchbase-capella_app_endpoint_resync.r1: Creation complete after 8s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## DESTROY
### Stop the resync operation managed by Terraform

Command: `terraform destroy`

Example output:
```
 $ terraform destroy                

couchbase-capella_app_endpoint_resync.r1: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_resync.r1 will be destroyed
  - resource "couchbase-capella_app_endpoint_resync" "r1" {
      - app_endpoint_name      = "api" -> null
      - app_service_id         = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collections_processing = {
          - "_default" = [
              - "_default",
            ]
        } -> null
      - docs_changed           = 0 -> null
      - docs_processed         = 0 -> null
      - organization_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id             = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - start_time             = "2025-09-24T14:52:03Z" -> null
      - state                  = "completed" -> null
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_endpoint_resync.r1: Destroying...
couchbase-capella_app_endpoint_resync.r1: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
## IMPORT
### Import an existing resync into Terraform state

Ensure the resource block has the correct `organization_id`, `project_id`, `cluster_id`, `app_service_id` and `app_endpoint_name` configured. Then run:

Command: `terraform import couchbase-capella_app_endpoint_resync.example app_endpoint_name=<name>,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000`

Example output
```
 $ terraform import  couchbase-capella_app_endpoint_resync.r1 app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_app_endpoint_resync.r1: Importing from ID "app_endpoint_name=api,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_resync.r1: Import prepared!
  Prepared couchbase-capella_app_endpoint_resync for import
couchbase-capella_app_endpoint_resync.r1: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```