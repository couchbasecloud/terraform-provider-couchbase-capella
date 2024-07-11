output "existing_project_event" {
  value = data.couchbase-capella_project_event.existing_project_event
}

data "couchbase-capella_project_event" "existing_project_event" {
  organization_id = var.organization_id
  id              = var.event_id
  project_id      = var.project_id
}