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

output "new_advanced_database_credential" {
  value     = couchbase-capella_database_credential.new_advanced_database_credential
  sensitive = true
}

output "advanced_database_credential_id" {
  value = couchbase-capella_database_credential.new_advanced_database_credential.id
}

# An advanced database credential is assigned capella user roles for
# fine-grained RBAC access instead of bucket-level access permissions.
# The user roles must already exist in the cluster.
resource "couchbase-capella_database_credential" "new_advanced_database_credential" {
  name            = var.advanced_database_credential.database_credential_name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  password        = var.advanced_database_credential.password
  credential_type = "advanced"
  user_roles      = var.user_roles
}

