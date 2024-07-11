output "existing_project_events" {
  value = data.couchbase-capella_project_events.existing_project_events
}

data "couchbase-capella_project_events" "existing_project_events" {
  organization_id = var.organization_id
  project_id = var.project_id
}