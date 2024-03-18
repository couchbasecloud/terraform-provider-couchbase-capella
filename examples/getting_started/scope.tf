# Stores the scope name in an output variable.
# Can be viewed using `terraform output scope` command
output "scope" {
  value = couchbase-capella_scope.new_scope
}

resource "couchbase-capella_scope" "new_scope" {
  scope_name      = var.scope.scope_name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = var.bucket_id
}