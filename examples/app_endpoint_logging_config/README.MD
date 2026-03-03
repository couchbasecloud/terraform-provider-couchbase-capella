# Capella App Endpoint Logging Config Example

This example shows how to manage an App Endpoint Logging Config for an App Service in Capella. This logging config controls what logs to stream to the log collector set up on the App Service.

The logging config can only be managed if the App Service has Log Streaming enabled. See the `app_service_log_streaming` resource for how to manage this through Terraform.

This modifies the logging config for the selected App Endpoint. It uses the organization ID, project ID, cluster ID, App Service ID, and App Endpoint name to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a logging config for a Capella app endpoint as stated in the `create_app_endpoint_logging_config.tf` file.
2. UPDATE: Update the logging config using Terraform.
3. IMPORT: Import a logging config that exists in Capella but not in the Terraform state file.
4. DELETE: Remove the logging config from the Terraform state file.

If you wish to use the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & READ
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config will be created
  + resource "couchbase-capella_app_endpoint_log_streaming_config" "new_app_endpoint_log_streaming_config" {
      + app_endpoint_name = "test_app_endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + log_keys          = [
          + "HTTP",
          + "Import",
        ]
      + log_level         = "info"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_endpoint_log_streaming_config = {
      + app_endpoint_name = "test_app_endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + log_keys          = [
          + "HTTP",
          + "Import",
        ]
      + log_level         = "info"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }
```

### Apply the Plan, in order to create the App Endpoint Logging Config Resource

Command: `terraform apply`

Sample Output:
```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config will be created
  + resource "couchbase-capella_app_endpoint_log_streaming_config" "new_app_endpoint_log_streaming_config" {
      + app_endpoint_name = "test_app_endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + log_keys          = [
          + "HTTP",
          + "Import",
        ]
      + log_level         = "info"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_endpoint_log_streaming_config = {
      + app_endpoint_name = "test_app_endpoint"
      + app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      + log_keys          = [
          + "HTTP",
          + "Import",
        ]
      + log_level         = "info"
      + organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Creating...
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_endpoint_log_streaming_config = {
  "app_endpoint_name" = "test_app_endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "log_keys" = toset([
    "HTTP",
    "Import",
  ])
  "log_level" = "info"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
``` 
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config
```

## IMPORT
### Remove the resource `new_app_endpoint_log_streaming_config` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config`

Sample Output:
``` 
Removed couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the logging config modifications still exist in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config app_endpoint_name=<app_endpoint_name>,app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config app_endpoint_name=test_app_endpoint,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Importing from ID "app_endpoint_name=test_app_endpoint,app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Import prepared!
  Prepared couchbase-capella_app_endpoint_log_streaming_config for import
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the app endpoint name i.e. the name of the app endpoint to which the logging config belongs.
The second ID is the app service ID i.e. the ID of the app service to which the app endpoint belongs.
The third ID is the cluster ID i.e. the ID of the cluster to which the app service belongs.
The fourth ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fifth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit terraform.tfvars file to change the app endpoint logging config configuration settings.

Command: `terraform plan`

Sample Output:

``` 
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_log_streaming_config" "new_app_endpoint_log_streaming_config" {
      ~ log_keys          = [
          - "HTTP",
          - "Import",
          + "Auth",
        ]
      ~ log_level         = "info" -> "warn"
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_app_endpoint_log_streaming_config = {
      ~ log_keys          = [
          - "HTTP",
          - "Import",
          + "Auth",
        ]
      ~ log_level         = "info" -> "warn"
        # (5 unchanged attributes hidden)
    }
```

command: `terrafom apply`

Sample Output:

```
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config will be updated in-place
  ~ resource "couchbase-capella_app_endpoint_log_streaming_config" "new_app_endpoint_log_streaming_config" {
      ~ log_keys          = [
          - "HTTP",
          - "Import",
          + "Auth",
        ]
      ~ log_level         = "info" -> "warn"
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_app_endpoint_log_streaming_config = {
      ~ log_keys          = [
          - "HTTP",
          - "Import",
          + "Auth",
        ]
      ~ log_level         = "info" -> "warn"
        # (5 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Modifying...
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_app_endpoint_log_streaming_config = {
  "app_endpoint_name" = "test_app_endpoint"
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "log_keys" = toset([
    "Auth",
  ])
  "log_level" = "warn"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
}
```

## DESTROY

Note that destroying this resource results in a simple no-op removal of this resource from the state file. The App Endpoint will still retain the logging config modifications that have been made.

### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
Changes to Outputs:
  - new_app_endpoint_log_streaming_config = {
      - app_endpoint_name = "test_app_endpoint"
      - app_service_id    = "ffffffff-aaaa-1414-eeee-000000000000"
      - cluster_id        = "ffffffff-aaaa-1414-eeee-000000000000"
      - log_keys          = [
          - "Auth",
        ]
      - log_level         = "warn"
      - organization_id   = "ffffffff-aaaa-1414-eeee-000000000000"
      - project_id        = "ffffffff-aaaa-1414-eeee-000000000000"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Destroying...
couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```