resource "couchbase-capella_query_indexes" "idx1" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id

  bucket_name     = couchbase-capella_bucket.new_bucket.name
  scope_name      = couchbase-capella_scope.new_scope.scope_name
  collection_name = couchbase-capella_collection.new_collection.collection_name

  index_name = "idx1"
  index_keys = var.index_keys
  where      = var.where

  with = {
    defer_build = var.with.defer_build
  }
}

resource "couchbase-capella_query_indexes" "idx2" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id

  bucket_name     = couchbase-capella_bucket.new_bucket.name
  scope_name      = couchbase-capella_scope.new_scope.scope_name
  collection_name = couchbase-capella_collection.new_collection.collection_name

  index_name = "idx2"
  index_keys = var.index_keys
  where      = var.where

  with = {
    defer_build = var.with.defer_build
  }
}

resource "couchbase-capella_query_indexes" "idx3" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id

  bucket_name     = couchbase-capella_bucket.new_bucket.name
  scope_name      = couchbase-capella_scope.new_scope.scope_name
  collection_name = couchbase-capella_collection.new_collection.collection_name

  index_name = "idx3"
  index_keys = var.index_keys
  where      = var.where

  with = {
    defer_build = var.with.defer_build
  }
}