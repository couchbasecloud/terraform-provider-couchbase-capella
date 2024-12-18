resource "couchbase-capella_query_indexes" "idx" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id

  bucket_name     = couchbase-capella_bucket.new_bucket.name
  scope_name      = couchbase-capella_scope.new_scope.scope_name
  collection_name = couchbase-capella_collection.new_collection.collection_name

  index_name = var.index_name
  index_keys = var.index_keys
  where      = var.where

  with = {
    defer_build   = var.with.defer_build
    num_replica   = var.with.num_replica
    num_partition = var.with.num_partition
  }
}

output "idx" {
  value = couchbase-capella_query_indexes.idx
}