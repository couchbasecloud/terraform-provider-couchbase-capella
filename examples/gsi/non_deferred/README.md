# Capella Create Secondary Index Example

This example shows how to create secondary indexes on a Capella cluster.

This creates 5 non-deferred indexes in test.test.test.  Non-deferred indexes are best suited when creating small indexes.

To run configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new index using `create_indexes.tf` file.
2. LIST: List indexees using `list_indexes.tf` file.
3. DELETE: Delete the indexes.
5. IMPORT: Import an index to the state file.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE INDEXES
### Ensure you create indexes.json which is a file with indexes
### Place this file in the same directory as the terraform script.

Command: `terraform apply`

Sample Output:
```
$ terraform apply
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

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
          + defer_build   = false
          + num_partition = 8
          + num_replica   = 0
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
          + defer_build   = false
          + num_partition = 8
          + num_replica   = 0
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
          + defer_build   = false
          + num_partition = 8
          + num_replica   = 0
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
          + defer_build   = false
          + num_partition = 8
          + num_replica   = 0
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
          + defer_build   = false
          + num_partition = 8
          + num_replica   = 0
        }
    }

Plan: 5 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.new_indexes["idx2"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx5"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx3"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx1"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx4"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx3"]: Creation complete after 4s
couchbase-capella_query_indexes.new_indexes["idx2"]: Creation complete after 9s
couchbase-capella_query_indexes.new_indexes["idx5"]: Still creating... [10s elapsed]
couchbase-capella_query_indexes.new_indexes["idx1"]: Still creating... [10s elapsed]
couchbase-capella_query_indexes.new_indexes["idx4"]: Still creating... [10s elapsed]
couchbase-capella_query_indexes.new_indexes["idx1"]: Creation complete after 14s
couchbase-capella_query_indexes.new_indexes["idx4"]: Creation complete after 19s
couchbase-capella_query_indexes.new_indexes["idx5"]: Still creating... [20s elapsed]
couchbase-capella_query_indexes.new_indexes["idx5"]: Creation complete after 24s

Apply complete! Resources: 5 added, 0 changed, 0 destroyed.
```

## LIST INDEXES IN test.test.test

