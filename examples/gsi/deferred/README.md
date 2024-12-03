# Capella Create Secondary Index Example

This example shows how to create secondary indexes on a Capella cluster.

This creates 5 deferred indexes in api.metrics.memory  

Deferred indexes are best suited for large indexes (like partitioned indexes), or when creating many indexes. When in doubt, use deferred indexes :)

If creating indexes in bulk, it is strongly recommended to build indexes in batches of 100 or less.

To run configure your Couchbase Capella provider as described in README in the root of this project.

## CREATE INDEXES
### Ensure you create indexes.json which is a file with indexes
### Place this file in the same directory as the terraform script.

Command: `terraform apply`

Sample Output:
```
export TF_LOG=INFO

terraform apply
2024-12-03T11:25:28.928-0800 [INFO]  Terraform version: 1.6.6
2024-12-03T11:25:28.929-0800 [INFO]  Go runtime version: go1.21.5
2024-12-03T11:25:28.929-0800 [INFO]  CLI args: []string{"terraform", "apply"}
2024-12-03T11:25:28.930-0800 [INFO]  Loading CLI configuration from /Users/hiteshwalia/.terraformrc
2024-12-03T11:25:28.936-0800 [INFO]  CLI command args: []string{"apply"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2024-12-03T11:25:28.974-0800 [INFO]  backend/local: starting Apply operation
2024-12-03T11:25:28.997-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:25:29.117-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:25:29.116-0800
2024-12-03T11:25:29.207-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:25:29.228-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:25:29.228-0800
2024-12-03T11:25:29.268-0800 [INFO]  backend/local: apply calling Plan
2024-12-03T11:25:29.271-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:25:29.293-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:25:29.292-0800
2024-12-03T11:25:29.315-0800 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=c09711c5-9ab0-820f-fd76-ef31f576af78 tf_rpc=ConfigureProvider @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/provider/provider.go:72 timestamp=2024-12-03T11:25:29.315-0800
2024-12-03T11:25:29.315-0800 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/provider/provider.go:159 authentication_token="***" host=http://localhost:8084 success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=c09711c5-9ab0-820f-fd76-ef31f576af78 @module=couchbase_capella tf_rpc=ConfigureProvider timestamp=2024-12-03T11:25:29.315-0800

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.couchbase-capella_query_index_monitor.mon_indexes will be read during apply
  # (depends on a resource or a module with changes pending)
 <= data "couchbase-capella_query_index_monitor" "mon_indexes" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + indexes         = [
          + "idx1",
          + "idx2",
          + "idx3",
          + "idx4",
          + "idx5",
        ]
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
    }

  # couchbase-capella_query_indexes.build_idx will be created
  + resource "couchbase-capella_query_indexes" "build_idx" {
      + bucket_name     = "api"
      + build_indexes   = [
          + "idx1",
          + "idx2",
          + "idx3",
          + "idx4",
          + "idx5",
        ]
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + num_replica = (known after apply)
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx1"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "field1",
        ]
      + index_name      = "idx1"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + defer_build = true
          + num_replica = (known after apply)
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx2"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "field2",
        ]
      + index_name      = "idx2"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + defer_build = true
          + num_replica = (known after apply)
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx3"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "field3",
        ]
      + index_name      = "idx3"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + defer_build = true
          + num_replica = (known after apply)
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx4"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "field5",
        ]
      + index_name      = "idx4"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + defer_build = true
          + num_replica = (known after apply)
        }
    }

  # couchbase-capella_query_indexes.new_indexes["idx5"] will be created
  + resource "couchbase-capella_query_indexes" "new_indexes" {
      + bucket_name     = "api"
      + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + collection_name = "memory"
      + index_keys      = [
          + "field5",
        ]
      + index_name      = "idx5"
      + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
      + scope_name      = "metrics"
      + status          = (known after apply)
      + with            = {
          + defer_build = true
          + num_replica = (known after apply)
        }
    }

Plan: 6 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + new_indexes = {
      + idx1 = {
          + bucket_name     = "api"
          + build_indexes   = null
          + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + collection_name = "memory"
          + index_keys      = [
              + "field1",
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
              + defer_build   = true
              + num_partition = null
              + num_replica   = (known after apply)
            }
        }
      + idx2 = {
          + bucket_name     = "api"
          + build_indexes   = null
          + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + collection_name = "memory"
          + index_keys      = [
              + "field2",
            ]
          + index_name      = "idx2"
          + is_primary      = null
          + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
          + partition_by    = null
          + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + scope_name      = "metrics"
          + status          = (known after apply)
          + where           = null
          + with            = {
              + defer_build   = true
              + num_partition = null
              + num_replica   = (known after apply)
            }
        }
      + idx3 = {
          + bucket_name     = "api"
          + build_indexes   = null
          + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + collection_name = "memory"
          + index_keys      = [
              + "field3",
            ]
          + index_name      = "idx3"
          + is_primary      = null
          + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
          + partition_by    = null
          + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + scope_name      = "metrics"
          + status          = (known after apply)
          + where           = null
          + with            = {
              + defer_build   = true
              + num_partition = null
              + num_replica   = (known after apply)
            }
        }
      + idx4 = {
          + bucket_name     = "api"
          + build_indexes   = null
          + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + collection_name = "memory"
          + index_keys      = [
              + "field5",
            ]
          + index_name      = "idx4"
          + is_primary      = null
          + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
          + partition_by    = null
          + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + scope_name      = "metrics"
          + status          = (known after apply)
          + where           = null
          + with            = {
              + defer_build   = true
              + num_partition = null
              + num_replica   = (known after apply)
            }
        }
      + idx5 = {
          + bucket_name     = "api"
          + build_indexes   = null
          + cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + collection_name = "memory"
          + index_keys      = [
              + "field5",
            ]
          + index_name      = "idx5"
          + is_primary      = null
          + organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
          + partition_by    = null
          + project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
          + scope_name      = "metrics"
          + status          = (known after apply)
          + where           = null
          + with            = {
              + defer_build   = true
              + num_partition = null
              + num_replica   = (known after apply)
            }
        }
    }

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

2024-12-03T11:25:32.721-0800 [INFO]  backend/local: apply calling Apply
2024-12-03T11:25:32.726-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:25:32.755-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:25:32.755-0800
2024-12-03T11:25:32.780-0800 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_rpc=ConfigureProvider @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/provider/provider.go:72 @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=f60ef08e-1b63-f03d-1a0f-e20288fc97cd timestamp=2024-12-03T11:25:32.780-0800
2024-12-03T11:25:32.780-0800 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: tf_rpc=ConfigureProvider @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/provider/provider.go:159 @module=couchbase_capella host=http://localhost:8084 success=true tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella authentication_token="***" tf_req_id=f60ef08e-1b63-f03d-1a0f-e20288fc97cd timestamp=2024-12-03T11:25:32.780-0800
couchbase-capella_query_indexes.new_indexes["idx3"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx2"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx1"]: Creating...
couchbase-capella_query_indexes.new_indexes["idx5"]: Creating...
2024-12-03T11:25:32.796-0800 [INFO]  Starting apply for couchbase-capella_query_indexes.new_indexes["idx2"]
couchbase-capella_query_indexes.new_indexes["idx4"]: Creating...
2024-12-03T11:25:32.796-0800 [INFO]  Starting apply for couchbase-capella_query_indexes.new_indexes["idx3"]
2024-12-03T11:25:32.796-0800 [INFO]  Starting apply for couchbase-capella_query_indexes.new_indexes["idx4"]
2024-12-03T11:25:32.797-0800 [INFO]  Starting apply for couchbase-capella_query_indexes.new_indexes["idx5"]
2024-12-03T11:25:32.797-0800 [INFO]  Starting apply for couchbase-capella_query_indexes.new_indexes["idx1"]
couchbase-capella_query_indexes.new_indexes["idx1"]: Creation complete after 5s
couchbase-capella_query_indexes.new_indexes["idx4"]: Creation complete after 6s
couchbase-capella_query_indexes.new_indexes["idx5"]: Creation complete after 7s
couchbase-capella_query_indexes.new_indexes["idx2"]: Creation complete after 8s
couchbase-capella_query_indexes.new_indexes["idx3"]: Creation complete after 9s
couchbase-capella_query_indexes.build_idx: Creating...
2024-12-03T11:25:42.082-0800 [INFO]  Starting apply for couchbase-capella_query_indexes.build_idx
couchbase-capella_query_indexes.build_idx: Still creating... [10s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [20s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [30s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [40s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [50s elapsed]
couchbase-capella_query_indexes.build_idx: Still creating... [1m0s elapsed]
couchbase-capella_query_indexes.build_idx: Creation complete after 1m5s
data.couchbase-capella_query_index_monitor.mon_indexes: Reading...
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [10s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [20s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [30s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [40s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [50s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [1m0s elapsed]
2024-12-03T11:27:51.693-0800 [INFO]  provider.terraform-provider-couchbase-capella: All indexes are ready. Please run "terraform apply --refresh-only" to update state.: @module=couchbase_capella tf_data_source_type=couchbase-capella_query_index_monitor tf_req_id=fffb0e2c-46be-5f59-0b3f-60db1f7a4fd5 @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/datasources/gsi_monitor.go:137 tf_rpc=ReadDataSource tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp=2024-12-03T11:27:51.691-0800
data.couchbase-capella_query_index_monitor.mon_indexes: Read complete after 1m5s

Apply complete! Resources: 6 added, 0 changed, 0 destroyed.

Outputs:

new_indexes = {
  "idx1" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field1",
    ])
    "index_name" = "idx1"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Created"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx2" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field2",
    ])
    "index_name" = "idx2"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Created"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx3" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field3",
    ])
    "index_name" = "idx3"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Created"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx4" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field5",
    ])
    "index_name" = "idx4"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Created"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx5" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field5",
    ])
    "index_name" = "idx5"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Created"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
}
```

