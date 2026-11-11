output "existing_app_service_maintenance_window" {
  value = data.couchbase-capella_app_service_maintenance_window.existing_app_service_maintenance_window
}

data "couchbase-capella_app_service_maintenance_window" "existing_app_service_maintenance_window" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}
