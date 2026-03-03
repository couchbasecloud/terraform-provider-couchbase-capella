resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  output_type = var.output_type
  credentials = var.credentials
}
