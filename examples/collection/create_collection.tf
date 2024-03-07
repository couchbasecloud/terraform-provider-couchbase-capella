output "new_collection" {
  value = couchbase-capella_collection.new_collection
}

resource "couchbase-capella_collection" "new_collection" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
  scope_name      = var.scope_name
  collection_name = var.collection.collection_name
  max_ttl         = var.collection.max_ttl
}