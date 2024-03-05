# Capella Audit Log Settings Example

This example shows how to create and manage audit log settings in Capella.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1.  LIST AND SAVE STATE: List existing auditlogsettings in Capella as stated in the `get_auditlogsettings.tf` file.
2.  CREATE: Create a new auditlogsettings in Capella as stated in the `create_auditlogsettings.tf` file.


If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## LIST

Commands: `terraform plan`
          `terraform apply`

``` 
terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/me/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state
│ to become incompatible with published releases.
╵
data.couchbase-capella_auditlogsettings.existing_auditlogsettings: Reading...
data.couchbase-capella_auditlogsettings.existing_auditlogsettings: Read complete after 0s

Changes to Outputs:
  + existing_auditlogsettings = {
      + auditenabled    = false
      + cluster_id      = "0fad8c6d-1a42-44b6-ada8-3c7a917114a5"
      + disabledusers   = []
      + enabledeventids = []
      + organization_id = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
      + project_id      = "6b22a328-3ad8-41c3-8d83-e6242d170423"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real
infrastructure.

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions
if you run "terraform apply" now.

terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/me/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state
│ to become incompatible with published releases.
╵
data.couchbase-capella_auditlogsettings.existing_auditlogsettings: Reading...
data.couchbase-capella_auditlogsettings.existing_auditlogsettings: Read complete after 0s

Changes to Outputs:
  + existing_auditlogsettings = {
      + auditenabled    = false
      + cluster_id      = "0fad8c6d-1a42-44b6-ada8-3c7a917114a5"
      + disabledusers   = []
      + enabledeventids = []
      + organization_id = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
      + project_id      = "6b22a328-3ad8-41c3-8d83-e6242d170423"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real
infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogsettings = {
  "auditenabled" = false
  "cluster_id" = "0fad8c6d-1a42-44b6-ada8-3c7a917114a5"
  "disabledusers" = tolist([])
  "enabledeventids" = tolist([])
  "organization_id" = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
  "project_id" = "6b22a328-3ad8-41c3-8d83-e6242d170423"
}
Outputs:

existing_auditlogsettings = {
  "auditenabled" = false
  "cluster_id" = "4575d05b-cc8f-4d85-b3e8-291d827ec812"
  "disabledusers" = tolist([])
  "enabledeventids" = tolist([])
  "organization_id" = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
  "project_id" = "6b22a328-3ad8-41c3-8d83-e6242d170423"
}
```

## CREATE

Command: `terraform apply`

``` 
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/me/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.couchbase-capella_auditlogsettings.existing_auditlogsettings: Reading...
data.couchbase-capella_auditlogsettings.existing_auditlogsettings: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_auditlogsettings.new_auditlogsettings will be created
  + resource "couchbase-capella_auditlogsettings" "new_auditlogsettings" {
      + auditenabled    = true
      + cluster_id      = "0fad8c6d-1a42-44b6-ada8-3c7a917114a5"
      + disabledusers   = []
      + enabledeventids = [
          + 20488,
          + 20490,
          + 20491,
        ]
      + organization_id = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
      + project_id      = "6b22a328-3ad8-41c3-8d83-e6242d170423"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_auditlogsettings.new_auditlogsettings: Creating...
couchbase-capella_auditlogsettings.new_auditlogsettings: Creation complete after 2s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

existing_auditlogsettings = {
  "auditenabled" = true
  "cluster_id" = "0fad8c6d-1a42-44b6-ada8-3c7a917114a5"
  "disabledusers" = tolist([])
  "enabledeventids" = tolist([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
  "project_id" = "6b22a328-3ad8-41c3-8d83-e6242d170423"
}
new_auditlogsettings = {
  "auditenabled" = true
  "cluster_id" = "0fad8c6d-1a42-44b6-ada8-3c7a917114a5"
  "disabledusers" = tolist([])
  "enabledeventids" = tolist([
    20488,
    20490,
    20491,
  ])
  "organization_id" = "3c549058-7ca4-4e89-9efb-8ed61a5f3d3f"
  "project_id" = "6b22a328-3ad8-41c3-8d83-e6242d170423"
}

```
