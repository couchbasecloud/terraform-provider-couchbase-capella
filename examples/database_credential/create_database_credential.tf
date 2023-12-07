output "new_database_credential" {
  value     = couchbase-capella_database_credential.new_database_credential
  sensitive = true
}

output "database_credential_id" {
  value = couchbase-capella_database_credential.new_database_credential.id
}

resource "couchbase-capella_database_credential" "new_database_credential" {
  name            = var.database_credential.database_credential_name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  password        = var.database_credential.password
  access          = var.access
}

