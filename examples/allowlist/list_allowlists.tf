output "projects_list" {
  value = data.capella_allowlist.existing_allowlists
}

data "capella_allowlist" "existing_allowlists" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
