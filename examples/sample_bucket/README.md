# Capella Sample Buckets Example
This example shows how to create and manage sample Buckets in Capella.

This creates a new bucket in the selected Capella cluster and lists existing sample buckets in the cluster. It uses the cluster ID to create and list buckets.

To run, configure your Couchbase Capella provider as described in README in the root of this project.


# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new sample bucket in Capella as stated in the `create_sample_bucket.tf` file.
2. UPDATE: Update the bucket configuration using Terraform.
3. LIST: List existing sample buckets in Capella as stated in the `list_sample_buckets.tf` file.
4. IMPORT: Import a bucket that exists in Capella but not in the terraform state file.
5. DELETE: Delete the newly created bucket from Capella.
c
If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.


## CREATE & LIST
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:

```
Terraform plan 
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/code/Lagher0/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_sample_buckets.existing_sample_buckets: Reading...
data.couchbase-capella_sample_buckets.existing_sample_buckets: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_sample_bucket.new_sample_bucket will be created
  + resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      + durability_level           = "none"
      + eviction_policy            = (known after apply)
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 200
      + name                       = "gamesim-sample"
      + organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_sample_bucket   = {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      + durability_level           = "none"
      + eviction_policy            = (known after apply)
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 200
      + name                       = "gamesim-sample"
      + organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }
  + samplebucket_id    = (known after apply)
  + samplebuckets_list = {
      + cluster_id      = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      + data            = null
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
    }

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```


### Apply the Plan, in order to create a new sample Bucket

Command: `terraform apply`

```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/workspace/code/Lagher0/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_sample_buckets.existing_sample_buckets: Reading...
data.couchbase-capella_sample_buckets.existing_sample_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_sample_bucket.new_sample_bucket will be created
  + resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      + durability_level           = "none"
      + eviction_policy            = (known after apply)
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 200
      + name                       = "gamesim-sample"
      + organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_sample_bucket   = {
      + bucket_conflict_resolution = "seqno"
      + cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      + durability_level           = "none"
      + eviction_policy            = (known after apply)
      + flush                      = false
      + id                         = (known after apply)
      + memory_allocation_in_mb    = 200
      + name                       = "gamesim-sample"
      + organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      + replicas                   = 1
      + stats                      = (known after apply)
      + storage_backend            = (known after apply)
      + time_to_live_in_seconds    = 0
      + type                       = "couchbase"
    }
  + samplebucket_id    = (known after apply)
  + samplebuckets_list = {
      + cluster_id      = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      + data            = null
      + organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_sample_bucket.new_sample_bucket: Creating...
couchbase-capella_sample_bucket.new_sample_bucket: Creation complete after 1s [id=Z2FtZXNpbS1zYW1wbGU=]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_sample_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "Z2FtZXNpbS1zYW1wbGU="
  "memory_allocation_in_mb" = 200
  "name" = "gamesim-sample"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "replicas" = 1
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
samplebucket_id = "Z2FtZXNpbS1zYW1wbGU="
samplebuckets_list = {
  "cluster_id" = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
  "data" = tolist(null) /* of object */
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
}
```



### Note the Bucket ID for the new sample Bucket
Command: `terraform output new_sample_bucket`

Sample Output:
```
terraform output new_sample_bucket
{
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "Z2FtZXNpbS1zYW1wbGU="
  "memory_allocation_in_mb" = 200
  "name" = "gamesim-sample"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "replicas" = 1
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


### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
 terraform state list
data.couchbase-capella_sample_buckets.existing_sample_buckets
couchbase-capella_sample_bucket.new_sample_bucket
```


## IMPORT
### Remove the resource `new_sample_bucket` from the Terraform State file

Command: `terraform state rm couchbase-capella_sample_bucket.new_sample_bucket`

Sample Output:

```
terraform state rm couchbase-capella_sample_bucket.new_sample_bucket
Removed couchbase-capella_sample_bucket.new_sample_bucket
Successfully removed 1 resource instance(s).
```


### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_sample_bucket.new_sample_bucket id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

In this case, the complete command is:
`terraform import couchbase-capella_sample_bucket.new_sample_bucket  id=Z2FtZXNpbS1zYW1wbGU=,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,cluster_id=17619f3c-08f5-40a3-8c0c-d2e5b263a025`

