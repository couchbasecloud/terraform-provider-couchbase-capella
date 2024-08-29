# Capella Private Endpoint Service Example

This example shows how to enable/disable Private Endpoint Service in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Enable private endpoint service on the cluster as shown in `enable_service.tf` file.
2. LIST: Get private endpoint service status as shown in the `get_status.tf` file.
3. DELETE: Disable private endpoint service on the cluster.
4. IMPORT: Import private endpoint status into the state file.

If you check the `terraform.template.tfvars` file - you can see that we need 3 main variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_private_endpoint_service.new_service will be created
  + resource "couchbase-capella_private_endpoint_service" "new_service" {
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

couchbase-capella_private_endpoint_service.new_service: Creating...
couchbase-capella_private_endpoint_service.new_service: Still creating... [10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [1m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [2m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [3m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [4m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [5m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [6m50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still creating... [7m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Creation complete after 7m0s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## LIST

Command: `terraform plan`

Sample Output:
```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵
data.couchbase-capella_private_endpoint_service.service_status: Reading...
data.couchbase-capella_private_endpoint_service.service_status: Read complete after 0s

Changes to Outputs:
  + existing_auditlogexport = {
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + enabled         = true
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

```

## DELETE
### Remove the resource block in enable_service.tf

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
│ incompatible with published releases.
╵
couchbase-capella_private_endpoint_service.new_service: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_private_endpoint_service.new_service will be destroyed
  # (because couchbase-capella_private_endpoint_service.new_service is not in configuration)
  - resource "couchbase-capella_private_endpoint_service" "new_service" {
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

couchbase-capella_private_endpoint_service.new_service: Destroying...
couchbase-capella_private_endpoint_service.new_service: Still destroying... [10s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still destroying... [20s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still destroying... [30s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still destroying... [40s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still destroying... [50s elapsed]
couchbase-capella_private_endpoint_service.new_service: Still destroying... [1m0s elapsed]
couchbase-capella_private_endpoint_service.new_service: Destruction complete after 1m0s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.
```

## IMPORT
### Ensure a resource block is configured as shown in enable_service.tf

Command: `terraform import couchbase-capella_private_endpoint_service.new_service cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
terraform import couchbase-capella_private_endpoint_service.new_service \
cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_private_endpoint_service.new_service: Importing from ID "cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_private_endpoint_service.new_service: Import prepared!
  Prepared couchbase-capella_private_endpoint_service for import
couchbase-capella_private_endpoint_service.new_service: Refreshing state...
2024-06-03T21:08:59.649-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_private_endpoint_service.new_service during refresh.
      - .enabled: was null, but now cty.False
      - .organization_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .project_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .cluster_id: was cty.StringVal("cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"), but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```
