# Capella Bucket Example

This example shows how to create and manage Buckets in Capella.

This creates a new bucket in the selected Capella cluster. It uses the organization ID, projectId and clusterId to do so.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new bucket in an existing Capella cluster as stated in the `create_bucket.tf` file.

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
terraform plan    
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + cluster_id            = "96f2e933-cf5e-407a-b9c7-926f706f89ef"
      + conflict_resolution   = "seqno"
      + durability_level      = "majorityAndPersistActive"
      + eviction_policy       = "fullEviction"
      + flush                 = true
      + id                    = (known after apply)
      + memory_allocationinmb = 105
      + name                  = "test_bucket"
      + organization_id       = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id            = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + replicas              = 2
      + stats                 = (known after apply)
      + storage_backend       = "couchstore"
      + ttl                   = 100
      + type                  = "couchbase"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_bucket = {
      + cluster_id            = "96f2e933-cf5e-407a-b9c7-926f706f89ef"
      + conflict_resolution   = "seqno"
      + durability_level      = "majorityAndPersistActive"
      + eviction_policy       = "fullEviction"
      + flush                 = true
      + id                    = (known after apply)
      + memory_allocationinmb = 105
      + name                  = "test_bucket"
      + organization_id       = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id            = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + replicas              = 2
      + stats                 = (known after apply)
      + storage_backend       = "couchstore"
      + ttl                   = 100
      + type                  = "couchbase"
    }

─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

### Apply the Plan, in order to create a new Bucket

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be created
  + resource "capella_bucket" "new_bucket" {
      + cluster_id            = "96f2e933-cf5e-407a-b9c7-926f706f89ef"
      + conflict_resolution   = "seqno"
      + durability_level      = "majorityAndPersistActive"
      + eviction_policy       = "fullEviction"
      + flush                 = true
      + id                    = (known after apply)
      + memory_allocationinmb = 105
      + name                  = "test_bucket"
      + organization_id       = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id            = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + replicas              = 2
      + stats                 = (known after apply)
      + storage_backend       = "couchstore"
      + ttl                   = 100
      + type                  = "couchbase"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_bucket = {
      + cluster_id            = "96f2e933-cf5e-407a-b9c7-926f706f89ef"
      + conflict_resolution   = "seqno"
      + durability_level      = "majorityAndPersistActive"
      + eviction_policy       = "fullEviction"
      + flush                 = true
      + id                    = (known after apply)
      + memory_allocationinmb = 105
      + name                  = "test_bucket"
      + organization_id       = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      + project_id            = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      + replicas              = 2
      + stats                 = (known after apply)
      + storage_backend       = "couchstore"
      + ttl                   = 100
      + type                  = "couchbase"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

capella_bucket.new_bucket: Creating...
capella_bucket.new_bucket: Still creating... [10s elapsed]
capella_bucket.new_bucket: Creation complete after 13s [id=dGVzdF9idWNrZXQ=]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

new_bucket = {
  "cluster_id" = "96f2e933-cf5e-407a-b9c7-926f706f89ef"
  "conflict_resolution" = "seqno"
  "durability_level" = "majorityAndPersistActive"
  "eviction_policy" = "fullEviction"
  "flush" = true
  "id" = "dGVzdF9idWNrZXQ="
  "memory_allocationinmb" = 105
  "name" = "test_bucket"
  "organization_id" = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  "project_id" = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
  "replicas" = 2
  "stats" = {
    "disk_used_in_mib" = 0
    "item_count" = 0
    "memory_used_in_mib" = 0
    "ops_per_second" = 0
  }
  "storage_backend" = "couchstore"
  "ttl" = 100
  "type" = "couchbase"
}
```

### Finally, destroy the resources created by Terraform

Command: `terraform destroy`

Sample Output:

```
terraform destroy
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/couchabasecloud/capella in /Users/nidhi.kumar/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵
capella_bucket.new_bucket: Refreshing state... [id=dGVzdF9idWNrZXQ=]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # capella_bucket.new_bucket will be destroyed
  - resource "capella_bucket" "new_bucket" {
      - cluster_id            = "96f2e933-cf5e-407a-b9c7-926f706f89ef" -> null
      - conflict_resolution   = "seqno" -> null
      - durability_level      = "majorityAndPersistActive" -> null
      - eviction_policy       = "fullEviction" -> null
      - flush                 = true -> null
      - id                    = "dGVzdF9idWNrZXQ=" -> null
      - memory_allocationinmb = 105 -> null
      - name                  = "test_bucket" -> null
      - organization_id       = "6af08c0a-8cab-4c1c-b257-b521575c16d0" -> null
      - project_id            = "f14134f2-7943-4e7b-b2c5-fc2071728b6e" -> null
      - replicas              = 2 -> null
      - stats                 = {
          - disk_used_in_mib   = 0 -> null
          - item_count         = 0 -> null
          - memory_used_in_mib = 0 -> null
          - ops_per_second     = 0 -> null
        } -> null
      - storage_backend       = "couchstore" -> null
      - ttl                   = 100 -> null
      - type                  = "couchbase" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - new_bucket = {
      - cluster_id            = "96f2e933-cf5e-407a-b9c7-926f706f89ef"
      - conflict_resolution   = "seqno"
      - durability_level      = "majorityAndPersistActive"
      - eviction_policy       = "fullEviction"
      - flush                 = true
      - id                    = "dGVzdF9idWNrZXQ="
      - memory_allocationinmb = 105
      - name                  = "test_bucket"
      - organization_id       = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
      - project_id            = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
      - replicas              = 2
      - stats                 = {
          - disk_used_in_mib   = 0
          - item_count         = 0
          - memory_used_in_mib = 0
          - ops_per_second     = 0
        }
      - storage_backend       = "couchstore"
      - ttl                   = 100
      - type                  = "couchbase"
    } -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

capella_bucket.new_bucket: Destroying... [id=dGVzdF9idWNrZXQ=]
capella_bucket.new_bucket: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```

