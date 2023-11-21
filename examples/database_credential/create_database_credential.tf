output "new_database_credential" {
  value     = capella_database_credential.new_database_credential
  sensitive = true
}

output "database_credential_id" {
  value = capella_database_credential.new_database_credential.id
}

resource "capella_database_credential" "new_database_credential" {
  name            = var.database_credential_name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  password        = var.password
  access          = var.access
}

