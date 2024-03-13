# Capella Buckets Example

This example shows how to create and manage Buckets in Capella.

This creates a new bucket in the selected Capella cluster and lists existing buckets in the cluster. It uses the cluster ID to create and list buckets.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new bucket in Capella as stated in the `create_bucket.tf` file.
2. UPDATE: Update the bucket configuration using Terraform.
3. LIST: List existing buckets in Capella as stated in the `list_buckets.tf` file.
4. IMPORT: Import a bucket that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created bucket from Capella.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE & LIST
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_buckets.existing_buckets: Reading...
data.capella_buckets.existing_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + audit                      = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      + replicas                   = 1
      + stats                      = {
          + disk_used_in_mib   = (known after apply)
          + item_count         = (known after apply)
          + memory_used_in_mib = (known after apply)
          + ops_per_second     = (known after apply)
        }
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + buckets_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = [
          + {
              + audit                      = {
                  + created_at  = "0001-01-01 00:00:00 +0000 UTC"
                  + created_by  = ""
                  + modified_at = "0001-01-01 00:00:00 +0000 UTC"
                  + modified_by = ""
                  + version     = 0
                }
              + bucket_conflict_resolution = "seqno"
              + cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + durability_level           = "none"
              + eviction_policy            = "valueOnly"
              + flush                      = false
              + id                         = "dHJhdmVsLXNhbXBsZQ=="
              + memory_allocation_in_mb    = 200
              + name                       = "travel-sample"
              + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + replicas                   = 1
              + stats                      = {
                  + disk_used_in_mib   = 11
                  + item_count         = 0
                  + memory_used_in_mib = 39
                  + ops_per_second     = 0
                }
              + storage_backend            = "couchstore"
              + time_to_live_in_seconds    = 0
              + type                       = "couchbase"
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_bucket   = {
      + audit                      = (known after apply)
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Bucket

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_buckets.existing_buckets: Reading...
data.capella_buckets.existing_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + audit                      = {
          + created_at  = (known after apply)
          + created_by  = (known after apply)
          + modified_at = (known after apply)
          + modified_by = (known after apply)
          + version     = (known after apply)
        }
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      + replicas                   = 1
      + stats                      = {
          + disk_used_in_mib   = (known after apply)
          + item_count         = (known after apply)
          + memory_used_in_mib = (known after apply)
          + ops_per_second     = (known after apply)
        }
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + buckets_list = {
      + cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + data            = [
          + {
              + audit                      = {
                  + created_at  = "0001-01-01 00:00:00 +0000 UTC"
                  + created_by  = ""
                  + modified_at = "0001-01-01 00:00:00 +0000 UTC"
                  + modified_by = ""
                  + version     = 0
                }
              + bucket_conflict_resolution = "seqno"
              + cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              + durability_level           = "none"
              + eviction_policy            = "valueOnly"
              + flush                      = false
              + id                         = "dHJhdmVsLXNhbXBsZQ=="
              + memory_allocation_in_mb    = 200
              + name                       = "travel-sample"
              + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              + project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              + replicas                   = 1
              + stats                      = {
                  + disk_used_in_mib   = 134
                  + item_count         = 63288
                  + memory_used_in_mib = 165
                  + ops_per_second     = 0
                }
              + storage_backend            = "couchstore"
              + time_to_live_in_seconds    = 0
              + type                       = "couchbase"
            },
        ]
      + organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    }
  + new_bucket   = {
      + audit                      = (known after apply)
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      + durability_level           = "none"
      + eviction_policy            = "fullEviction"
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 100
      + name                       = "new_terraform_bucket"
      + organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      + project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = "couchstore"
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_bucket.new_bucket: Creating...
capella_bucket.new_bucket: Creation complete after 8s [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

buckets_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "0001-01-01 00:00:00 +0000 UTC"
        "created_by" = ""
        "modified_at" = "0001-01-01 00:00:00 +0000 UTC"
        "modified_by" = ""
        "version" = 0
      }
      "bucket_conflict_resolution" = "seqno"
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "durability_level" = "none"
      "eviction_policy" = "valueOnly"
      "flush" = false
      "id" = "dHJhdmVsLXNhbXBsZQ=="
      "memory_allocation_in_mb" = 200
      "name" = "travel-sample"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "replicas" = 1
      "stats" = {
        "disk_used_in_mib" = 134
        "item_count" = 63288
        "memory_used_in_mib" = 165
        "ops_per_second" = 0
      }
      "storage_backend" = "couchstore"
      "time_to_live_in_seconds" = 0
      "type" = "couchbase"
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_bucket = {
  "audit" = {
    "created_at" = "0001-01-01 00:00:00 +0000 UTC"
    "created_by" = ""
    "modified_at" = "0001-01-01 00:00:00 +0000 UTC"
    "modified_by" = ""
    "version" = 0
  }
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "memory_allocation_in_mb" = 100
  "name" = "new_terraform_bucket"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
  "replicas" = 1
  "stats" = {
    "disk_used_in_mib" = 1
    "item_count" = 0
    "memory_used_in_mib" = 38
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
```

### Note the Bucket ID for the new Bucket
Command: `terraform output new_bucket`

Sample Output:
```
$ terraform output new_bucket
{
  "audit" = {
    "created_at" = "0001-01-01 00:00:00 +0000 UTC"
    "created_by" = ""
    "modified_at" = "0001-01-01 00:00:00 +0000 UTC"
    "modified_by" = ""
    "version" = 0
  }
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "memory_allocation_in_mb" = 100
  "name" = "new_terraform_bucket"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
  "replicas" = 1
  "stats" = {
    "disk_used_in_mib" = 1
    "item_count" = 0
    "memory_used_in_mib" = 38
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
```

In this case, the bucket ID for my new bucket is `bmV3X3RlcnJhZm9ybV9idWNrZXQ=`

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
$ terraform state list
data.couchbase-capella_buckets.existing_buckets
couchbase-capella_bucket.new_bucket
```

## IMPORT
### Remove the resource `new_bucket` from the Terraform State file

Command: `terraform state rm couchbase-capella_bucket.new_bucket`

Sample Output:
```
$ terraform state rm couchbase-capella_bucket.new_bucket
Removed capella_bucket.new_bucket
Successfully removed 1 resource instance(s).
```

Please note, this command will only remove the resource from the Terraform State file, but in reality, the resource exists in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_bucket.new_bucket id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_bucket.new_bucket id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d`

Sample Output:
```
$ terraform import couchbase-capella_bucket.new_bucket id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d
capella_bucket.new_bucket: Importing from ID "id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d"...
capella_bucket.new_bucket: Import prepared!
  Prepared capella_bucket for import
capella_bucket.new_bucket: Refreshing state... [id=id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=,cluster_id=f499a9e6-e5a1-4f3e-95a7-941a41d046e6,project_id=958ad6b5-272d-49f0-babd-cc98c6b54a81,organization_id=0783f698-ac58-4018-84a3-31c3b6ef785d]
data.capella_buckets.existing_buckets: Reading...
data.capella_buckets.existing_buckets: Read complete after 1s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the bucket ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the bucket belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_buckets.existing_buckets: Reading...
capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
data.capella_buckets.existing_buckets: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the bucket configuration settings.

Command: `terraform apply -var 'bucket={durability_level="majority",name="new_terraform_bucket"}'`

Sample Output:
```
$ terraform apply -var 'bucket={durability_level="majority",name="new_terraform_bucket"}'
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_buckets.existing_buckets: Reading...
capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
data.capella_buckets.existing_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be updated in-place
  ~ resource "capella_bucket" "new_bucket" {
      ~ audit                      = {
          ~ created_at  = "0001-01-01 00:00:00 +0000 UTC" -> (known after apply)
          + created_by  = (known after apply)
          ~ modified_at = "0001-01-01 00:00:00 +0000 UTC" -> (known after apply)
          + modified_by = (known after apply)
          ~ version     = 0 -> (known after apply)
        }
      ~ durability_level           = "none" -> "majority"
        id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
        name                       = "new_terraform_bucket"
      ~ stats                      = {
          ~ disk_used_in_mib   = 17 -> (known after apply)
          ~ item_count         = 0 -> (known after apply)
          ~ memory_used_in_mib = 54 -> (known after apply)
          ~ ops_per_second     = 0 -> (known after apply)
        }
        # (11 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ new_bucket = {
      ~ audit                      = {
          - created_at  = "0001-01-01 00:00:00 +0000 UTC"
          - created_by  = ""
          - modified_at = "0001-01-01 00:00:00 +0000 UTC"
          - modified_by = ""
          - version     = 0
        } -> (known after apply)
      ~ durability_level           = "none" -> "majority"
        id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
        name                       = "new_terraform_bucket"
      ~ stats                      = {
          - disk_used_in_mib   = 17
          - item_count         = 0
          - memory_used_in_mib = 54
          - ops_per_second     = 0
        } -> (known after apply)
        # (11 unchanged elements hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_bucket.new_bucket: Modifying... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_bucket.new_bucket: Modifications complete after 2s [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

buckets_list = {
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "data" = tolist([
    {
      "audit" = {
        "created_at" = "0001-01-01 00:00:00 +0000 UTC"
        "created_by" = ""
        "modified_at" = "0001-01-01 00:00:00 +0000 UTC"
        "modified_by" = ""
        "version" = 0
      }
      "bucket_conflict_resolution" = "seqno"
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "durability_level" = "none"
      "eviction_policy" = "fullEviction"
      "flush" = false
      "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
      "memory_allocation_in_mb" = 100
      "name" = "new_terraform_bucket"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "replicas" = 1
      "stats" = {
        "disk_used_in_mib" = 17
        "item_count" = 0
        "memory_used_in_mib" = 54
        "ops_per_second" = 0
      }
      "storage_backend" = "couchstore"
      "time_to_live_in_seconds" = 0
      "type" = "couchbase"
    },
    {
      "audit" = {
        "created_at" = "0001-01-01 00:00:00 +0000 UTC"
        "created_by" = ""
        "modified_at" = "0001-01-01 00:00:00 +0000 UTC"
        "modified_by" = ""
        "version" = 0
      }
      "bucket_conflict_resolution" = "seqno"
      "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      "durability_level" = "none"
      "eviction_policy" = "valueOnly"
      "flush" = false
      "id" = "dHJhdmVsLXNhbXBsZQ=="
      "memory_allocation_in_mb" = 200
      "name" = "travel-sample"
      "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      "replicas" = 1
      "stats" = {
        "disk_used_in_mib" = 134
        "item_count" = 63288
        "memory_used_in_mib" = 165
        "ops_per_second" = 0
      }
      "storage_backend" = "couchstore"
      "time_to_live_in_seconds" = 0
      "type" = "couchbase"
    },
  ])
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
}
new_bucket = {
  "audit" = {
    "created_at" = "0001-01-01 00:00:00 +0000 UTC"
    "created_by" = ""
    "modified_at" = "0001-01-01 00:00:00 +0000 UTC"
    "modified_by" = ""
    "version" = 0
  }
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
  "durability_level" = "majority"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
  "memory_allocation_in_mb" = 100
  "name" = "new_terraform_bucket"
  "organization_id" = "0783f698-ac58-4018-84a3-31c3b6ef785d"
  "project_id" = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
  "replicas" = 1
  "stats" = {
    "disk_used_in_mib" = 17
    "item_count" = 0
    "memory_used_in_mib" = 54
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
```

## DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:
```
$ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchbasecloud/capella in /Users/$USER/workspace/terraform-provider-capella
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.capella_buckets.existing_buckets: Reading...
capella_bucket.new_bucket: Refreshing state... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
data.capella_buckets.existing_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be destroyed
  - resource "capella_bucket" "new_bucket" {
      - audit                      = {
          - created_at  = "0001-01-01 00:00:00 +0000 UTC" -> null
          - modified_at = "0001-01-01 00:00:00 +0000 UTC" -> null
          - version     = 0 -> null
        }
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6" -> null
      - durability_level           = "majority" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ=" -> null
      - memory_allocation_in_mb    = 100 -> null
      - name                       = "new_terraform_bucket" -> null
      - organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d" -> null
      - project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81" -> null
      - replicas                   = 1 -> null
      - stats                      = {
          - disk_used_in_mib   = 17 -> null
          - item_count         = 0 -> null
          - memory_used_in_mib = 54 -> null
          - ops_per_second     = 0 -> null
        }
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - buckets_list = {
      - cluster_id      = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      - data            = [
          - {
              - audit                      = {
                  - created_at  = "0001-01-01 00:00:00 +0000 UTC"
                  - created_by  = ""
                  - modified_at = "0001-01-01 00:00:00 +0000 UTC"
                  - modified_by = ""
                  - version     = 0
                }
              - bucket_conflict_resolution = "seqno"
              - cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - durability_level           = "majority"
              - eviction_policy            = "fullEviction"
              - flush                      = false
              - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
              - memory_allocation_in_mb    = 100
              - name                       = "new_terraform_bucket"
              - organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              - replicas                   = 1
              - stats                      = {
                  - disk_used_in_mib   = 17
                  - item_count         = 0
                  - memory_used_in_mib = 54
                  - ops_per_second     = 0
                }
              - storage_backend            = "couchstore"
              - time_to_live_in_seconds    = 0
              - type                       = "couchbase"
            },
          - {
              - audit                      = {
                  - created_at  = "0001-01-01 00:00:00 +0000 UTC"
                  - created_by  = ""
                  - modified_at = "0001-01-01 00:00:00 +0000 UTC"
                  - modified_by = ""
                  - version     = 0
                }
              - bucket_conflict_resolution = "seqno"
              - cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
              - durability_level           = "none"
              - eviction_policy            = "valueOnly"
              - flush                      = false
              - id                         = "dHJhdmVsLXNhbXBsZQ=="
              - memory_allocation_in_mb    = 200
              - name                       = "travel-sample"
              - organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
              - project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
              - replicas                   = 1
              - stats                      = {
                  - disk_used_in_mib   = 134
                  - item_count         = 63288
                  - memory_used_in_mib = 165
                  - ops_per_second     = 0
                }
              - storage_backend            = "couchstore"
              - time_to_live_in_seconds    = 0
              - type                       = "couchbase"
            },
        ]
      - organization_id = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id      = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
    } -> null
  - new_bucket   = {
      - audit                      = {
          - created_at  = "0001-01-01 00:00:00 +0000 UTC"
          - created_by  = ""
          - modified_at = "0001-01-01 00:00:00 +0000 UTC"
          - modified_by = ""
          - version     = 0
        }
      - bucket_conflict_resolution = "seqno"
      - cluster_id                 = "f499a9e6-e5a1-4f3e-95a7-941a41d046e6"
      - durability_level           = "majority"
      - eviction_policy            = "fullEviction"
      - flush                      = false
      - id                         = "bmV3X3RlcnJhZm9ybV9idWNrZXQ="
      - memory_allocation_in_mb    = 100
      - name                       = "new_terraform_bucket"
      - organization_id            = "0783f698-ac58-4018-84a3-31c3b6ef785d"
      - project_id                 = "958ad6b5-272d-49f0-babd-cc98c6b54a81"
      - replicas                   = 1
      - stats                      = {
          - disk_used_in_mib   = 17
          - item_count         = 0
          - memory_used_in_mib = 54
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

capella_bucket.new_bucket: Destroying... [id=bmV3X3RlcnJhZm9ybV9idWNrZXQ=]
capella_bucket.new_bucket: Destruction complete after 3s

Destroy complete! Resources: 1 destroyed.
```
