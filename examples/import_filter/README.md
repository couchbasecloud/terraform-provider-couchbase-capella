# Capella Import Filter Example

This example shows how to create and manage import filters for App Endpoints in Capella.

This creates a new Import Filter for a specific collection within an App Endpoint. Import Filters are 

The default Import Filter is `function(doc){channel(doc.channels);}` for the default collection and `function(doc){channel(collectionName);}` for named collections.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

1. CREATE: Create a new Import Filter in an existing Capella App Endpoint as stated in the `create_import_filter.tf` file.
2. UPDATE: Update the app service configuration using Terraform.
3. LIST: List existing app services in Capella as stated in the `list_access_functions.tf` file.
4. IMPORT: Import an app services that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created app service from Capella.

## Example Walkthrough

### 1. CREATE: View the plan for the resources that Terraform will create

#### View the plan
Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_import_filter.example_import_filter will be created
  + resource "couchbase-capella_import_filter" "example_import_filter" {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + import_filter   = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
      + app_endpoint_name        = "api"
      + scope           = "_default"
      + collection      = "_default"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```
#### Apply the terraform config to create the resource
Command: `terraform apply`

Sample Output:
```
 $ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_import_filter.example_import_filter will be created
  + resource "couchbase-capella_import_filter" "example_import_filter" {
      + app_service_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id         = "ffffffff-aaaa-1414-eeee-000000000000"
      + import_filter      = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
      + app_endpoint_name  = "api"
      + scope              = "_default"
      + collection         = "_default"
      + organization_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id         = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_import_filter.example_import_filter: Creating...
couchbase-capella_import_filter.example_import_filter: Still creating... [00m10s elapsed]
couchbase-capella_import_filter.example_import_filter: Creation complete after 14s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

### View the current list of resources that are present in Terraform State

Command: `terraform show`

Sample Output:
```
$ terraform show
# couchbase-capella_import_filter.example_import_filter:
resource "couchbase-capella_import_filter" "example_import_filter" {
    app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
    cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    import_filter     = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }"
    app_endpoint_name = "api"
    scope             = "_default"
    collection        = "_default"
    organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
    project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
}
```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
couchbase-capella_import_filter.example_import_filter
```


### 2. Update the Import Filter

Make changes to the `import_filter` attribute in `create_import_filter.tf` and run:

Command: `terraform apply`

Sample Output:
```
$ terraform apply  


couchbase-capella_import_filter.example_import_filter: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_import_filter.example_import_filter will be updated in-place
  ~ resource "couchbase-capella_import_filter" "example_import_filter" {
      ~ import_filter   = "function(doc) { if (doc.type != 'mobile') { return false; } return true; }" -> "function(doc) { if (doc.type != 'edge') { return false; } return true; }"
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_import_filter.example_import_filter: Modifying...
couchbase-capella_import_filter.example_import_filter: Modifications complete after 7s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```


## IMPORT
### Remove the resource `example_access_function` from the Terraform State file

Command: `terraform state rm couchbase-capella_import_filter.example_import_filter`

Sample Output:
``` 
$ terraform state rm couchbase-capella_import_filter.example_import_filter
Removed couchbase-capella_import_filter.example_import_filter
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import  couchbase-capella_import_filter.coll2_import_filter app_endpoint_name=api,scope_name=_default,collection_name=_default,app_service_id=<appservice_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import  couchbase-capella_import_filter.coll2_import_filter app_endpoint_name=api,scope_name=_default,collection_name=_default,id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 

$ terraform import  couchbase-capella_import_filter.example_import_filterterraform import  couchbase-capella_import_filter.coll2_import_filter app_endpoint_name=api,scope_name=_default,collection_name=_default,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_import_filter.example_import_filter: Importing from ID "app_endpoint_name=api,scope_name=_default,collection_name=_default,organization_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_import_filter.example_import_filter: Import prepared!
  Prepared couchbase-capella_import_filter for import
couchbase-capella_import_filter.example_import_filter: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

``` 

### 5. Delete the resources that Terraform manages

Command: `terraform destroy`

Sample Output:
```
 $ terraform destroy

couchbase-capella_import_filter.example_import_filter: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_import_filter.example_import_filter will be destroyed
  - resource "couchbase-capella_import_filter" "example_import_filter" {
      - app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - import_filter     = "function(doc) { if (doc.type != 'edge') { return false; } return true; }" -> null
      - app_endpoint_name = "api" -> null
      - scope             = "_default" -> null
      - collection        = "_default" -> null
      - organization_id   = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id        = "ffffffff-aaaa-1414-eeee-000000000000" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_import_filter.example_import_filter: Destroying...
couchbase-capella_import_filter.example_import_filter: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```

