data "couchbase-capella_app_service_audit_log_export_job" "existing_app_service_audit_log_export_job" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
  id              = "<export_id>"
}
