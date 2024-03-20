# Stores the scope name in an output variable.
# Can be viewed using `terraform output scope` command
output "scope" {
  value = couchbase-capella_scope.new_scope
}

resource "couchbase-capella_scope" "new_scope" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  bucket_id       = couchbase-capella_bucket.new_bucket.id
  scope_name      = var.scope.scope_name
}