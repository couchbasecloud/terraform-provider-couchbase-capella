# Capella Free Tier Buckets Examples

This example shows how to manage Free Tier Buckets in Capella

This creates a new free-tier bucket in the selected Capella free-tier cluster and lists existing free-tier buckets in the cluster. It uses the cluster ID to create and list free-tier buckets.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new free-tier bucket in Capella as stated in the `create_free_tier_bucket.tf` file.
2. LIST: List the buckets in the Capella free-tier cluster using Terraform.
3. UPDATE: Update the free-tier bucket configuration using Terraform.
3. IMPORT: Import a free-tier bucket that exists in Capella but not in the terraform state file.
4. DELETE: Delete the newly created free-tier bucket from Capella.

## CREATE

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output
```
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_free_tier_bucket.new_free_tier_bucket will be created
  + resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
      + bucket_conflict_resolution = (known after apply)
      + cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + durability_level           = (known after apply)
      + eviction_policy            = (known after apply)
      + flush                      = (known after apply)
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 250
      + name                       = "test_bucket"
      + organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + replicas                   = (known after apply)
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = (known after apply)
      + type                       = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + free_tier_bucket_id  = (known after apply)
  + new_free_tier_bucket = {
      + bucket_conflict_resolution = (known after apply)
      + cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + durability_level           = (known after apply)
      + eviction_policy            = (known after apply)
      + flush                      = (known after apply)
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 250
      + name                       = "test_bucket"
      + organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + replicas                   = (known after apply)
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = (known after apply)
      + type                       = (known after apply)
    }

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Bucket

Command: `terraform apply`

Sample Output:
```
$terraform apply    
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_free_tier_bucket.new_free_tier_bucket will be created
  + resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
      + bucket_conflict_resolution = (known after apply)
      + cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + durability_level           = (known after apply)
      + eviction_policy            = (known after apply)
      + flush                      = (known after apply)
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 150
      + name                       = "test_bucket"
      + organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + replicas                   = (known after apply)
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = (known after apply)
      + type                       = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + free_tier_bucket_id  = (known after apply)
  + new_free_tier_bucket = {
      + bucket_conflict_resolution = (known after apply)
      + cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + durability_level           = (known after apply)
      + eviction_policy            = (known after apply)
      + flush                      = (known after apply)
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 150
      + name                       = "test_bucket"
      + organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + replicas                   = (known after apply)
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = (known after apply)
      + type                       = (known after apply)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_free_tier_bucket.new_free_tier_bucket: Creating...
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Still creating... [10s elapsed]
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Creation complete after 14s [id=dGVzdF9idWNrZXQ=]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

free_tier_bucket_id = "dGVzdF9idWNrZXQ="
new_free_tier_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "dGVzdF9idWNrZXQ="
  "memory_allocation_in_mb" = 150
  "name" = "test_bucket"
  "organization_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "project_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
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

```

## LIST

### List the free-tier buckets in the Capella free-tier cluster using Terraform

Command: `terraform plan`

Sample Output:
```
$terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_free_tier_buckets.existing_buckets: Reading...
data.couchbase-capella_free_tier_buckets.existing_buckets: Read complete after 2s

Changes to Outputs:
  + buckets_list = {
      + cluster_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + data            = [
          + {
              + bucket_conflict_resolution = "seqno"
              + cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              + durability_level           = "none"
              + eviction_policy            = "fullEviction"
              + flush                      = false
              + id                         = "dHJhdmVsLXNhbXBsZQ=="
              + memory_allocation_in_mb    = 200
              + name                       = "travel-sample"
              + organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              + project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              + replicas                   = 0
              + stats                      = {
                  + disk_used_in_mib   = 66
                  + item_count         = 63343
                  + memory_used_in_mib = 80
                  + ops_per_second     = 0
                }
              + storage_backend            = "couchstore"
              + time_to_live_in_seconds    = 0
              + type                       = "couchbase"
            },
        ]
      + organization_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + project_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

## Apply the plan in order to fetch all the free-tier buckets of the cluster

Command: `terraform apply`

Sample Output:
```
$terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_free_tier_buckets.existing_buckets: Reading...
data.couchbase-capella_free_tier_buckets.existing_buckets: Read complete after 1s

Changes to Outputs:
  + buckets_list = {
      + cluster_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + data            = [
          + {
              + bucket_conflict_resolution = "seqno"
              + cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              + durability_level           = "none"
              + eviction_policy            = "fullEviction"
              + flush                      = false
              + id                         = "dHJhdmVsLXNhbXBsZQ=="
              + memory_allocation_in_mb    = 200
              + name                       = "travel-sample"
              + organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              + project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              + replicas                   = 0
              + stats                      = {
                  + disk_used_in_mib   = 66
                  + item_count         = 63343
                  + memory_used_in_mib = 80
                  + ops_per_second     = 0
                }
              + storage_backend            = "couchstore"
              + time_to_live_in_seconds    = 0
              + type                       = "couchbase"
            },
        ]
      + organization_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      + project_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes


Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

buckets_list = {
  "cluster_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "data" = tolist([
    {
      "bucket_conflict_resolution" = "seqno"
      "cluster_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      "durability_level" = "none"
      "eviction_policy" = "fullEviction"
      "flush" = false
      "id" = "dHJhdmVsLXNhbXBsZQ=="
      "memory_allocation_in_mb" = 200
      "name" = "travel-sample"
      "organization_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      "project_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      "replicas" = 0
      "stats" = {
        "disk_used_in_mib" = 66
        "item_count" = 63343
        "memory_used_in_mib" = 80
        "ops_per_second" = 0
      }
      "storage_backend" = "couchstore"
      "time_to_live_in_seconds" = 0
      "type" = "couchbase"
    },
  ])
  "organization_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "project_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}
```

## UPDATE

### Update the free-tier bucket configuration using Terraform
In the below example, we are updating the memory allocation of the free-tier bucket from 200MB to 250MB.

Command: `terraform apply`

Sample Output:
```
$terraform apply              
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Refreshing state... [id=dGVzdF9idWNrZXQ=]
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place
Terraform will perform the following actions:
  # couchbase-capella_free_tier_bucket.new_free_tier_bucket will be updated in-place
  ~ resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
        id                         = "dGVzdF9idWNrZXQ="
      ~ memory_allocation_in_mb    = 360 -> 300
        name                       = "test_bucket"
        # (12 unchanged attributes hidden)
    }
Plan: 0 to add, 1 to change, 0 to destroy.
Changes to Outputs:
  ~ new_free_tier_bucket = {
        id                         = "dGVzdF9idWNrZXQ="
      ~ memory_allocation_in_mb    = 360 -> 300
        name                       = "test_bucket"
        # (12 unchanged attributes hidden)
    }
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.
  Enter a value: yes
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Modifying... [id=dGVzdF9idWNrZXQ=]
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Modifications complete after 6s [id=dGVzdF9idWNrZXQ=]
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
Outputs:
free_tier_bucket_id = "dGVzdF9idWNrZXQ="
new_free_tier_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "dGVzdF9idWNrZXQ="
  "memory_allocation_in_mb" = 300
  "name" = "test_bucket"
  "organization_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "project_id" = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  "replicas" = 0
  "stats" = {
    "disk_used_in_mib" = 8
    "item_count" = 0
    "memory_used_in_mib" = 28
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}

```


## Import

### Remove the free-tier bucket from the terraform state file

Command: `terraform state rm couchbase-capella_free_tier_bucket.new_free_tier_bucket`

Sample Output:
```
$terraform state list
couchbase-capella_free_tier_bucket.new_free_tier_bucket

$terraform state rm couchbase-capella_free_tier_bucket.new_free_tier_bucket
Removed couchbase-capella_free_tier_bucket.new_free_tier_bucket
Successfully removed 1 resource instance(s).
```


### Import a free-tier bucket that exists in Capella but not in the terraform state file

Command: `id=dGVzdF9idWNrZXQ=,organization_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,project_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,cluster_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee`

Sample Output
```
terraform import couchbase-capella_free_tier_bucket.new_free_tier_bucket id=dGVzdF9idWNrZXQ=,organization_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,project_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,cluster_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Importing from ID "id=dGVzdF9idWNrZXQ=,organization_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,project_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,cluster_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"...
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Import prepared!
  Prepared couchbase-capella_free_tier_bucket for import
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Refreshing state... [id=id=dGVzdF9idWNrZXQ=,organization_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,project_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee,cluster_id=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee]

Import successful!
The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.


```

## DESTROY

### Delete the newly created free-tier bucket from Capella

Command: `terraform destroy`

Sample Output:
```
$terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/couchbase-capella in /Users/saicharanrachamadugu/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_free_tier_buckets.existing_buckets: Reading...
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Refreshing state... [id=dGVzdF9idWNrZXQ=]
data.couchbase-capella_free_tier_buckets.existing_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_free_tier_bucket.new_free_tier_bucket will be destroyed
  - resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "dGVzdF9idWNrZXQ=" -> null
      - memory_allocation_in_mb    = 250 -> null
      - name                       = "test_bucket" -> null
      - organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" -> null
      - project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" -> null
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

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - buckets_list         = {
      - cluster_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      - data            = [
          - {
              - bucket_conflict_resolution = "seqno"
              - cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - durability_level           = "none"
              - eviction_policy            = "fullEviction"
              - flush                      = false
              - id                         = "U2FpY2hhcmFuLXRlc3Q="
              - memory_allocation_in_mb    = 100
              - name                       = "Saicharan-test"
              - organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
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
            },
          - {
              - bucket_conflict_resolution = "seqno"
              - cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - durability_level           = "none"
              - eviction_policy            = "fullEviction"
              - flush                      = false
              - id                         = "dGVzdF9idWNrZXQ="
              - memory_allocation_in_mb    = 250
              - name                       = "test_bucket"
              - organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
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
            },
          - {
              - bucket_conflict_resolution = "seqno"
              - cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - durability_level           = "none"
              - eviction_policy            = "fullEviction"
              - flush                      = false
              - id                         = "dHJhdmVsLXNhbXBsZQ=="
              - memory_allocation_in_mb    = 200
              - name                       = "travel-sample"
              - organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
              - replicas                   = 0
              - stats                      = {
                  - disk_used_in_mib   = 73
                  - item_count         = 63343
                  - memory_used_in_mib = 81
                  - ops_per_second     = 0
                }
              - storage_backend            = "couchstore"
              - time_to_live_in_seconds    = 0
              - type                       = "couchbase"
            },
        ]
      - organization_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      - project_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    } -> null
  - free_tier_bucket_id  = "dGVzdF9idWNrZXQ=" -> null
  - new_free_tier_bucket = {
      - bucket_conflict_resolution = "seqno"
      - cluster_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      - durability_level           = "none"
      - eviction_policy            = "fullEviction"
      - flush                      = false
      - id                         = "dGVzdF9idWNrZXQ="
      - memory_allocation_in_mb    = 250
      - name                       = "test_bucket"
      - organization_id            = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      - project_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
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

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_free_tier_bucket.new_free_tier_bucket: Destroying... [id=dGVzdF9idWNrZXQ=]
couchbase-capella_free_tier_bucket.new_free_tier_bucket: Destruction complete after 2s

Destroy complete! Resources: 1 destroyed.
```