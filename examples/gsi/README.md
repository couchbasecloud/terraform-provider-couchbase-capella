# Capella Create Secondary Index Example

This example shows how to create a secondary index on Capella.

This creates index idx_pe9 in test.test.test

To run configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Create a new index using `create_index.tf` file.
2. UPDATE: Update the number of replicas to 2.
3. LIST: List indexees using `list_indexes.tf` file.
4. DELETE: Delete the index.
5. IMPORT: Import an index to the state file.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE INDEX

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
data.couchbase-capella_network_peers.existing_network_peers: Reading...
data.couchbase-capella_network_peers.existing_network_peers: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_query_indexes.idx will be created
  + resource "couchbase-capella_query_indexes" "idx" {
      + bucket_name     = "test"
      + cluster_id      = "73b9e1ad-2619-4b0c-8a5d-25183c96c670"
      + collection_name = "test"
      + index_keys      = [
          + "sourceairport",
          + "destinationairport",
          + "stops",
          + "airline",
          + "id",
          + "ARRAY_COUNT(schedule)",
        ]
      + index_name      = "idx_pe9"
      + organization_id = "313e335d-cb58-4989-a9fb-ae983a1055e3"
      + partition_by    = [
          + "sourceairport",
          + "destinationairport",
        ]
      + project_id      = "42dcf8f9-3b14-44f4-9a71-a50de45175ea"
      + scope_name      = "test"
      + with            = {
          + defer_build   = false
          + num_partition = 8
          + num_replica   = 1
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.idx: Creating...
couchbase-capella_query_indexes.idx: Creation complete after 5s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

## UPDATE INDEX WITH 2 REPLICAS

Command: `terraform apply`

Sample Output:
```
couchbase-capella_query_indexes.idx: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # couchbase-capella_query_indexes.idx will be updated in-place
  ~ resource "couchbase-capella_query_indexes" "idx" {
      ~ with            = {
          ~ num_replica   = 0 -> 2
            # (2 unchanged attributes hidden)
        }
        # (9 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.idx: Modifying...
couchbase-capella_query_indexes.idx: Modifications complete after 0s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## LIST INDEXES IN test.test.test

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
data.couchbase-capella_query_indexes.list: Read complete after 0s

Changes to Outputs:
  + list_indexes = {
      + bucket_name     = "test"
      + cluster_id      = "73b9e1ad-2619-4b0c-8a5d-25183c96c670"
      + collection_name = "test"
      + data            = [
          + {
              + definition = "CREATE INDEX `idx1` ON `test`.`test`.`test`(`c0`) WITH {  \"nodes\":[ \"svc-dqi-node-001.fiflsazo9gjxjbrt.aws-guardians.nonprod-project-avengers.com:18091\",\"svc-dqi-node-003.fiflsazo9gjxjbrt.aws-guardians.nonprod-project-avengers.com:18091\" ], \"num_replica\":1 }"
              + index_name = "idx1 (replica 1)"
            },
          + {
              + definition = "CREATE INDEX `idx1` ON `test`.`test`.`test`(`c0`) WITH {  \"nodes\":[ \"svc-dqi-node-001.fiflsazo9gjxjbrt.aws-guardians.nonprod-project-avengers.com:18091\",\"svc-dqi-node-003.fiflsazo9gjxjbrt.aws-guardians.nonprod-project-avengers.com:18091\" ], \"num_replica\":1 }"
              + index_name = "idx1"
            },
        ]
      + organization_id = "313e335d-cb58-4989-a9fb-ae983a1055e3"
      + project_id      = "42dcf8f9-3b14-44f4-9a71-a50de45175ea"
      + scope_name      = "test"
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

## DELETE INDEX idx_pe9
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
couchbase-capella_query_indexes.idx: Refreshing state...

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # couchbase-capella_query_indexes.idx will be destroyed
  # (because couchbase-capella_query_indexes.idx is not in configuration)
  - resource "couchbase-capella_query_indexes" "idx" {
      - bucket_name     = "test" -> null
      - cluster_id      = "73b9e1ad-2619-4b0c-8a5d-25183c96c670" -> null
      - collection_name = "test" -> null
      - index_keys      = [
          - "sourceairport",
          - "destinationairport",
          - "stops",
          - "airline",
          - "id",
          - "ARRAY_COUNT(schedule)",
        ] -> null
      - index_name      = "idx_pe9" -> null
      - organization_id = "313e335d-cb58-4989-a9fb-ae983a1055e3" -> null
      - partition_by    = [
          - "sourceairport",
          - "destinationairport",
        ] -> null
      - project_id      = "42dcf8f9-3b14-44f4-9a71-a50de45175ea" -> null
      - scope_name      = "test" -> null
      - with            = {
          - defer_build   = false -> null
          - num_partition = 8 -> null
          - num_replica   = 0 -> null
        } -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

couchbase-capella_query_indexes.idx: Destroying...
couchbase-capella_query_indexes.idx: Destruction complete after 1s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.
```


## IMPORT INDEX

Command: `terraform import couchbase-capella_query_indexes.idx1 index_name=idx1,collection_name=test,scope_name=test,bucket_name=test,cluster_id=73b9e1ad-2619-4b0c-8a5d-25183c96c670,project_id=42dcf8f9-3b14-44f4-9a71-a50de45175ea,or
ganization_id=313e335d-cb58-4989-a9fb-ae983a1055e3`

Sample Output:
```
couchbase-capella_query_indexes.idx1: Importing from ID "index_name=idx1,collection_name=test,scope_name=test,bucket_name=test,cluster_id=73b9e1ad-2619-4b0c-8a5d-25183c96c670,project_id=42dcf8f9-3b14-44f4-9a71-a50de45175ea,organization_id=313e335d-cb58-4989-a9fb-ae983a1055e3"...
couchbase-capella_query_indexes.idx1: Import prepared!
  Prepared couchbase-capella_query_indexes for import
couchbase-capella_query_indexes.idx1: Refreshing state...
2024-09-25T12:43:09.513-0700 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.idx1 during refresh.
      - .collection_name: was null, but now cty.StringVal("test")
      - .index_keys: was null, but now cty.ListVal([]cty.Value{cty.StringVal("`c0`")})
      - .index_name: was cty.StringVal("index_name=idx1,collection_name=test,scope_name=test,bucket_name=test,cluster_id=73b9e1ad-2619-4b0c-8a5d-25183c96c670,project_id=42dcf8f9-3b14-44f4-9a71-a50de45175ea,organization_id=313e335d-cb58-4989-a9fb-ae983a1055e3"), but now cty.StringVal("idx1")
      - .is_primary: was null, but now cty.False
      - .project_id: was null, but now cty.StringVal("42dcf8f9-3b14-44f4-9a71-a50de45175ea")
      - .scope_name: was null, but now cty.StringVal("test")
      - .where: was null, but now cty.StringVal("")
      - .bucket_name: was null, but now cty.StringVal("test")
      - .organization_id: was null, but now cty.StringVal("313e335d-cb58-4989-a9fb-ae983a1055e3")
      - .cluster_id: was null, but now cty.StringVal("73b9e1ad-2619-4b0c-8a5d-25183c96c670")

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

