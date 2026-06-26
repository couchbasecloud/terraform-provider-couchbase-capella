output "database_roles_list" {
  value = data.couchbase-capella_database_roles.existing_roles
}

data "couchbase-capella_database_roles" "existing_roles" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

