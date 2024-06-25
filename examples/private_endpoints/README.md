# Capella Private Endpoint Example

This example shows how to accept and reject a private endpoint in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Accept a private endpoint as shown in `accept_endpoint.tf` file.
2. LIST: List private endpoints as shown in the `list_endpoints.tf` file.
3. DELETE: Reject a private endpoint.
4. IMPORT: Import private endpoint status into the state file.

If you check the `terraform.template.tfvars` file - you can see that we need 3 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE AND LIST

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_private_endpoints.list_endpoints: Reading...
data.couchbase-capella_private_endpoints.list_endpoints: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_private_endpoints.accept_endpoint will be created
  + resource "couchbase-capella_private_endpoints" "accept_endpoint" {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + endpoint_id     = "vpce-7"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + status          = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + list_endpoints = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + data            = [
          + {
              + id     = "vpce-1"
              + status = "rejected"
            },
          + {
              + id     = "vpce-2"
              + status = "rejected"
            },
          + {
              + id     = "vpce-3"
              + status = "rejected"
            },
          + {
              + id     = "vpce-4"
              + status = "rejected"
            },
          + {
              + id     = "vpce-5"
              + status = "rejected"
            },
          + {
              + id     = "vpce-6"
              + status = "rejected"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_private_endpoints.accept_endpoint: Creating...
couchbase-capella_private_endpoints.accept_endpoint: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

list_endpoints = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "data" = tolist([
    {
      "id" = "vpce-1"
      "status" = "rejected"
    },
    {
      "id" = "vpce-2"
      "status" = "rejected"
    },
    {
      "id" = "vpce-3"
      "status" = "rejected"
    },
    {
      "id" = "vpce-4"
      "status" = "rejected"
    },
    {
      "id" = "vpce-5"
      "status" = "rejected"
    },
    {
      "id" = "vpce-6"
      "status" = "rejected"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

## DELETE

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_private_endpoints.accept_endpoint: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_private_endpoints.accept_endpoint will be destroyed
  # (because couchbase-capella_private_endpoints.accept_endpoint is not in configuration)
  - resource "couchbase-capella_private_endpoints" "accept_endpoint" {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - endpoint_id     = "vpce-7" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - status          = "linked" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - list_endpoints = {
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - data            = [
          - {
              - id     = "vpce-1"
              - status = "rejected"
            },
          - {
              - id     = "vpce-2"
              - status = "rejected"
            },
          - {
              - id     = "vpce-3"
              - status = "rejected"
            },
          - {
              - id     = "vpce-4"
              - status = "rejected"
            },
          - {
              - id     = "vpce-5"
              - status = "rejected"
            },
          - {
              - id     = "6"
              - status = "rejected"
            },
        ]
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_private_endpoints.accept_endpoint: Destroying...
couchbase-capella_private_endpoints.accept_endpoint: Destruction complete after 1s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.
```

## IMPORT

Command: `terraform import couchbase-capella_private_endpoint_service.new_service endpoint_id=<endpoint_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
terraform import couchbase-capella_private_endpoints.accept_endpoint endpoint_id=vpce-9,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_private_endpoints.accept_endpoint: Importing from ID "endpoint_id=vpce-9,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_private_endpoints.accept_endpoint: Import prepared!
Prepared couchbase-capella_private_endpoints for import
couchbase-capella_private_endpoints.accept_endpoint: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```