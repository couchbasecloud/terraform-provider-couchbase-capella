output "new_cluster_alert_configuration" {
  value = couchbase-capella_cluster_alert_configuration.new_cluster_alert_configuration
}

resource "couchbase-capella_cluster_alert_configuration" "new_cluster_alert_configuration" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  configs = {
    cpu_usage_high = {
      frequency_hours    = 3
      warning_threshold  = 80
      critical_threshold = 90
    }
    CBHardOutOfMemoryErrors = {
      frequency_hours = 1
    }
  }

  recipients = ["oncall@example.com"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "2027-01-01T00:00:00Z"
    }
  }
}
