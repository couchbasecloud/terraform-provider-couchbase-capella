# Stores the events details in an output variable.
# Can be viewed using `terraform output existing_events` command
output "existing_events" {
  value = data.couchbase-capella_events.existing_events
}

data "couchbase-capella_events" "existing_events" {
  organization_id = var.organization_id
}