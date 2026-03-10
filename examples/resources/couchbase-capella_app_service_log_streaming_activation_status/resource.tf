resource "couchbase-capella_app_service_log_streaming_activation_status" "app_service_log_streaming_activation_status" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  state           = "<target-state>"
}
