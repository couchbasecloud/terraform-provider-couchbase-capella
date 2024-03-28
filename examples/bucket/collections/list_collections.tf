output "collections_list" {
  value = data.couchbase-capella_collections.existing_collections
}

data "couchbase-capella_collections" "existing_collections" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
  scope_name      = var.scope_name
}
