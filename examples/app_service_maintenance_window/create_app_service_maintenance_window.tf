output "new_app_service_maintenance_window" {
  value = couchbase-capella_app_service_maintenance_window.new_app_service_maintenance_window
}

resource "couchbase-capella_app_service_maintenance_window" "new_app_service_maintenance_window" {
  organization_id  = var.organization_id
  project_id       = var.project_id
  cluster_id       = var.cluster_id
  app_service_id   = var.app_service_id
  preferred_day    = var.maintenance_window.preferred_day
  preferred_hour   = var.maintenance_window.preferred_hour
  preferred_minute = var.maintenance_window.preferred_minute
}
