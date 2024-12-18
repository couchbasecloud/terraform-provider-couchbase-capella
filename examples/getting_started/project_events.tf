# Stores the project events details in an output variable.
# Can be viewed using `terraform output existing_project_events` command
output "existing_project_events" {
  value = data.couchbase-capella_project_events.existing_project_events
}

data "couchbase-capella_project_events" "existing_project_events" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
}