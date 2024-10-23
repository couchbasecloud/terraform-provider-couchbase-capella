# Capella Flush Bucket Example

This example shows how to use flush buckets in Capella.


# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new flush bucket and a new bucket resources as shown in the flush_bucket.tf and create_bucket.tf files respectively. There is no corresponding Capella resource for flush bucket.
2. UPDATE: Does nothing due to flush bucket resource consisting entirely of id's. Will only update the bucket resource. To re-execute the flush bucket call you will need to re-create the flush bucket terraform resource.
3. DELETE: Delete the flush bucket and bucket resources from terraform state.

## CREATE
### View the plan for the resources that Terraform will create


Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_bucket.new_bucket will be created
  + resource "couchbase-capella_bucket" "new_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = true
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000002"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # couchbase-capella_flush.new_flush will be created
  + resource "couchbase-capella_flush" "new_flush" {
      + bucket_id       = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_id  = (known after apply)
  + new_bucket = {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = true
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000002"
      + replicas                   = 0
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }
  + new_flush  = {
      + bucket_id       = (known after apply)
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

  # couchbase-capella_bucket.new_bucket will be created
  + resource "couchbase-capella_bucket" "new_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = true
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000002"
      + replicas                   = 0
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

  # couchbase-capella_flush.new_flush will be created
  + resource "couchbase-capella_flush" "new_flush" {
      + bucket_id       = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + bucket_id  = (known after apply)
  + new_bucket = {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = true
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id                 = "ffffffff-aaaa-1414-eeee-000000000002"
      + replicas                   = 0
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }
  + new_flush  = {
      + bucket_id       = (known after apply)
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000001"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000002"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_bucket.new_bucket: Creating...
couchbase-capella_bucket.new_bucket: Creation complete after 7s [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
couchbase-capella_flush.new_flush: Creating...
couchbase-capella_flush.new_flush: Still creating... [10s elapsed]
couchbase-capella_flush.new_flush: Creation complete after 10s

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

bucket_id = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
new_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = true
  "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "memory_allocation_in_mb" = 100
  "name" = "new_terraform_bucket"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
  "replicas" = 0
  "stats" = {
    "disk_used_in_mib" = 0
    "item_count" = 0
    "memory_used_in_mib" = 0
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
new_flush = {
  "bucket_id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
}

```

## UPDATE

### Run terraform apply again

Sample Output:
```
terraform apply

couchbase-capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
couchbase-capella_flush.new_flush: Refreshing state...

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_bucket.new_bucket has changed
  ~ resource "couchbase-capella_bucket" "new_bucket" {
        id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
        name                       = "new_terraform_bucket"
      ~ stats                      = {
          ~ disk_used_in_mib   = 0 -> 8
          ~ memory_used_in_mib = 0 -> 31
            # (2 unchanged attributes hidden)
        }
        # (12 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or respond to these changes.

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Changes to Outputs:
  ~ new_bucket = {
        id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
        name                       = "new_terraform_bucket"
      ~ stats                      = {
          ~ disk_used_in_mib   = 0 -> 8
          ~ memory_used_in_mib = 0 -> 31
            # (2 unchanged attributes hidden)
        }
        # (12 unchanged attributes hidden)
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

bucket_id = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
new_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = true
  "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "memory_allocation_in_mb" = 100
  "name" = "new_terraform_bucket"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
  "replicas" = 0
  "stats" = {
    "disk_used_in_mib" = 8
    "item_count" = 0
    "memory_used_in_mib" = 31
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
new_flush = {
  "bucket_id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
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


couchbase-capella_flush.new_flush: Refreshing state...
couchbase-capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush must be replaced
-/+ resource "couchbase-capella_flush" "new_flush" {
      ~ bucket_id       = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> "dHJhdmVsLXNhbXBsZQ%3D%3D" # forces replacement
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_flush = {
      ~ bucket_id       = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> "dHJhdmVsLXNhbXBsZQ%3D%3D"
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
couchbase-capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush must be replaced
-/+ resource "couchbase-capella_flush" "new_flush" {
      ~ bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> "Z2FtZXNpbS1zYW1wbGU%3couchbase-capella_flush.new_flush: Refreshing state...
couchbase-capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_flush.new_flush must be replaced
-/+ resource "couchbase-capella_flush" "new_flush" {
      ~ bucket_id       = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> "dHJhdmVsLXNhbXBsZQ%3D%3D" # forces replacement
        # (3 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_flush  = {
      ~ bucket_id       = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> "dHJhdmVsLXNhbXBsZQ%3D%3D"
        # (3 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_flush.new_flush: Destroying...
couchbase-capella_flush.new_flush: Destruction complete after 0s
couchbase-capella_flush.new_flush: Creating...
couchbase-capella_flush.new_flush: Creation complete after 5s

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

bucket_id = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
new_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = true
  "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "memory_allocation_in_mb" = 100
  "name" = "new_terraform_bucket"
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000001"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000002"
  "replicas" = 0
  "stats" = {
    "disk_used_in_mib" = 8
    "item_count" = 0
    "memory_used_in_mib" = 31
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
new_flush = {
  "bucket_id" = "dHJhdmVsLXNhbXBsZQ%3D%3D"
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
couchbase-capella_bucket.new_bucket
couchbase-capella_flush.new_flush
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy

╵
couchbase-capella_flush.new_flush: Refreshing state...
couchbase-capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_bucket.new_bucket will be destroyed
  - resource "couchbase-capella_bucket" "new_bucket" {
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = true -> null
      - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> null
      - memory_allocation_in_mb    = 100 -> null
      - name                       = "new_terraform_bucket" -> null
      - organization_id            = "ffffffff-aaaa-1414-eeee-000000000001" -> null
      - project_id                 = "ffffffff-aaaa-1414-eeee-000000000002" -> null
      - replicas                   = 0 -> null
      - stats                      = {
          - disk_used_in_mib   = 8 -> null
          - item_count         = 0 -> null
          - memory_used_in_mib = 31 -> null
          - ops_per_second     = 0 -> null
        } -> null
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

  # couchbase-capella_flush.new_flush will be destroyed
  - resource "couchbase-capella_flush" "new_flush" {
      - bucket_id       = "dHJhdmVsLXNhbXBsZQ%3D%3D" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000001" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000002" -> null
    }

Plan: 0 to add, 0 to change, 2 to destroy.

Changes to Outputs:
  - bucket_id  = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> null
  - new_bucket = {
      - bucket_conflict_resolution = "seqno"
      - cluster_id                 = "ffffffff-aaaa-1414-eeee-000000000000"
      - durability_level           = "none"
      - eviction_policy            = "fullEviction"
      - flush                      = true
      - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
      - memory_allocation_in_mb    = 100
      - name                       = "new_terraform_bucket"
      - organization_id            = "ffffffff-aaaa-1414-eeee-000000000001"
      - project_id                 = "ffffffff-aaaa-1414-eeee-000000000002"
      - replicas                   = 0
      - stats                      = {
          - disk_used_in_mib   = 8
          - item_count         = 0
          - memory_used_in_mib = 31
          - ops_per_second     = 0
        }
      - storage_backend            = "couchstore"
      - time_to_live_in_seconds    = 0
      - type                       = "couchbase"
    } -> null
  - new_flush  = {
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
couchbase-capella_bucket.new_bucket: Destroying... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
couchbase-capella_flush.new_flush: Destruction complete after 0s
couchbase-capella_bucket.new_bucket: Destruction complete after 2s

Destroy complete! Resources: 2 destroyed.

```
