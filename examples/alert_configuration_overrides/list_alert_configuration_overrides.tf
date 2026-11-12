output "existing_alert_configuration_overrides" {
  value = data.couchbase-capella_alert_configuration_overrides.existing_alert_configuration_overrides
}

data "couchbase-capella_alert_configuration_overrides" "existing_alert_configuration_overrides" {
  organization_id = var.organization_id
}
