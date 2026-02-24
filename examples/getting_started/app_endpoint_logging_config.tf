resource "couchbase-capella_app_endpoint_log_streaming_config" "new_app_endpoint_log_streaming_config" {
  organization_id  = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.endpoint1.name

  log_level = var.app_endpoint_log_streaming_config.log_level
  log_keys = var.app_endpoint_log_streaming_config.log_keys

  depends_on = [
    couchbase-capella_project.new_project, 
    couchbase-capella_cluster.new_cluster, 
    couchbase-capella_app_service.new_app_service, 
    couchbase-capella_app_endpoint.endpoint1
    ]
}