output "existing_events" {
  value = data.couchbase-capella_events.existing_events
}

data "couchbase-capella_events" "existing_events" {
  organization_id = var.organization_id
}

# Example of using optional fields when fetching events. Use this if you want to use optional fields.

#  data "couchbase-capella_events" "existing_events" {
#    organization_id = var.organization_id
#    project_ids     = var.events.project_ids
#    cluster_ids     = var.events.cluster_ids
#    user_ids        = var.events.user_ids
#    severity_levels = var.events.severity_levels
#    tags            = var.events.tags
#    from            = var.events.from
#    to              = var.events.to
#    sort_by         = var.events.sort_by
#    sort_direction  = var.events.sort_direction
#    page            = var.events.page
#    per_page        = var.events.per_page
#  }