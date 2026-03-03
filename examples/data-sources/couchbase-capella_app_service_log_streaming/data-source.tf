data "couchbase-capella_app_service_log_streaming" "existing_app_service_log_streaming" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
}
