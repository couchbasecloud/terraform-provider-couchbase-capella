output "scopes_list" {
  value = data.couchbase-capella_scopes.existing_scopes
}

data "couchbase-capella_scopes" "existing_scopes" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
}