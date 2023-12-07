output "projects_list" {
  value = data.couchbase-capella_projects.existing_projects
}

data "couchbase-capella_projects" "existing_projects" {
  organization_id = var.organization_id
}
