# Capella App Service On/Off Example

This example shows how to switch the app service to on or off state on demand in Capella.

This switches the selected Capella App service to on or off state. It uses the organization ID, project ID, cluster ID and app service ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Switch the app service to on/off state as mentioned in the `create_app_service_onoff_ondemand.tf` file.
2. UPDATE: Update the current app service state to on/off.
3. IMPORT: Import an app service state that exists in Capella but not in the terraform state file.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand will be created
  + resource "couchbase-capella_app_service_onoff_ondemand" "new_app_service_onoff_ondemand" {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + state           = "on"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service_onoff_ondemand = {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + state           = "on"
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────


```


### Apply the Plan, in order to switch the app service to on/off state

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand will be created
  + resource "couchbase-capella_app_service_onoff_ondemand" "new_app_service_onoff_ondemand" {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + state           = "on"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_app_service_onoff_ondemand = {
      + app_service_id  = "ffffffff-aaaa-1414-eeee-000000000000"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + state           = "on"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Creating...
couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_app_service_onoff_ondemand = {
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "state" = "on"
}

```

## IMPORT
### Remove the resource `new_app_service_onoff_ondemand` from the Terraform State file

Command: `terraform state rm couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand`

Sample Output:
``` 
$  terraform state rm couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand
Removed couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand
Successfully removed 1 resource instance(s).
```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand app_service_id=<app_service_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand cluster_id=07d9b739-0ede-4f1b-859b-1b21d0647fc6,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 
$ terraform import couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Importing from ID "app_service_id=ffffffff-aaaa-1414-eeee-000000000000,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Import prepared!
  Prepared couchbase-capella_app_service_onoff_ondemand for import
couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the app service ID i.e. the ID of the app service that we want to switch on/off to.
The second ID is the cluster ID i.e. the ID of the cluster to which app service belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

## UPDATE
### Let us edit terraform.tfvars file to change the app service state.

command: `terrafom apply`

Sample Output:

```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Refreshing state...

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand has changed
  ~ resource "couchbase-capella_app_service_onoff_ondemand" "new_app_service_onoff_ondemand" {
      ~ state           = "<null>" -> "\"<null>\""
        # (4 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand will be updated in-place
  ~ resource "couchbase-capella_app_service_onoff_ondemand" "new_app_service_onoff_ondemand" {
      ~ state           = "\"<null>\"" -> "off"
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_app_service_onoff_ondemand = {
      ~ state           = "<null>" -> "off"
        # (4 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Modifying...
couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_app_service_onoff_ondemand = {
  "app_service_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "state" = "off"
}

```
