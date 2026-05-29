# Retrieve a single project by ID
data "couchbase-capella_project" "existing_project" {
  organization_id = "<organization_id>"
  id              = "<project_id>"
}

output "existing_project" {
  value = data.couchbase-capella_project.existing_project
}

# List all projects in an organization
data "couchbase-capella_projects" "existing_projects" {
  organization_id = "<organization_id>"
}

output "existing_projects" {
  value = data.couchbase-capella_projects.existing_projects
}