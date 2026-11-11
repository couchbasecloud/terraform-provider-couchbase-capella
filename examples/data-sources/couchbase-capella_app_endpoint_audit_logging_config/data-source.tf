data "couchbase-capella_app_endpoint_audit_logging_config" "existing_app_endpoint_audit_logging_config" {
  organization_id   = "<organization_id>"
  project_id        = "<project_id>"
  cluster_id        = "<cluster_id>"
  app_service_id    = "<app_service_id>"
  app_endpoint_name = "<app_endpoint_name>"
}
