data "couchbase-capella_app_service_audit_log_export_job" "existing_app_service_audit_log_export_job" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  id              = var.export_id
}

output "existing_app_service_audit_log_export_job" {
  value = data.couchbase-capella_app_service_audit_log_export_job.existing_app_service_audit_log_export_job
}
