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

resource "couchbase-capella_query_indexes" "build_idx" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  bucket_name     = "api"
  scope_name      = "metrics"
  collection_name = "memory"
  build_indexes = = ["idx1","idx2", "idx3", "idx4", "idx5",]
}


# couchbase-capella_query_indexes.new_indexes["idx1"] will be created -
 resource "couchbase-capella_query_indexes" "new_indexes" {
    bucket_name     = "api"
    cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    collection_name = "memory"
    index_keys      = ["field1"]
    index_name      = "idx1"
    organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
    project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
    scope_name      = "metrics"
    with            = {
      defer_build = true
        }
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