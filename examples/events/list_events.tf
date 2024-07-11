output "existing_events" {
  value = data.couchbase-capella_events.existing_events
}

data "couchbase-capella_events" "existing_events" {
  organization_id = var.organization_id
}