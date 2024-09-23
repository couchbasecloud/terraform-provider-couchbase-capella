# Stores the event details in an output variable.
# Can be viewed using `terraform output existing_event` command
output "existing_event" {
  value = data.couchbase-capella_event.existing_event
}

data "couchbase-capella_event" "existing_event" {
  organization_id = var.organization_id
  id              = data.couchbase-capella_events.existing_events.data[0].id
  depends_on      = [data.couchbase-capella_events.existing_events]
}