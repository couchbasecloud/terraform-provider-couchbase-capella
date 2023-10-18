output "database_credentials_list" {
  value = data.capella_database_credentials.existing_credentials
}

data "capella_database_credentials" "existing_credentials" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
