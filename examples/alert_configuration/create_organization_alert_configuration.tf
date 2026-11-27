output "new_organization_alert_configuration" {
  value = couchbase-capella_organization_alert_configuration.new_organization_alert_configuration
}

resource "couchbase-capella_organization_alert_configuration" "new_organization_alert_configuration" {
  organization_id = var.organization_id

  configs = {
    cpu_usage_high = {
      frequency_hours = 6
    }
  }

  recipients = ["platform@couchbase.com"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "2027-01-01T00:00:00Z"
    }
  }
}
