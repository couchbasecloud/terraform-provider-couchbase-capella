# Capella App Service Private Endpoint Example

This example shows how to accept a private endpoint connection request on an App Service in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Accept a private endpoint connection request as shown in `accept_endpoint.tf` file.
2. DELETE: Remove the private endpoint from the App Service.
3. IMPORT: Import the private endpoint into the state file.

Before you run this example, enable the private endpoint service on the App Service and create a private endpoint on your cloud service provider that requests a connection to it. The endpoint shows as `pendingAcceptance` on the App Service.

If you check the `terraform.template.tfvars` file - you can see that we need 6 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE

Command: `terraform apply`

Terraform accepts the connection request and waits until the endpoint status is `linked`.

Sample Output:
```
terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_private_endpoint.accept_endpoint will be created
  + resource "couchbase-capella_app_service_private_endpoint" "accept_endpoint" {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + endpoint_id     = "vpce-000000000000aaaaa"
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

couchbase-capella_app_service_private_endpoint.accept_endpoint: Creating...
couchbase-capella_app_service_private_endpoint.accept_endpoint: Creation complete after 35s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## IMPORT

Command: `terraform import couchbase-capella_app_service_private_endpoint.accept_endpoint endpoint_id=<endpoint_id>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
couchbase-capella_app_service_private_endpoint.accept_endpoint: Importing from ID "endpoint_id=vpce-000000000000aaaaa,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_service_private_endpoint.accept_endpoint: Import prepared!
  Prepared couchbase-capella_app_service_private_endpoint for import
couchbase-capella_app_service_private_endpoint.accept_endpoint: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

## DELETE

Command: `terraform destroy`

Destroying the resource removes the private endpoint from the App Service. If the endpoint was linked, the connection is terminated. If it was still pending, the request is rejected.

Sample Output:
```
terraform destroy

Terraform will perform the following actions:

  # couchbase-capella_app_service_private_endpoint.accept_endpoint will be destroyed
  - resource "couchbase-capella_app_service_private_endpoint" "accept_endpoint" {
      - app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - endpoint_id     = "vpce-000000000000aaaaa" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - service_name    = "com.amazonaws.vpce.us-east-1.vpce-svc-11223344556677aab" -> null
      - status          = "linked" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_service_private_endpoint.accept_endpoint: Destroying...
couchbase-capella_app_service_private_endpoint.accept_endpoint: Destruction complete after 2s

Destroy complete! Resources: 1 destroyed.
```
