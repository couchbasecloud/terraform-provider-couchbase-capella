data "couchbase-capella_app_endpoint_audit_logging_config" "existing_app_endpoint_audit_logging_config" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint_name = var.app_endpoint_name
}

output "existing_app_endpoint_audit_logging_config" {
  value = data.couchbase-capella_app_endpoint_audit_logging_config.existing_app_endpoint_audit_logging_config
}
