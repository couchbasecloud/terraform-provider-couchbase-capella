resource "couchbase-capella_app_service_audit_log_export_job" "new_app_service_audit_log_export_job" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  start           = var.start
  end             = var.end
}

output "new_app_service_audit_log_export_job" {
  value = couchbase-capella_app_service_audit_log_export_job.new_app_service_audit_log_export_job
}
