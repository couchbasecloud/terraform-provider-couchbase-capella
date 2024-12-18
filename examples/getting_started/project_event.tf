# Stores the project event details in an output variable.
# Can be viewed using `terraform output existing_project_event` command
output "existing_project_event" {
  value = data.couchbase-capella_project_event.existing_project_event
}

data "couchbase-capella_project_event" "existing_project_event" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  id              = data.couchbase-capella_project_events.existing_project_events.data[0].id
  depends_on      = [data.couchbase-capella_project_events.existing_project_events]
}