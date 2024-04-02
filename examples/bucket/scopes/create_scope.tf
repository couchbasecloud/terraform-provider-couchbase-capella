output "new_scope" {
  value = couchbase-capella_scope.new_scope
}

resource "couchbase-capella_scope" "new_scope" {
  scope_name      = var.scope
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
}
