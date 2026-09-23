output "new_database_role" {
  value = couchbase-capella_database_role.new_database_role
}

output "database_role_id" {
  value = couchbase-capella_database_role.new_database_role.id
}

resource "couchbase-capella_database_role" "new_database_role" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = var.database_role.name
  description     = var.database_role.description
  access          = var.access
}

