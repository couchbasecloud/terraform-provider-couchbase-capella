# Capella App Service Log Streaming Activation Status Resource Example

This example shows how to control the activation status of App Services Log Streaming.

This controls whether Log Streaming is paused or enabled (resume Log Streaming if it is paused) for a Capella App Service. Log Streaming must have already been set up on the App Service for this resource to be used to control the activation status. It uses the organization ID, project ID, cluster ID and App Service ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create the App Service  as stated in the `create_app_service_log_streaming_activation_status.tf` file.
2. UPDATE: Update the App Service Log Streaming activation status using Terraform.
3. IMPORT: Import App Service Log Streaming so the activation status can be controlled with Terraform.
4. DELETE: No-op remove the App Service Log Streaming activation status from the Terraform state file.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
This will create the App Service Log Streaming activation status resource in Terraform. If the `state` provided matches the actual state of the log streaming config state in Capella, then there will be no changes to be made and the resource will be added (similar to an import). However, if the `state` provided does not match the state, then Terraform will pause or enable log streaming depend on the `state` provided in Terraform.

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status will be created
  + resource "couchbase-capella_app_service_log_streaming_activation_status" "app_service_log_streaming_activation_status" {
      + app_service_id  = "9a34f4a3-b05e-4e6e-a249-c472f9f915dd"
      + cluster_id      = "669cb38a-ea74-42df-94fe-bdca9d332e1c"
      + organization_id = "a16247ef-8356-4a09-9edb-7fd2507a2881"
      + project_id      = "e5dbe56f-43bb-482b-a161-80bb7430dbbd"
      + state           = "enabled"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Apply the Plan, in order to create a new entry for the activation status

Command: `terraform apply`

Sample Output:
```
$ terraform apply
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status will be created
  + resource "couchbase-capella_app_service_log_streaming_activation_status" "app_service_log_streaming_activation_status" {
      + app_service_id  = "9a34f4a3-b05e-4e6e-a249-c472f9f915dd"
      + cluster_id      = "669cb38a-ea74-42df-94fe-bdca9d332e1c"
      + organization_id = "a16247ef-8356-4a09-9edb-7fd2507a2881"
      + project_id      = "e5dbe56f-43bb-482b-a161-80bb7430dbbd"
      + state           = "enabled"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Creating...
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still creating... [00m10s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still creating... [00m20s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still creating... [00m30s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still creating... [00m40s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still creating... [00m50s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Creation complete after 56s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

### List the resources that are present in the Terraform State file

Command: `terraform state list`

Sample Output:
``` 
$ terraform state list
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status
```

## IMPORT
### Remove the resource `app_service_log_streaming_activation_status` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status`

Sample Output:
``` 
$ terraform state rm couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status

```
Please note, this command will only remove the resource from the Terraform State file, but in reality, Log Streaming is still enabled for the App Service in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status app_service_id=9a34f4a3-b05e-4e6e-a249-c472f9f915dd,cluster_id=669cb38a-ea74-42df-94fe-bdca9d332e1c,project_id=e5dbe56f-43bb-482b-a161-80bb7430dbbd,organization_id=a16247ef-8356-4a09-9edb-7fd2507a2881`

Sample Output:
``` 
$ terraform import couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status app_service_id=9a34f4a3-b05e-4e6e-a249-c472f9f915dd,cluster_id=669cb38a-ea74-42df-94fe-bdca9d332e1c,project_id=e5dbe56f-43bb-482b-a161-80bb7430dbbd,organization_id=a16247ef-8356-4a09-9edb-7fd2507a2881
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Importing from ID "app_service_id=9a34f4a3-b05e-4e6e-a249-c472f9f915dd,cluster_id=669cb38a-ea74-42df-94fe-bdca9d332e1c,project_id=e5dbe56f-43bb-482b-a161-80bb7430dbbd,organization_id=a16247ef-8356-4a09-9edb-7fd2507a2881"...
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Import prepared!
  Prepared couchbase-capella_app_service_log_streaming_activation_status for import
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Refreshing state...

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
### Let us edit the terraform.tfvars file to change the Log Streaming state

Command: `terraform apply -var-file=terraform.template.tfvars`

``` 
$ terraform apply -var-file=terraform.template.tfvars
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status will be updated in-place
  ~ resource "couchbase-capella_app_service_log_streaming_activation_status" "app_service_log_streaming_activation_status" {
      ~ state           = "enabled" -> "paused"
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Modifying...
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still modifying... [00m10s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still modifying... [00m20s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still modifying... [00m30s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still modifying... [00m40s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still modifying... [00m50s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Still modifying... [01m00s elapsed]
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Modifications complete after 1m2s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## DESTROY
This will result in a no-op destroy of the log streaming activation status resource in Terraform, with no actual change to the log streaming state.
### Finally, delete the log streaming activation status resource and destroy all resources managed by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status will be destroyed
  - resource "couchbase-capella_app_service_log_streaming_activation_status" "app_service_log_streaming_activation_status" {
      - app_service_id  = "9a34f4a3-b05e-4e6e-a249-c472f9f915dd" -> null
      - cluster_id      = "669cb38a-ea74-42df-94fe-bdca9d332e1c" -> null
      - organization_id = "a16247ef-8356-4a09-9edb-7fd2507a2881" -> null
      - project_id      = "e5dbe56f-43bb-482b-a161-80bb7430dbbd" -> null
      - state           = "paused" -> null
    }
    
Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes
  
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Destroying...
couchbase-capella_app_service_log_streaming_activation_status.app_service_log_streaming_activation_status: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```
