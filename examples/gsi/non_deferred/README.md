# Capella Create Secondary Index Example

This example shows how to create secondary indexes on a Capella cluster.

This creates 1 non-deferred indexes in api.metrics.memory.  Non-deferred indexes are best suited when creating 1 small index.

To run configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new index using `create_index.tf` file.
2. LIST: List indexees using `list_indexes.tf` file.
3. DELETE: Delete the indexes.
5. IMPORT: Import an index to the state file.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE INDEX

Command: `terraform apply`

Sample Output:
```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_query_indexes.idx will be created
  + resource "couchbase-capella_query_indexes" "idx" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "ram",
        ]
      + index_name      = "idx1"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + num_replica = (known after apply)
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + idx = {
      + bucket_name     = "api"
      + build_indexes   = null
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "ram",
        ]
      + index_name      = "idx1"
      + is_primary      = null
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + partition_by    = null
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + where           = null
      + with            = {
          + defer_build   = null
          + num_partition = null
          + num_replica   = (known after apply)
        }
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.idx: Creating...
couchbase-capella_query_indexes.idx: Creation complete after 4s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

idx = {
  "bucket_name" = "api"
  "build_indexes" = toset(null) /* of string */
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "memory"
  "index_keys" = tolist([
    "ram",
  ])
  "index_name" = "idx1"
  "is_primary" = tobool(null)
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "partition_by" = tolist(null) /* of string */
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "metrics"
  "status" = "Ready"
  "where" = tostring(null)
  "with" = {
    "defer_build" = tobool(null)
    "num_partition" = tonumber(null)
    "num_replica" = 1
  }
}
```

## LIST INDEXES IN api.metrics.memory

Command: `terraform plan`
Sample Output:
```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_query_indexes.list: Reading...
couchbase-capella_query_indexes.idx: Refreshing state...
data.couchbase-capella_query_indexes.list: Read complete after 0s

Changes to Outputs:
  - idx          = {
      - bucket_name     = "api"
      - build_indexes   = null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - collection_name = "memory"
      - index_keys      = [
          - "ram",
        ]
      - index_name      = "idx1"
      - is_primary      = null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      - partition_by    = null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      - scope_name      = "metrics"
      - status          = "Ready"
      - where           = null
      - with            = {
          - defer_build   = null
          - num_partition = null
          - num_replica   = 1
        }
    } -> null
  + list_indexes = {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + data            = [
          + {
              + definition = "CREATE INDEX `idx1` ON `api`.`metrics`.`memory`(`ram`) WITH {  \"nodes\":[ \"node1.acme.com:18091\",\"node2.acme.com:18091\" ], \"num_replica\":1 }"
              + index_name = "idx1 (replica 1)"
            },
          + {
              + definition = "CREATE INDEX `idx1` ON `api`.`metrics`.`memory`(`ram`) WITH {  \"nodes\":[ \"node1.acme.com:18091\",\"node2.acme.com:18091\" ], \"num_replica\":1 }"
              + index_name = "idx1"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

## DELETE INDEX
### Remove resource block from script

Command: `terraform apply`
Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible
│ with published releases.
╵
couchbase-capella_query_indexes.idx: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_query_indexes.idx will be destroyed
  # (because couchbase-capella_query_indexes.idx is not in configuration)
  - resource "couchbase-capella_query_indexes" "idx" {
      - bucket_name     = "api" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "memory" -> null
      - index_keys      = [
          - "ram",
        ] -> null
      - index_name      = "idx1" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "metrics" -> null
      - with            = {
          - num_replica = 1 -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.idx: Destroying...
couchbase-capella_query_indexes.idx: Destruction complete after 0s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.
```


## IMPORT INDEX

Command: `import couchbase-capella_query_indexes.idx1 index_name=idx1,collection_name=memory,scope_name=metrics,bucket_name=api,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
couchbase-capella_query_indexes.idx1: Importing from ID "index_name=idx1,collection_name=memory,scope_name=metrics,bucket_name=api,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_query_indexes.idx1: Import prepared!
  Prepared couchbase-capella_query_indexes for import
couchbase-capella_query_indexes.idx1: Refreshing state...
2024-12-03T10:36:19.939-0800 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.idx1 during refresh.
      - .cluster_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .scope_name: was null, but now cty.StringVal("metrics")
      - .where: was null, but now cty.StringVal("")
      - .bucket_name: was null, but now cty.StringVal("api")
      - .index_name: was cty.StringVal("index_name=idx1,collection_name=memory,scope_name=metrics,bucket_name=api,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"), but now cty.StringVal("idx1")
      - .organization_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .project_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .status: was null, but now cty.StringVal("")
      - .index_keys: was null, but now cty.ListVal([]cty.Value{cty.StringVal("`ram`")})
      - .is_primary: was null, but now cty.False
      - .collection_name: was null, but now cty.StringVal("memory")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

