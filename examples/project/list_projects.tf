output "projects_list" {
  value = data.capella_projects.existing_projects
}

data "capella_projects" "existing_projects" {
  organization_id = var.organization_id
}
