data "couchbase-capella_project" "existing_project" {
  organization_id = "<organization_id>"
  id              = "<project_id>"
}

output "existing_project" {
  value = data.couchbase-capella_project.existing_project
}