```
terraform import couchbase-capella_sample_bucket.new_sample_bucket id=Z2FtZXNpbS1zYW1wbGU=,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,cluster_id=17619f3c-08f5-40a3-8c0c-d2e5b263a025
couchbase-capella_sample_bucket.new_sample_bucket: Importing from ID "id=Z2FtZXNpbS1zYW1wbGU=,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,cluster_id=17619f3c-08f5-40a3-8c0c-d2e5b263a025"...
couchbase-capella_sample_bucket.new_sample_bucket: Import prepared!
  Prepared couchbase-capella_sample_bucket for import
data.couchbase-capella_sample_buckets.existing_sample_buckets: Reading...
couchbase-capella_sample_bucket.new_sample_bucket: Refreshing state... [id=id=Z2FtZXNpbS1zYW1wbGU=,organization_id=6af08c0a-8cab-4c1c-b257-b521575c16d0,project_id=c1fade1a-9f27-4a3c-af73-d1b2301890e3,cluster_id=17619f3c-08f5-40a3-8c0c-d2e5b263a025]
data.couchbase-capella_sample_buckets.existing_sample_buckets: Read complete after 2s

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```


Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the sample bucket ID i.e. the ID of the resource that we want to import.
The second ID is the cluster ID i.e. the ID of the cluster to which the sample bucket belongs.
The third ID is the project ID i.e. the ID of the project to which the cluster belongs.
The fourth ID is the organization ID i.e. the ID of the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:

```
terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in $HOME/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_sample_buckets.existing_sample_buckets: Reading...
couchbase-capella_sample_bucket.new_sample_bucket: Refreshing state... [id=Z2FtZXNpbS1zYW1wbGU=]
data.couchbase-capella_sample_buckets.existing_sample_buckets: Read complete after 1s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

## UPDATE
### Let us edit the terraform.tfvars file to change the bucket configuration settings.

Sample buckets does not support update functionality. To update the terraform state it recreates the
sample bucket with the given changes

Command: `terraform apply -var 'samplebucket={name="travel-sample"}'`

Sample Output:

```

 terraform apply -var 'samplebucket={name="travel-sample"}'
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in $HOME/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_sample_buckets.existing_sample_buckets: Reading...
couchbase-capella_sample_bucket.new_sample_bucket: Refreshing state... [id=Z2FtZXNpbS1zYW1wbGU=]
data.couchbase-capella_sample_buckets.existing_sample_buckets: Read complete after 1s

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_sample_bucket.new_sample_bucket has changed
  ~ resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
        id                         = "Z2FtZXNpbS1zYW1wbGU="
        name                       = "gamesim-sample"
      ~ stats                      = {
          ~ item_count         = 196 -> 390
          ~ memory_used_in_mib = 20 -> 42
            # (2 unchanged attributes hidden)
        }
        # (12 unchanged attributes hidden)
    }


