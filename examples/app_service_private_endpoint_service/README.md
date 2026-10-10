# Capella App Service Private Endpoint Service Example

This example shows how to enable/disable the Private Endpoint Service for an App Service (Sync Gateway) in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Enable the private endpoint service on the App Service as shown in `enable_service.tf` file.
2. LIST: Get the App Service private endpoint service status as shown in the `get_status.tf` file.
3. DELETE: Disable the private endpoint service on the App Service.
4. IMPORT: Import the private endpoint status into the state file.

If you check the `terraform.template.tfvars` file - you can see that we need 5 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE

Command: `terraform apply`

Sample Output:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_private_endpoint_service.new_service will be created
  + resource "couchbase-capella_app_service_private_endpoint_service" "new_service" {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + enabled         = true
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_private_endpoint_service.new_service: Creating...
couchbase-capella_app_service_private_endpoint_service.new_service: Still creating... [30s elapsed]
couchbase-capella_app_service_private_endpoint_service.new_service: Creation complete after 45s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## LIST

Command: `terraform plan`

Sample Output:
```
data.couchbase-capella_app_service_private_endpoint_service.service_status: Reading...
data.couchbase-capella_app_service_private_endpoint_service.service_status: Read complete after 0s

Changes to Outputs:
  + service_status = {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + enabled         = true
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + state           = "enabled"
      + target_state    = "enabled"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

## DELETE
### Remove the resource block in enable_service.tf

Command: `terraform apply`

Sample Output:
```
couchbase-capella_app_service_private_endpoint_service.new_service: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_service_private_endpoint_service.new_service will be destroyed
  # (because couchbase-capella_app_service_private_endpoint_service.new_service is not in configuration)
  - resource "couchbase-capella_app_service_private_endpoint_service" "new_service" {
      - app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - enabled         = true -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_private_endpoint_service.new_service: Destroying...
couchbase-capella_app_service_private_endpoint_service.new_service: Destruction complete after 20s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.
```

## IMPORT
### Ensure a resource block is configured as shown in enable_service.tf

Command: `terraform import couchbase-capella_app_service_private_endpoint_service.new_service app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
couchbase-capella_app_service_private_endpoint_service.new_service: Importing from ID "app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_service_private_endpoint_service.new_service: Import prepared!
  Prepared couchbase-capella_app_service_private_endpoint_service for import
couchbase-capella_app_service_private_endpoint_service.new_service: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```
