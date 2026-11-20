output "existing_organization_alert_configuration" {
  value = data.couchbase-capella_organization_alert_configuration.existing_organization_alert_configuration
}

data "couchbase-capella_organization_alert_configuration" "existing_organization_alert_configuration" {
  organization_id = var.organization_id
}