### Logs show indexes ready: 
### 2024-12-03T11:27:51.693-0800 [INFO]  provider.terraform-provider-couchbase-capella: All indexes are ready. Please run "terraform apply --refresh-only" to update state.: @module=couchbase_capella tf_data_source_type=couchbase-capella_query_index_monitor tf_req_id=fffb0e2c-46be-5f59-0b3f-60db1f7a4fd5 @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/datasources/gsi_monitor.go:137 tf_rpc=ReadDataSource tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp=2024-12-03T11:27:51.691-0800

### The output shows indexes are in Created state.  This is correct because that is the state before sending the build statement.
### After the build the indexes will be in Ready state.  Please refresh the state on terraform.

Command: `terraform apply --refresh-only`

Sample Output:
```
2024-12-03T11:39:09.191-0800 [INFO]  Terraform version: 1.6.6
2024-12-03T11:39:09.191-0800 [INFO]  Go runtime version: go1.21.5
2024-12-03T11:39:09.191-0800 [INFO]  CLI args: []string{"terraform", "apply", "--refresh-only"}
2024-12-03T11:39:09.192-0800 [INFO]  Loading CLI configuration from /Users/hiteshwalia/.terraformrc
2024-12-03T11:39:09.196-0800 [INFO]  CLI command args: []string{"apply", "--refresh-only"}
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - couchbasecloud/couchbase-capella in /Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
2024-12-03T11:39:09.228-0800 [INFO]  backend/local: starting Apply operation
2024-12-03T11:39:09.242-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:39:09.322-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:39:09.321-0800
2024-12-03T11:39:09.395-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:39:09.415-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:39:09.415-0800
2024-12-03T11:39:09.448-0800 [INFO]  backend/local: apply calling Plan
2024-12-03T11:39:09.451-0800 [INFO]  provider: configuring client automatic mTLS
2024-12-03T11:39:09.475-0800 [INFO]  provider.terraform-provider-couchbase-capella: configuring server automatic mTLS: timestamp=2024-12-03T11:39:09.475-0800
2024-12-03T11:39:09.502-0800 [INFO]  provider.terraform-provider-couchbase-capella: Configuring the Capella Client: tf_req_id=3cd02dd7-581b-0e39-1166-06425747081a tf_rpc=ConfigureProvider @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/provider/provider.go:72 @module=couchbase_capella tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp=2024-12-03T11:39:09.501-0800
2024-12-03T11:39:09.502-0800 [INFO]  provider.terraform-provider-couchbase-capella: Configured Capella client: @module=couchbase_capella authentication_token="***" tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella tf_req_id=3cd02dd7-581b-0e39-1166-06425747081a tf_rpc=ConfigureProvider @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/provider/provider.go:159 host=http://localhost:8084 success=true timestamp=2024-12-03T11:39:09.502-0800
couchbase-capella_query_indexes.new_indexes["idx5"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx2"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx1"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx4"]: Refreshing state...
couchbase-capella_query_indexes.new_indexes["idx3"]: Refreshing state...
2024-12-03T11:39:09.932-0800 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.new_indexes["idx3"] during refresh.
      - .status: was cty.StringVal("Created"), but now cty.StringVal("Ready")
2024-12-03T11:39:10.748-0800 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.new_indexes["idx4"] during refresh.
      - .status: was cty.StringVal("Created"), but now cty.StringVal("Ready")
2024-12-03T11:39:11.755-0800 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.new_indexes["idx5"] during refresh.
      - .status: was cty.StringVal("Created"), but now cty.StringVal("Ready")
2024-12-03T11:39:12.744-0800 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.new_indexes["idx1"] during refresh.
      - .status: was cty.StringVal("Created"), but now cty.StringVal("Ready")
2024-12-03T11:39:13.741-0800 [WARN]  Provider "registry.terraform.io/couchbasecloud/couchbase-capella" produced an unexpected new value for couchbase-capella_query_indexes.new_indexes["idx2"] during refresh.
      - .status: was cty.StringVal("Created"), but now cty.StringVal("Ready")
couchbase-capella_query_indexes.build_idx: Refreshing state...
data.couchbase-capella_query_index_monitor.mon_indexes: Reading...
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [10s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [20s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [30s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [40s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [50s elapsed]
data.couchbase-capella_query_index_monitor.mon_indexes: Still reading... [1m0s elapsed]
2024-12-03T11:40:18.052-0800 [INFO]  provider.terraform-provider-couchbase-capella: All indexes are ready. Please run "terraform apply --refresh-only" to update state.: tf_rpc=ReadDataSource @caller=/Users/hiteshwalia/GolandProjects/terraform-provider-couchbase-capella/internal/datasources/gsi_monitor.go:137 tf_req_id=e71d13cb-78bb-25e4-6a09-5347c71fd047 @module=couchbase_capella tf_data_source_type=couchbase-capella_query_index_monitor tf_provider_addr=hashicorp.com/couchbasecloud/couchbase-capella timestamp=2024-12-03T11:40:18.043-0800
data.couchbase-capella_query_index_monitor.mon_indexes: Read complete after 1m4s

Note: Objects have changed outside of Terraform

Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:

  # couchbase-capella_query_indexes.new_indexes["idx1"] has changed
  ~ resource "couchbase-capella_query_indexes" "new_indexes" {
      ~ status          = "Created" -> "Ready"
        # (9 unchanged attributes hidden)
    }

  # couchbase-capella_query_indexes.new_indexes["idx2"] has changed
  ~ resource "couchbase-capella_query_indexes" "new_indexes" {
      ~ status          = "Created" -> "Ready"
        # (9 unchanged attributes hidden)
    }

  # couchbase-capella_query_indexes.new_indexes["idx3"] has changed
  ~ resource "couchbase-capella_query_indexes" "new_indexes" {
      ~ status          = "Created" -> "Ready"
        # (9 unchanged attributes hidden)
    }

  # couchbase-capella_query_indexes.new_indexes["idx4"] has changed
  ~ resource "couchbase-capella_query_indexes" "new_indexes" {
      ~ status          = "Created" -> "Ready"
        # (9 unchanged attributes hidden)
    }

  # couchbase-capella_query_indexes.new_indexes["idx5"] has changed
  ~ resource "couchbase-capella_query_indexes" "new_indexes" {
      ~ status          = "Created" -> "Ready"
        # (9 unchanged attributes hidden)
    }


This is a refresh-only plan, so Terraform will not take any actions to undo these. If you were expecting these changes then you can apply this plan to record the updated
values in the Terraform state without changing any remote objects.

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Changes to Outputs:
  ~ new_indexes = {
      ~ idx1 = {
          ~ status          = "Created" -> "Ready"
            # (13 unchanged attributes hidden)
        }
      ~ idx2 = {
          ~ status          = "Created" -> "Ready"
            # (13 unchanged attributes hidden)
        }
      ~ idx3 = {
          ~ status          = "Created" -> "Ready"
            # (13 unchanged attributes hidden)
        }
      ~ idx4 = {
          ~ status          = "Created" -> "Ready"
            # (13 unchanged attributes hidden)
        }
      ~ idx5 = {
          ~ status          = "Created" -> "Ready"
            # (13 unchanged attributes hidden)
        }
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Would you like to update the Terraform state to reflect these detected changes?
  Terraform will write these changes to the state without modifying any real infrastructure.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

2024-12-03T11:40:22.959-0800 [INFO]  backend/local: apply calling Apply

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

new_indexes = {
  "idx1" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field1",
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
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx2" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field2",
    ])
    "index_name" = "idx2"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Ready"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx3" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field3",
    ])
    "index_name" = "idx3"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Ready"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx4" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field5",
    ])
    "index_name" = "idx4"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Ready"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
  "idx5" = {
    "bucket_name" = "api"
    "build_indexes" = toset(null) /* of string */
    "cluster_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "collection_name" = "memory"
    "index_keys" = tolist([
      "field5",
    ])
    "index_name" = "idx5"
    "is_primary" = tobool(null)
    "organization_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "partition_by" = tolist(null) /* of string */
    "project_id" = "ffffffff-aaaa-1414-eeee-000000000000"
    "scope_name" = "metrics"
    "status" = "Ready"
    "where" = tostring(null)
    "with" = {
      "defer_build" = true
      "num_partition" = tonumber(null)
      "num_replica" = 1
    }
  }
}
```