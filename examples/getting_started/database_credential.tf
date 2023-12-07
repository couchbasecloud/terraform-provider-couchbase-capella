# Stores the database_credential details in an output variable.
# Can be viewed using `terraform output database_credential` command
output "database_credential" {
  value     = couchbase-capella_database_credential.new_database_credential
  sensitive = true
}

resource "couchbase-capella_database_credential" "new_database_credential" {
  name            = var.database_credential_name
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  password        = var.password
  access          = var.access
}

