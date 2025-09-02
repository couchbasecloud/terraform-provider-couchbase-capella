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

Command: `terraform plan`

Sample Output:
```
$ terraform plan                                                                                                                                                                                          
```

### View the current list of resources that are present in Terraform State

Command: `terraform show`

Sample Output:
```
 $ terraform show

```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
 $ terraform state list                                                                

```


### 2. Update the Import Filter

Make changes to the `access_control_function` attribute in `create_access_control_function.tf` and run:

Command: `terraform apply`

Sample Output:
```
 $ terraform apply                     1 ↵

```


## IMPORT
### Remove the resource `example_access_function` from the Terraform State file

Command: `terraform state rm couchbase-capella_access_control_function.example_access_function`

Sample Output:
``` 
terraform state rm couchbase-capella_access_control_function.examp
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_access_control_function.example_access_function id=<appservice_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_access_control_function.example_access_function id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 
terraform import 
``` 

### 5. Delete the resources that Terraform manages

Command: `terraform destroy`

Sample Output:
```
 $ terraform destroy

```


## Prerequisites

- Couchbase Capella organization with appropriate permissions
- Existing App Service and App Endpoint
- Valid scope and collection within the App Endpoint

## Important Notes

- Import Filters are applied at the collection level
- Functions must be valid JavaScript
- Changes to Import Filters affect all documents in the collection
- The resource requires `app_endpoint_name` rather than ID for identification 