data "couchbase-capella_app_service_audit_log_export_jobs" "list_app_service_audit_log_export_jobs" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  # Optional. Omit to return every export job of the App Service.
  export_id = var.export_id
}

output "app_service_audit_log_export_jobs" {
  value = data.couchbase-capella_app_service_audit_log_export_jobs.list_app_service_audit_log_export_jobs
}
