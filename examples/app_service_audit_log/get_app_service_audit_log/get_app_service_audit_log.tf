data "couchbase-capella_app_service_audit_log" "existing_app_service_audit_log" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}

output "existing_app_service_audit_log" {
  value = data.couchbase-capella_app_service_audit_log.existing_app_service_audit_log
}
