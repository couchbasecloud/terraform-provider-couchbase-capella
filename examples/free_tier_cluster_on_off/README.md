# Capella Free Tier Cluster On/Off Example

This example demonstrates how to use the Capella Terraform provider to turn off and on a free-tier cluster

To run, configure your Couchbase capella provider as described in README in the root of this project

# Example Walkthrough


In this example, we are going to do the following.

1. CREATE: Switch the free-tier cluster and assosciated app-service if any to on/off
2. UPDATE: Update the current free-tier cluster state to on/off
3. IMPORT: Import a free-tier cluster state that exists in Capella but not in terraform

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform plan will create

Command: `terraform plan`

Sample Output:
```
 $terraform plan 
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off will be created
  + resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
      + cluster_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + organization_id = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + project_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + state           = "off"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_free_tier_cluster_on_off = {
      + cluster_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + organization_id = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + project_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + state           = "off"
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the plan to create the resources

Command: `terraform apply`

Sample Output:
```
 $terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off will be created
  + resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
      + cluster_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + organization_id = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + project_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + state           = "off"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_free_tier_cluster_on_off = {
      + cluster_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + organization_id = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + project_id      = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
      + state           = "off"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Creating...
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still creating... [10s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still creating... [20s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still creating... [30s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still creating... [40s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still creating... [50s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still creating... [1m0s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Creation complete after 1m5s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_free_tier_cluster_on_off = {
  "cluster_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
  "organization_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
  "project_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
  "state" = "off"
}

```

### Verify the resource state

Command: `terraform output new_free_tier_cluster_on_off`

Sample Output:
```
$terraform output new_free_tier_cluster_on_off

{
"cluster_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
"organization_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
"project_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
"state" = "off"
}
```

## IMPORT

### Import a resource that exists in Capella but not in terraform

Command: `terraform import couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off cluster_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc,organization_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc,project_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc`

Sample Output:
```
$terraform import couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off cluster_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc,organization_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc,project_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Importing from ID "cluster_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc,organization_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc,project_id=ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"...
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Import prepared!
  Prepared couchbase-capella_free_tier_cluster_on_off for import
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```
## UPDATE
### Let us edit terraform.tfvars file to change the free-tier cluster state.

Command: `terraform plan`

Sample Output:

```
$terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off will be updated in-place
  ~ resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
      ~ state           = "off" -> "on"
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_free_tier_cluster_on_off = {
      ~ state           = "off" -> "on"
        # (3 unchanged attributes hidden)
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the plan to update the resources

Command: `terraform apply`

Sample Output:
```
$terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off will be updated in-place
  ~ resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
      ~ state           = "off" -> "on"
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_free_tier_cluster_on_off = {
      ~ state           = "off" -> "on"
        # (3 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Modifying...
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [10s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [20s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [30s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [40s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [50s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [1m0s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [1m10s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [1m20s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [1m30s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [1m40s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Still modifying... [1m50s elapsed]
couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off: Modifications complete after 1m55s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

new_free_tier_cluster_on_off = {
  "cluster_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
  "organization_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
  "project_id" = "ffffffff-aaaaaaaa-bbbbbbbbb-cccccc"
  "state" = "on"
}
```
