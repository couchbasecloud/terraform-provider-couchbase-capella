resource "couchbase-capella_project_alert_configuration" "example" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"

  configs = {
    disk_usage_high = {
      warning_threshold  = 80
      critical_threshold = 90
    }
  }

  recipients = ["project-team@couchbase.com"]

  muted_alerts = {
    cpu_usage_high = {
      action      = "mute"
      muted_until = "2027-01-01T00:00:00Z"
    }
  }
}
