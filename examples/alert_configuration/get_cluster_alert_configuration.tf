output "existing_cluster_alert_configuration" {
  value = data.couchbase-capella_cluster_alert_configuration.existing_cluster_alert_configuration
}

data "couchbase-capella_cluster_alert_configuration" "existing_cluster_alert_configuration" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
