resource "couchbase-capella_query_indexes" "idx" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  bucket_name     = var.bucket_name
  scope_name      = var.scope_name
  collection_name = var.collection_name

  index_name = var.index_name
  index_keys = var.index_keys
}

output "idx" {
  value = couchbase-capella_query_indexes.idx
}
