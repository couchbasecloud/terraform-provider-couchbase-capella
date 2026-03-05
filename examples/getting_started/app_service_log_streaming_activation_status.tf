resource "couchbase-capella_app_service_log_streaming_activation_status" "app_service_log_streaming_activation_status" {
  organization_id  = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id = couchbase-capella_app_service.new_app_service.id

  state = var.app_service_log_streaming_activation_status.state

  depends_on = [couchbase-capella_app_service_log_streaming.app_service_log_streaming]
}