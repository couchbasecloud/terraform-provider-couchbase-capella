---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "couchbase-capella_query_indexes Data Source - terraform-provider-couchbase-capella"
subcategory: ""
description: |-
  The data source for retrieving Query Indexes in Couchbase Capella.
---

# couchbase-capella_query_indexes (Data Source)

The data source for retrieving Query Indexes in Couchbase Capella.

## Example Usage

```terraform
data "couchbase-capella_query_indexes" "list" {
  organization_id = <organization_id>
  project_id      = <project_id>
  cluster_id      = <cluster_id>
  bucket_name     = "api"
  scope_name      = "metrics"
  collection_name = "memory"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bucket_name` (String) The name of the bucket where the indexes exist. Specifies the bucket portion of the keyspace.
- `cluster_id` (String) The GUID4 ID of the cluster.
- `organization_id` (String) The GUID4 ID of the organization.
- `project_id` (String) The GUID4 ID of the project.

### Optional

- `collection_name` (String) Specifies the collection portion of the keyspace. If unspecified, this will be the default collection.
- `scope_name` (String) The name of the scope where the indexes exist. Specifies the scope portion of the keyspace. If unspecified, this will be the default scope.

### Read-Only

- `data` (Attributes List) List of indexes in the specified keyspace. (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `condition` (String) The WHERE clause condition for the index.
- `index_key` (List of String) List of document fields being indexed.
- `is_primary` (Boolean) Specifies whether this is a primary index.
- `keyspace_id` (String) The full keyspace identifier for the index (bucket.scope.collection).
- `name` (String) The name of the index.
- `partition` (List of String) List of fields the index is partitioned by.
- `partition_count` (Number) Number of partitions for the index.
- `replica_count` (Number) Number of index replicas.
- `state` (String) The current state of the index. For example 'Created', 'Ready', etc.
