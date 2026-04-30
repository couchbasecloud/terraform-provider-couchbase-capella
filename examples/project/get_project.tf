output "existing_project" {
  value = data.couchbase-capella_project.existing_project
}

data "couchbase-capella_project" "existing_project" {
  organization_id = var.organization_id
  id              = var.project_id
}
