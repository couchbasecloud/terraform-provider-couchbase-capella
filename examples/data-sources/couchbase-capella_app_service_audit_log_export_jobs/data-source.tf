data "couchbase-capella_app_service_audit_log_export_jobs" "app_service_audit_log_export_jobs" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  # Optional. Omit to return every export job of the App Service.
  export_id = "<export_id>"
}
