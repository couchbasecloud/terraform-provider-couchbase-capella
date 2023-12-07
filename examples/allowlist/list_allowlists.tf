output "allowlists_list" {
  value = data.couchbase-capella_allowlists.existing_allowlists
}

data "couchbase-capella_allowlists" "existing_allowlists" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
