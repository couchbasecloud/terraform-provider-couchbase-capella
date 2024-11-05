output "list_indexes" {
  value = data.couchbase-capella_query_indexes.list
}

data "couchbase-capella_query_indexes" "list" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_name     = var.bucket_name
  scope_name      = var.scope_name
  collection_name = var.collection_name
}