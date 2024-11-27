# Capella Create Deferred Secondary Index Example

This example shows how to create deferred secondary indexes on a Capella cluster.

This creates 5 deferred indexes in test.test.test.  Deferred indexes are best suited for large indexes, or when creating many indexes.
When in doubt, use deferred indexes.

If creating indexes in bulk, it is strongly recommended to build indexes in batches of 100 or less.

To try this example, configure your Couchbase Capella provider as described in README in the root of this project.

## CREATE INDEXES
**Ensure you create `indexes.json`, which is a file with indexes.**
**Place this file in the same directory as the terraform script.**

Command: `terraform apply`

Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published
│ releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.couchbase-capella_query_index_monitor.mon_indexes will be read during apply
  # (depends on a resource or a module with changes pending)
 <= data "couchbase-capella_query_index_monitor" "mon_indexes" {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + indexes         = [
          + "idx1",
          + "idx2",
          + "idx3",
          + "idx4",
          + "idx5",
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
    }

  # couchbase-capella_query_indexes.build_idx will be created
  + resource "couchbase-capella_query_indexes" "build_idx" {
      + bucket_name     = "test"
      + build_indexes   = [
          + "idx1",
          + "idx2",
          + "idx3",
          + "idx4",
          + "idx5",
        ]
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
    }

  # couchbase-capella_query_indexes.new_indexes["idx1"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + index_keys      = [
          + "field1",
        ]
      + index_name      = "idx1"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
      + with            = {
          + defer_build = true
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx2"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + index_keys      = [
          + "field2",
        ]
      + index_name      = "idx2"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
      + with            = {
          + defer_build = true
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx3"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + index_keys      = [
          + "field3",
        ]
      + index_name      = "idx3"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
      + with            = {
          + defer_build = true
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx4"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + index_keys      = [
          + "field5",
        ]
      + index_name      = "idx4"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
      + with            = {
          + defer_build = true
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx5"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + index_keys      = [
          + "field5",
        ]
      + index_name      = "idx5"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
      + with            = {
          + defer_build = true
        }
    }

Plan: 6 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + watch_indexes = {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + indexes         = [
          + "idx1",
          + "idx2",
          + "idx3",
          + "idx4",
          + "idx5",
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.new_indexes["idx1"]: Creating...
couchbase-capella_query_indexes.build_idx: Creating...
couchbase-capella_query_indexes.new_indexes["idx2"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx5"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx3"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx4"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx1"]: Creation complete after 1s
couchbase-capella_query_indexes.new_indexes["idx2"]: Creation complete after 2s
couchbase-capella_query_indexes.new_indexes["idx5"]: Creation complete after 3s
couchbase-capella_query_indexes.new_indexes["idx4"]: Creation complete after 4s
couchbase-capella_query_indexes.new_indexes["idx3"]: Creation complete after 5s
couchbase-capella_query_indexes.build_idx: Still creating... [10s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [20s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [30s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [40s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [50s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [1m0s elapsed]
couchbase-capella_query_indexes.build_idx: Creation complete after 1m1s
data.couchbase-capella_query_index_monitor.mon_indexes: Reading...
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [10s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [20s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [30s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [40s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [50s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [1m0s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Read complete after 1m1s

Apply complete! Resources: 6 added, 0 changed, 0 destroyed.

Outputs:

watch_indexes = {
  "bucket_name" = "test"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "test"
  "indexes" = toset([
    "idx1",
    "idx2",
    "idx3",
    "idx4",
    "idx5",
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "test"
}
```
