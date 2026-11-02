resource "couchbase-capella_app_service_audit_log_export_job" "new_app_service_audit_log_export_job" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
  start           = "2026-11-01T00:00:00Z"
  end             = "2026-11-02T00:00:00Z"
}
