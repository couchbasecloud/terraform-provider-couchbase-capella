output "existing_alerts" {
  value = data.couchbase-capella_alerts.existing_alerts
}

data "couchbase-capella_alerts" "existing_alerts" {
  organization_id = var.organization_id
}

# Example of using optional fields when fetching alerts. Use this if you want to use optional fields.

#  data "couchbase-capella_alerts" "existing_alerts" {
#    organization_id = var.organization_id
#    status          = var.alerts.status
#    cluster_ids     = var.alerts.cluster_ids
#    project_ids     = var.alerts.project_ids
#    severity_levels = var.alerts.severity_levels
#    from            = var.alerts.from
#    to              = var.alerts.to
#  }
