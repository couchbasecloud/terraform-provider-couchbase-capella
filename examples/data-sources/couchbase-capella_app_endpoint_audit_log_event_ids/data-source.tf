data "couchbase-capella_app_endpoint_audit_log_event_ids" "existing_app_endpoint_audit_log_event_ids" {
  organization_id   = "<organization_id>"
  project_id        = "<project_id>"
  cluster_id        = "<cluster_id>"
  app_service_id    = "<app_service_id>"
  app_endpoint_name = "<app_endpoint_name>"
}
