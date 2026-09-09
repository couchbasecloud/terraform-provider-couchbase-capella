# Capella App Service Private Endpoints Example

This example shows how to accept and list private endpoint connections for an App Service (Sync Gateway) in Capella. The App Service private endpoint service (see the `app_service_private_endpoint_service` example) must already be enabled, and the corresponding CSP-side private endpoint (VPC endpoint, Private Link, or PSC connection) must already exist and be pending acceptance.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Accept a pending private endpoint connection as shown in `accept_endpoint.tf` file.
2. LIST: List all private endpoints connected to the App Service's private endpoint service as shown in the `list_endpoints.tf` file.
3. DELETE: Reject (remove) the accepted private endpoint connection.
4. IMPORT: Import an accepted private endpoint connection into the state file.

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

  # couchbase-capella_app_service_private_endpoints.accept_endpoint will be created
  + resource "couchbase-capella_app_service_private_endpoints" "accept_endpoint" {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + endpoint_id     = "vpce-0123456789abcdef0"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + service_name    = (known after apply)
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_private_endpoints.accept_endpoint: Creating...
couchbase-capella_app_service_private_endpoints.accept_endpoint: Creation complete after 5s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## LIST

Command: `terraform plan`

Sample Output:
```
data.couchbase-capella_app_service_private_endpoints.list_endpoints: Reading...
data.couchbase-capella_app_service_private_endpoints.list_endpoints: Read complete after 0s

Changes to Outputs:
  + list_endpoints = {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + id           = "vpce-0123456789abcdef0"
              + service_name = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
              + status       = "linked"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.
```

## DELETE
### Remove the resource block in accept_endpoint.tf

Command: `terraform apply`

Sample Output:
```
couchbase-capella_app_service_private_endpoints.accept_endpoint: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_service_private_endpoints.accept_endpoint will be destroyed
  # (because couchbase-capella_app_service_private_endpoints.accept_endpoint is not in configuration)
  - resource "couchbase-capella_app_service_private_endpoints" "accept_endpoint" {
      - app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - endpoint_id     = "vpce-0123456789abcdef0" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - service_name    = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0" -> null
      - status          = "linked" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_private_endpoints.accept_endpoint: Destroying...
couchbase-capella_app_service_private_endpoints.accept_endpoint: Destruction complete after 3s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.
```

## IMPORT
### Ensure a resource block is configured as shown in accept_endpoint.tf

Command: `terraform import couchbase-capella_app_service_private_endpoints.accept_endpoint endpoint_id=<endpoint_id>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
couchbase-capella_app_service_private_endpoints.accept_endpoint: Importing from ID "endpoint_id=vpce-0123456789abcdef0,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_service_private_endpoints.accept_endpoint: Import prepared!
  Prepared couchbase-capella_app_service_private_endpoints for import
couchbase-capella_app_service_private_endpoints.accept_endpoint: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```