Unless you have made equivalent changes to your configuration, or ignored the relevant attributes using ignore_changes, the following plan may include actions to undo or respond to these changes.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # couchbase-capella_sample_bucket.new_sample_bucket must be replaced
-/+ resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      ~ eviction_policy            = "fullEviction" -> (known after apply)
      ~ id                         = "Z2FtZXNpbS1zYW1wbGU=" -> (known after apply)
      ~ name                       = "gamesim-sample" -> "travel-sample" # forces replacement
      ~ stats                      = {
          ~ disk_used_in_mib   = 0 -> (known after apply)
          ~ item_count         = 390 -> (known after apply)
          ~ memory_used_in_mib = 42 -> (known after apply)
          ~ ops_per_second     = 0 -> (known after apply)
        } -> (known after apply)
      ~ storage_backend            = "couchstore" -> (known after apply)
        # (10 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  ~ new_sample_bucket   = {
      ~ eviction_policy            = "fullEviction" -> (known after apply)
      ~ id                         = "Z2FtZXNpbS1zYW1wbGU=" -> (known after apply)
      ~ name                       = "gamesim-sample" -> "travel-sample"
      ~ stats                      = {
          - disk_used_in_mib   = 0
          - item_count         = 196
          - memory_used_in_mib = 20
          - ops_per_second     = 0
        } -> (known after apply)
      ~ storage_backend            = "couchstore" -> (known after apply)
        # (10 unchanged attributes hidden)
    }
  ~ samplebucket_id    = "Z2FtZXNpbS1zYW1wbGU=" -> (known after apply)
  ~ samplebuckets_list = {
      ~ data            = [
          ~ {
              ~ id                         = "dHJhdmVsLXNhbXBsZQ==" -> "Z2FtZXNpbS1zYW1wbGU="
              ~ name                       = "travel-sample" -> "gamesim-sample"
              ~ stats                      = {
                  ~ disk_used_in_mib   = 15 -> 0
                  ~ item_count         = 163 -> 390
                  ~ memory_used_in_mib = 72 -> 42
                    # (1 unchanged attribute hidden)
                }
                # (12 unchanged attributes hidden)
            },
        ]
        # (3 unchanged attributes hidden)
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_sample_bucket.new_sample_bucket: Destroying... [id=Z2FtZXNpbS1zYW1wbGU=]
couchbase-capella_sample_bucket.new_sample_bucket: Destruction complete after 1s
couchbase-capella_sample_bucket.new_sample_bucket: Creating...
couchbase-capella_sample_bucket.new_sample_bucket: Creation complete after 0s [id=dHJhdmVsLXNhbXBsZQ==]

Apply complete! Resources: 1 added, 0 changed, 1 destroyed.

Outputs:

new_sample_bucket = {
  "bucket_conflict_resolution" = "seqno"
  "cluster_id" = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
  "durability_level" = "none"
  "eviction_policy" = "fullEviction"
  "flush" = false
  "id" = "dHJhdmVsLXNhbXBsZQ=="
  "memory_allocation_in_mb" = 200
  "name" = "travel-sample"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  "replicas" = 1
  "stats" = {
    "disk_used_in_mib" = 0
    "item_count" = 163
    "memory_used_in_mib" = 33
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "time_to_live_in_seconds" = 0
  "type" = "couchbase"
}
samplebucket_id = "dHJhdmVsLXNhbXBsZQ=="
samplebuckets_list = {
  "cluster_id" = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
  "data" = tolist([
    {
      "bucket_conflict_resolution" = "seqno"
      "cluster_id" = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      "durability_level" = "none"
      "eviction_policy" = "fullEviction"
      "flush" = false
      "id" = "Z2FtZXNpbS1zYW1wbGU="
      "memory_allocation_in_mb" = 200
      "name" = "gamesim-sample"
      "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      "replicas" = 1
      "stats" = {
        "disk_used_in_mib" = 0
        "item_count" = 390
        "memory_used_in_mib" = 42
        "ops_per_second" = 0
      }
      "storage_backend" = "couchstore"
      "time_to_live_in_seconds" = 0
      "type" = "couchbase"
    },
  ])
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
}
```

# DESTROY
### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:

```
➜  sample_bucket git:(AV-70846_add_import_sample_data_set_apis) ✗ terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in $HOME/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_sample_buckets.existing_sample_buckets: Reading...
couchbase-capella_sample_bucket.new_sample_bucket: Refreshing state... [id=dHJhdmVsLXNhbXBsZQ==]
data.couchbase-capella_sample_buckets.existing_sample_buckets: Read complete after 1s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_sample_bucket.new_sample_bucket will be destroyed
  - resource "couchbase-capella_sample_bucket" "new_sample_bucket" {
      - bucket_conflict_resolution = "seqno" -> null
      - cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025" -> null
      - durability_level           = "none" -> null
      - eviction_policy            = "fullEviction" -> null
      - flush                      = false -> null
      - id                         = "dHJhdmVsLXNhbXBsZQ==" -> null
      - memory_allocation_in_mb    = 200 -> null
      - name                       = "travel-sample" -> null
      - organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3" -> null
      - replicas                   = 1 -> null
      - stats                      = {
          - disk_used_in_mib   = 120 -> null
          - item_count         = 63288 -> null
          - memory_used_in_mib = 165 -> null
          - ops_per_second     = 0 -> null
        } -> null
      - storage_backend            = "couchstore" -> null
      - time_to_live_in_seconds    = 0 -> null
      - type                       = "couchbase" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_sample_bucket   = {
      - bucket_conflict_resolution = "seqno"
      - cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      - durability_level           = "none"
      - eviction_policy            = "fullEviction"
      - flush                      = false
      - id                         = "dHJhdmVsLXNhbXBsZQ=="
      - memory_allocation_in_mb    = 200
      - name                       = "travel-sample"
      - organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
      - replicas                   = 1
      - stats                      = {
          - disk_used_in_mib   = 0
          - item_count         = 163
          - memory_used_in_mib = 33
          - ops_per_second     = 0
        }
      - storage_backend            = "couchstore"
      - time_to_live_in_seconds    = 0
      - type                       = "couchbase"
    } -> null
  - samplebucket_id    = "dHJhdmVsLXNhbXBsZQ==" -> null
  - samplebuckets_list = {
      - cluster_id      = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
      - data            = [
          - {
              - bucket_conflict_resolution = "seqno"
              - cluster_id                 = "17619f3c-08f5-40a3-8c0c-d2e5b263a025"
              - durability_level           = "none"
              - eviction_policy            = "fullEviction"
              - flush                      = false
              - id                         = "dHJhdmVsLXNhbXBsZQ=="
              - memory_allocation_in_mb    = 200
              - name                       = "travel-sample"
              - organization_id            = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
              - project_id                 = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
              - replicas                   = 1
              - stats                      = {
                  - disk_used_in_mib   = 120
                  - item_count         = 63288
                  - memory_used_in_mib = 165
                  - ops_per_second     = 0
                }
              - storage_backend            = "couchstore"
              - time_to_live_in_seconds    = 0
              - type                       = "couchbase"
            },
        ]
      - organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

couchbase-capella_sample_bucket.new_sample_bucket: Destroying... [id=dHJhdmVsLXNhbXBsZQ==]
couchbase-capella_sample_bucket.new_sample_bucket: Destruction complete after 2s

Destroy complete! Resources: 1 destroyed.
```