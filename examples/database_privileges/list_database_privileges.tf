data "couchbase-capella_database_privileges" "list_database_privileges" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

output "database_privileges" {
  value = data.couchbase-capella_database_privileges.list_database_privileges
}

