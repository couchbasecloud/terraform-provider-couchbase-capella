# Capella App Service Log Streaming Resource Example

This example shows how to setup and manage App Services Log Streaming.

This enables Log Streaming for a Capella App Service. It uses the organization ID, project ID, cluster ID and App Service ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Setup Log Streaming for an App Service as stated in the `create_app_service_log_streaming.tf` file.
2. UPDATE: Update the App Service Log Streaming configuration using Terraform.
3. IMPORT: Import Log Streaming that has already been set up for an App Service but is not in the Terraform state file.
4. DELETE: Disable Log Streaming for an App Service.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming.app_service_log_streaming will be created
  + resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
      + app_service_id  = "8f3fda40-56a2-4105-9eab-c4df0f3c13e3"
      + cluster_id      = "58c74f93-5ae6-4167-9904-741903180ec7"
      + config_state    = (known after apply)
      + credentials     = (sensitive value)
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + output_type     = "generic_http"
      + project_id      = "93aa3b7e-b5a0-4ce6-bb9d-426f1e928c49"
      + streaming_state = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new App Service entry

Command: `terraform apply`

Sample Output:
```
$ terraform apply
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming.app_service_log_streaming will be created
  + resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
      + app_service_id  = "8f3fda40-56a2-4105-9eab-c4df0f3c13e3"
      + cluster_id      = "58c74f93-5ae6-4167-9904-741903180ec7"
      + config_state    = (known after apply)
      + credentials     = (sensitive value)
      + organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76"
      + output_type     = "generic_http"
      + project_id      = "93aa3b7e-b5a0-4ce6-bb9d-426f1e928c49"
      + streaming_state = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_log_streaming.app_service_log_streaming: Creating...
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still creating... [00m10s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still creating... [00m20s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still creating... [00m30s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still creating... [00m40s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still creating... [00m50s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Creation complete after 58s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

### List the resources that are present in the Terraform State file

Command: `terraform state list`

Sample Output:
``` 
$ terraform state list
couchbase-capella_app_service_log_streaming.app_service_log_streaming
```

## IMPORT
### Remove the resource `app_service_log_streaming` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_service_log_streaming.app_service_log_streaming`

Sample Output:
``` 
$ terraform state rm couchbase-capella_app_service_log_streaming.app_service_log_streaming
Removed couchbase-capella_app_service_log_streaming.app_service_log_streaming
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, Log Streaming is still enabled for the App Service in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_service_log_streaming.app_service_log_streaming app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_app_service_log_streaming.app_service_log_streaming app_service_id=8f3fda40-56a2-4105-9eab-c4df0f3c13e3,cluster_id=58c74f93-5ae6-4167-9904-741903180ec7,project_id=93aa3b7e-b5a0-4ce6-bb9d-426f1e928c49,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76`

Sample Output:
``` 
$ terraform import couchbase-capella_app_service_log_streaming.app_service_log_streaming app_service_id=8f3fda40-56a2-4105-9eab-c4df0f3c13e3,cluster_id=58c74f93-5ae6-4167-9904-741903180ec7,project_id=93aa3b7e-b5a0-4ce6-bb9d-426f1e928c49,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Importing from ID "app_service_id=8f3fda40-56a2-4105-9eab-c4df0f3c13e3,cluster_id=58c74f93-5ae6-4167-9904-741903180ec7,project_id=93aa3b7e-b5a0-4ce6-bb9d-426f1e928c49,organization_id=adb4fb4c-1d98-4287-ac33-230742d2cc76"...
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Import prepared!
  Prepared couchbase-capella_app_service_log_streaming for import
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the app service ID i.e. the ID of the resource that we want to set up Log Streaming for.
The second ID is the cluster ID i.e. the ID of the cluster to which App Service belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit the terraform.tfvars file to change the Log Streaming credentials

Command: `terraform apply -var-file=terraform.template.tfvars`

``` 
$ terraform apply -var-file=terraform.template.tfvars
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming.app_service_log_streaming will be updated in-place
  ~ resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
      ~ config_state    = "enabled" -> (known after apply)
      + credentials     = (sensitive value)
      ~ streaming_state = "unknown" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_log_streaming.app_service_log_streaming: Modifying...
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still modifying... [00m10s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still modifying... [00m20s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still modifying... [00m30s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Modifications complete after 37s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DESTROY
### Finally, disable Log Streaming and destroy all resources managed by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming.app_service_log_streaming will be destroyed
  - resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
      - app_service_id  = "8f3fda40-56a2-4105-9eab-c4df0f3c13e3" -> null
      - cluster_id      = "58c74f93-5ae6-4167-9904-741903180ec7" -> null
      - config_state    = "enabled" -> null
      - credentials     = (sensitive value) -> null
      - organization_id = "adb4fb4c-1d98-4287-ac33-230742d2cc76" -> null
      - output_type     = "generic_http" -> null
      - project_id      = "93aa3b7e-b5a0-4ce6-bb9d-426f1e928c49" -> null
      - streaming_state = "unknown" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_app_service_log_streaming.app_service_log_streaming: Destroying...
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Still destroying... [00m10s elapsed]
couchbase-capella_app_service_log_streaming.app_service_log_streaming: Destruction complete after 19s

Destroy complete! Resources: 1 destroyed. 
```
