output "existing_project_alert_configuration" {
  value = data.couchbase-capella_project_alert_configuration.existing_project_alert_configuration
}

data "couchbase-capella_project_alert_configuration" "existing_project_alert_configuration" {
  organization_id = var.organization_id
  project_id      = var.project_id
}