Command: `terraform plan`
Sample Output:
```
terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/$USER/GolandProjects/terraform-provider-couchbase-capella/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_query_indexes.list: Reading...
couchbase-capella_query_indexes.new_indexes["idx5"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx2"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx3"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx1"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx4"]: Refreshing state...
data.couchbase-capella_query_indexes.list: Read complete after 1s

Changes to Outputs:
  + list_indexes = {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + data            = [
          + {
              + definition = "CREATE INDEX `idx1` ON `test`.`test`.`test`(`field1`)"
              + index_name = "idx1"
            },
          + {
              + definition = "CREATE INDEX `idx2` ON `test`.`test`.`test`(`field2`)"
              + index_name = "idx2"
            },
          + {
              + definition = "CREATE INDEX `idx3` ON `test`.`test`.`test`(`field3`)"
              + index_name = "idx3"
            },
          + {
              + definition = "CREATE INDEX `idx4` ON `test`.`test`.`test`(`field5`)"
              + index_name = "idx4"
            },
          + {
              + definition = "CREATE INDEX `idx5` ON `test`.`test`.`test`(`field5`)"
              + index_name = "idx5"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

## DELETE INDEXES
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
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
data.couchbase-capella_query_indexes.list: Reading...
couchbase-capella_query_indexes.new_indexes["idx5"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx3"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx4"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx2"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx1"]: Refreshing state...
data.couchbase-capella_query_indexes.list: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_query_indexes.new_indexes["idx1"] will be destroyed
  # (because couchbase-capella_query_indexes.new_indexes is not in configuration)
  - resource "couchbase-capella_query_indexes" "new_indexes" {
      - bucket_name     = "test" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "test" -> null
      - index_keys      = [
          - "field1",
        ] -> null
      - index_name      = "idx1" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "test" -> null
      - with            = {
          - defer_build   = false -> null
          - num_partition = 8 -> null
          - num_replica   = 0 -> null
        } -> null
    }

  # couchbase-capella_query_indexes.new_indexes["idx2"] will be destroyed
  # (because couchbase-capella_query_indexes.new_indexes is not in configuration)
  - resource "couchbase-capella_query_indexes" "new_indexes" {
      - bucket_name     = "test" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "test" -> null
      - index_keys      = [
          - "field2",
        ] -> null
      - index_name      = "idx2" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "test" -> null
      - with            = {
          - defer_build   = false -> null
          - num_partition = 8 -> null
          - num_replica   = 0 -> null
        } -> null
    }

  # couchbase-capella_query_indexes.new_indexes["idx3"] will be destroyed
  # (because couchbase-capella_query_indexes.new_indexes is not in configuration)
  - resource "couchbase-capella_query_indexes" "new_indexes" {
      - bucket_name     = "test" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "test" -> null
      - index_keys      = [
          - "field3",
        ] -> null
      - index_name      = "idx3" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "test" -> null
      - with            = {
          - defer_build   = false -> null
          - num_partition = 8 -> null
          - num_replica   = 0 -> null
        } -> null
    }

  # couchbase-capella_query_indexes.new_indexes["idx4"] will be destroyed
  # (because couchbase-capella_query_indexes.new_indexes is not in configuration)
  - resource "couchbase-capella_query_indexes" "new_indexes" {
      - bucket_name     = "test" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "test" -> null
      - index_keys      = [
          - "field5",
        ] -> null
      - index_name      = "idx4" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "test" -> null
      - with            = {
          - defer_build   = false -> null
          - num_partition = 8 -> null
          - num_replica   = 0 -> null
        } -> null
    }

  # couchbase-capella_query_indexes.new_indexes["idx5"] will be destroyed
  # (because couchbase-capella_query_indexes.new_indexes is not in configuration)
  - resource "couchbase-capella_query_indexes" "new_indexes" {
      - bucket_name     = "test" -> null
      - cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - collection_name = "test" -> null
      - index_keys      = [
          - "field5",
        ] -> null
      - index_name      = "idx5" -> null
      - organization_id = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - project_id      = "ffffffff-aaaa-1414-eeee-000000000000" -> null
      - scope_name      = "test" -> null
      - with            = {
          - defer_build   = false -> null
          - num_partition = 8 -> null
          - num_replica   = 0 -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 5 to destroy.

Changes to Outputs:
  + list_indexes = {
      + bucket_name     = "test"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "test"
      + data            = [
          + {
              + definition = "CREATE INDEX `idx1` ON `test`.`test`.`test`(`field1`)"
              + index_name = "idx1"
            },
          + {
              + definition = "CREATE INDEX `idx2` ON `test`.`test`.`test`(`field2`)"
              + index_name = "idx2"
            },
          + {
              + definition = "CREATE INDEX `idx3` ON `test`.`test`.`test`(`field3`)"
              + index_name = "idx3"
            },
          + {
              + definition = "CREATE INDEX `idx4` ON `test`.`test`.`test`(`field5`)"
              + index_name = "idx4"
            },
          + {
              + definition = "CREATE INDEX `idx5` ON `test`.`test`.`test`(`field5`)"
              + index_name = "idx5"
            },
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "test"
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.new_indexes["idx1"]: Destroying...
couchbase-capella_query_indexes.new_indexes["idx3"]: Destroying...
couchbase-capella_query_indexes.new_indexes["idx5"]: Destroying...
couchbase-capella_query_indexes.new_indexes["idx2"]: Destroying...
couchbase-capella_query_indexes.new_indexes["idx4"]: Destroying...
couchbase-capella_query_indexes.new_indexes["idx5"]: Destruction complete after 0s
couchbase-capella_query_indexes.new_indexes["idx3"]: Destruction complete after 0s
couchbase-capella_query_indexes.new_indexes["idx2"]: Destruction complete after 0s
couchbase-capella_query_indexes.new_indexes["idx1"]: Destruction complete after 0s
couchbase-capella_query_indexes.new_indexes["idx4"]: Destruction complete after 0s

Apply complete! Resources: 0 added, 0 changed, 5 destroyed.

Outputs:

list_indexes = {
  "bucket_name" = "test"
  "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "collection_name" = "test"
  "data" = toset([
    {
      "definition" = "CREATE INDEX `idx1` ON `test`.`test`.`test`(`field1`)"
      "index_name" = "idx1"
    },
    {
      "definition" = "CREATE INDEX `idx2` ON `test`.`test`.`test`(`field2`)"
      "index_name" = "idx2"
    },
    {
      "definition" = "CREATE INDEX `idx3` ON `test`.`test`.`test`(`field3`)"
      "index_name" = "idx3"
    },
    {
      "definition" = "CREATE INDEX `idx4` ON `test`.`test`.`test`(`field5`)"
      "index_name" = "idx4"
    },
    {
      "definition" = "CREATE INDEX `idx5` ON `test`.`test`.`test`(`field5`)"
      "index_name" = "idx5"
    },
  ])
  "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
  "scope_name" = "test"
}
```


## IMPORT INDEX

Command: `terraform import couchbase-capella_query_indexes.idx1 index_name=idx1,collection_name=test,scope_name=test,bucket_name=test,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,or
ganization_id=ffffffff-aaaa-1414-eeee-000000000000`

Sample Output:
```
couchbase-capella_query_indexes.idx1: Importing from ID "index_name=idx1,collection_name=test,scope_name=test,bucket_name=test,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"...
couchbase-capella_query_indexes.idx1: Import prepared!
  Prepared couchbase-capella_query_indexes for import
couchbase-capella_query_indexes.idx1: Refreshing state...
2024-09-25T12:43:09.513-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.idx1 during refresh.
      - .collection_name: was null, but now cty.StringVal("test")
      - .index_keys: was null, but now cty.ListVal([]cty.Value{cty.StringVal("`c0`")})
      - .index_name: was cty.StringVal("index_name=idx1,collection_name=test,scope_name=test,bucket_name=test,cluster_id=ffffffff-aaaa-1414-eeee-000000000000,project_id=ffffffff-aaaa-1414-eeee-000000000000,organization_id=ffffffff-aaaa-1414-eeee-000000000000"), but now cty.StringVal("idx1")
      - .is_primary: was null, but now cty.False
      - .project_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .scope_name: was null, but now cty.StringVal("test")
      - .where: was null, but now cty.StringVal("")
      - .bucket_name: was null, but now cty.StringVal("test")
      - .organization_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")
      - .cluster_id: was null, but now cty.StringVal("ffffffff-aaaa-1414-eeee-000000000000")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

