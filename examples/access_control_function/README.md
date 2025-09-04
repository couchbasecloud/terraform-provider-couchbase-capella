# Capella App Endpoint Access Control Function Example

This example shows how to create and manage App Endpoint Access Control Functions in Capella.

This creates a new access control function for a specific collection within an App Endpoint in the selected Capella organization. Access control functions are JavaScript functions that specify access control policies applied to documents in collections. Every document update is processed by this function.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new access control function in Capella as stated in the `create_access_control_function.tf` file.
2. UPDATE: Update the access control function using Terraform.
3. DELETE: Delete the newly created access control function from Capella.
4. IMPORT: Import an access control function that exists in Capella but not in the terraform state file.

If you check the `terraform.template.tfvars` file - you can see that we need several variables to run the terraform commands.
Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_access_control_function.acf will be created
  + resource "couchbase-capella_app_endpoint_access_control_function" "acf" {
      + access_control_function = "function (doc, oldDoc, meta) {channel(doc.channels); }"
      + app_endpoint            = "test-endpoint"
      + app_service_id          = "{appServiceId}"
      + cluster_id              = "{clusterId}"
      + collection              = "test"
      + organization_id         = "{orgId}"
      + project_id              = "{projectId}"
      + scope                   = "test"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_access_control_function.acf: Creating...
couchbase-capella_app_endpoint_access_control_function.acf: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## UPDATE

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_app_endpoint_access_control_function.acf: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_access_control_function.acf will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_access_control_function" "acf" {
      ~ access_control_function = "function (doc, oldDoc, meta) {channel(doc.channels); }" -> "function (currentDoc, oldDoc, meta) {channel(doc.channels); }"
        # (7 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_access_control_function.acf: Modifying...
couchbase-capella_app_endpoint_access_control_function.acf: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_app_endpoint_access_control_function.acf: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_access_control_function.acf will be destroyed
  - resource "couchbase-capella_app_endpoint_access_control_function" "acf" {
      - access_control_function = "function (currentDoc, oldDoc, meta) {channel(doc.channels); }" -> null
      - app_endpoint            = "test-endpoint" -> null
      - app_service_id          = "{appServiceId}" -> null
      - cluster_id              = "{clusterId}" -> null
      - collection              = "test" -> null
      - organization_id         = "{orgId}" -> null
      - project_id              = "{projectId}" -> null
      - scope                   = "test" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_endpoint_access_control_function.acf: Destroying...
couchbase-capella_app_endpoint_access_control_function.acf: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```


## IMPORT

Command: `terraform import couchbase-capella_app_endpoint_access_control_function.acf organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>,app_endpoint_name=<app_endpoint_name>,scope_name=<scope_name>,collection_name=<collection_name>`

Sample Output:
```
$ terraform import couchbase-capella_app_endpoint_access_control_function.acf organization_id={orgId},project_id={projectId},cluster_id={clusterId},app_service_id={appServiceId},app_endpoint=test-endpoint,scope_name=test,collection_name=test
couchbase-capella_app_endpoint_access_control_function.acf: Importing from ID "organization_id={orgId},project_id={projectId},cluster_id={clusterId},app_service_id={appServiceId},app_endpoint=test-endpoint,scope_name=test,collection_name=test"...
couchbase-capella_app_endpoint_access_control_function.acf: Import prepared!
  Prepared couchbase-capella_app_endpoint_access_control_function for import
couchbase-capella_app_endpoint_access_control_function.acf: Refreshing state...
2025-09-02T13:55:15.576-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_app_endpoint_access_control_function.acf during refresh.
      - .access_control_function: was null, but now cty.StringVal("function (currentDoc, oldDoc, meta) {channel(doc.channels); }")
      - .app_endpoint: was cty.StringVal("organization_id={orgId},project_id={projectId},cluster_id={clusterId},app_service_id={appServiceId},app_endpoint=test-endpoint,scope_name=test,collection_name=test"), but now cty.StringVal("test-endpoint")
      - .app_service_id: was null, but now cty.StringVal("{appServiceId}")
      - .cluster_id: was null, but now cty.StringVal("{clusterId}")
      - .collection: was null, but now cty.StringVal("test")
      - .organization_id: was null, but now cty.StringVal("{orgId}")
      - .project_id: was null, but now cty.StringVal("{projectId}")
      - .scope: was null, but now cty.StringVal("test")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```
