# Capella Flush Bucket Example

This example shows how to use flush buckets in Capella.


# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new flush bucket resource as shown in the flush_bucket.tf file. There is no corresponding Capella resource.
2. UPDATE: Does nothing due to flush bucket resource consisting entirely of id's.
3. DELETE: Delete the flush bucket resource from terraform state.

## CREATE
### View the plan for the resources that Terraform will create


Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush will be created
  + resource "couchbase-capella_flush" "new_flush" {
      + bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_flush = {
      + bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

```

### Apply the Plan, in order to flush document in the bucket

Command: `terraform apply`

Sample Output:
```
$ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush will be created
  + resource "couchbase-capella_flush" "new_flush" {
      + bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_flush = {
      + bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_flush.new_flush: Creating...
couchbase-capella_flush.new_flush: Still creating... [10s elapsed]
couchbase-capella_flush.new_flush: Creation complete after 11s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_flush = {
  "bucket_id" = "dHJhdmVsLXNhbXBsZQ%3D%3D"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
}

```

## UPDATE

### Run terraform apply again

Sample Output:
```
╵$ terraform apply

couchbase-capella_flush.new_flush: Refreshing state...

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

new_flush = {
  "bucket_id" = "dHJhdmVsLXNhbXBsZQ%3D%3D"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
}


### Note the output for the new flush bucket resource
Command: `terraform output new_flush`

Sample Output:
```

terraform output new_flush
{
  "bucket_id" = "dHJhdmVsLXNhbXBsZQ%3D%3D"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
}

```


### Change the bucket or one of the other id's.

Command: `terraform plan`

Sample Output:
```
╵$ terraform plan


Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush must be replaced
-/+ resource "couchbase-capella_flush" "new_flush" {
      ~ bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> "Z2FtZXNpbS1zYW1wbGU%3D" # forces replacement
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_flush = {
      ~ bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> "Z2FtZXNpbS1zYW1wbGU%3D"
        # (3 unchanged attributes hidden)
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.

```

Command: `terraform apply`

Sample Output:
```
╵$ terraform apply
couchbase-capella_flush.new_flush: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush must be replaced
-/+ resource "couchbase-capella_flush" "new_flush" {
      ~ bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> "Z2FtZXNpbS1zYW1wbGU%3D" # forces replacement
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_flush = {
      ~ bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> "Z2FtZXNpbS1zYW1wbGU%3D"
        # (3 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_flush.new_flush: Destroying...
couchbase-capella_flush.new_flush: Destruction complete after 0s
couchbase-capella_flush.new_flush: Creating...
couchbase-capella_flush.new_flush: Still creating... [10s elapsed]
couchbase-capella_flush.new_flush: Creation complete after 11s

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

new_flush = {
  "bucket_id" = "Z2FtZXNpbS1zYW1wbGU%3D"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
}


```


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```

$ terraform state list
couchbase-capella_flush.new_flush

```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy

couchbase-capella_flush.new_flush: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush will be destroyed
  - resource "couchbase-capella_flush" "new_flush" {
      - bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000001" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000002" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_flush = {
      - bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D"
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_flush.new_flush: Destroying...
couchbase-capella_flush.new_flush: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.

```
