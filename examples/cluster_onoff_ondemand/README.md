# Capella Cluster On/Off Example

This example shows how to switch the cluster to on or off state on demand in Capella.

This switches the selected Capella cluster to on or off state. It uses the organization ID, project ID, and cluster ID to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Switch the cluster to on/off state as mentioned in the `create_cluster_onoff_ondemand.tf` file.
2. UPDATE: Update the current cluster state to on/off.
3. IMPORT: Import a cluster state that exists in Capella but not in the terraform state file.

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

  # couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand will be created
  + resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + state                      = "on"
      + turn_on_linked_app_service = true
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_cluster_onoff_ondemand = {
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + state                      = "on"
      + turn_on_linked_app_service = true
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

```

### Apply the Plan, in order to switch the cluster to on/off state

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

  # couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand will be created
  + resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + state                      = "on"
      + turn_on_linked_app_service = true
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_cluster_onoff_ondemand = {
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + state                      = "on"
      + turn_on_linked_app_service = true
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Creating...
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Creation complete after 1s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_cluster_onoff_ondemand = {
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "state" = "on"
  "turn_on_linked_app_service" = true
}

```
### Note the output of the new Cluster state
Command: `terraform output new_cluster_onoff_ondemand`

Sample Output:
```
$ terraform output new_cluster_onoff_ondemand
{
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "state" = "on"
  "turn_on_linked_app_service" = true
}

```

## IMPORT
### Remove the resource `new_cluster_onoff_ondemand` from the Terraform State file

Command: `terraform state rm couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand`

Sample Output:
``` 
$  terraform state rm couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand
Removed couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand
Successfully removed 1 resource instance(s).

```
Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
``` 
$ terraform import couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Importing from ID "cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Import prepared!
  Prepared couchbase-capella_cluster_onoff_ondemand for import
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

Here, we pass the IDs as a single comma-separated string.
The first ID is the cluster ID i.e. the ID of the cluster to which state belongs.
The second ID is the project ID i.e. the ID of the project to which the cluster belongs.
The third ID is the organization ID i.e. the ID of the organization to which the project belongs.


## UPDATE
### Let us edit terraform.tfvars file to change the cluster state.

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
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Refreshing state...

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand has changed
  ~ resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
      ~ state                      = "<null>" -> "\"<null>\""
        # (4 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand will be updated in-place
  ~ resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
      ~ state                      = "\"<null>\"" -> "off"
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_cluster_onoff_ondemand = {
      ~ state                      = "<null>" -> "off"
        # (4 unchanged attributes hidden)
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────


```

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
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Refreshing state...

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

# couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand has changed
~ resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
~ state                      = "<null>" -> "\"<null>\""
# (4 unchanged attributes hidden)
}


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
~ update in-place

Terraform will perform the following actions:

# couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand will be updated in-place
~ resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
~ state                      = "\"<null>\"" -> "off"
# (4 unchanged attributes hidden)
}

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
~ new_cluster_onoff_ondemand = {
~ state                      = "<null>" -> "off"
# (4 unchanged attributes hidden)
}

Do you want to perform these actions?
Terraform will perform the actions described above.
Only 'yes' will be accepted to approve.

Enter a value: yes

couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Modifying...
couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand: Modifications complete after 1s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_cluster_onoff_ondemand = {
"cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
"organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
"project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
"state" = "off"
"turn_on_linked_app_service" = false
}

```
