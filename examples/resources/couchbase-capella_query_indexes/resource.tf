# Example json file -
  {
  "resource": {
    "couchbase-capella_indexes": {
      "idx1": {
        "organization_id": "ffffffff-aaaa-1414-eeee-000000000000",
        "project_id": "ffffffff-aaaa-1414-eeee-000000000000",
        "cluster_id": "ffffffff-aaaa-1414-eeee-000000000000",
        "bucket_name": "test",
        "scope_name": "metrics",
        "scope_name": "memory",
        "index_name": "idx1",
        "index_keys": [
          "field1"
        ],
        "with": {
          "defer_build": true
        }
      },
    "idx2": {
          "organization_id": "ffffffff-aaaa-1414-eeee-000000000000",
          "project_id": "ffffffff-aaaa-1414-eeee-000000000000",
          "cluster_id": "ffffffff-aaaa-1414-eeee-000000000000",
          "bucket_name": "test",
          "scope_name": "metrics",
          "scope_name": "memory",
          "index_name": "idx2",
          "index_keys": [
            "field2"
          ],
          "with": {
            "defer_build": true
             }
          }
        }
      }
    }

# Example for deferred index build -
locals {
  index_template = templatefile("${path.module}/indexes.json", {
    organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
    project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  })

  decoded_template = jsondecode(local.index_template)
  index_names      = [for idx, details in local.decoded_template.resource["couchbase-capella_indexes"] : details.index_name]
}

resource "couchbase-capella_query_indexes" "new_indexes" {
  for_each        = jsondecode(local.index_template).resource["couchbase-capella_indexes"]
  organization_id = each.value.organization_id
  project_id      = each.value.project_id
  cluster_id      = each.value.cluster_id
  bucket_name     = each.value.bucket_name
  scope_name      = each.value.scope_name
  collection_name = each.value.collection_name
  index_name      = each.value.index_name
  index_keys      = each.value.index_keys
  with = {
    defer_build = each.value.with.defer_build
  }
}

resource "couchbase-capella_query_indexes" "build_idx" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"

  bucket_name     = "api"
  scope_name      = "metrics"
  collection_name = "memory"

  build_indexes = local.index_names

  depends_on = [couchbase-capella_query_indexes.new_indexes]
}

data "couchbase-capella_query_index_monitor" "mon_indexes" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  bucket_name     = "api"
  scope_name      = "metrics"
  collection_name = "memory"
  indexes         = local.index_names

  depends_on = [couchbase-capella_query_indexes.build_idx]
}


# Example for non-deferred index build -
resource "couchbase-capella_query_indexes" "idx" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  bucket_name     = "api"
  scope_name      = "metrics"
  collection_name = "memory"
  index_name      = "idx1"
  index_keys      = ["ram"]
}