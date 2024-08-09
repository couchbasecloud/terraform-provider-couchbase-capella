output "existing_event" {
  value = data.couchbase-capella_event.existing_event
}

data "couchbase-capella_event" "existing_event" {
  organization_id = var.organization_id
  id              = var.event_id
}