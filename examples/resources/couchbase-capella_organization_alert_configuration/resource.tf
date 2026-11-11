resource "couchbase-capella_organization_alert_configuration" "example" {
  organization_id = "<organization_id>"

  configs = {
    cpu_usage_high = {
      frequency_hours    = 6
      warning_threshold  = 85
      critical_threshold = 95
    }
  }

  recipients = ["platform@example.com"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "2027-01-01T00:00:00Z"
    }
  }
}
