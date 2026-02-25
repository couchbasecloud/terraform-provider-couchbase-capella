output "new_app_endpoint_log_streaming_config" {
  value = couchbase-capella_app_endpoint_log_streaming_config.new_app_endpoint_log_streaming_config
}

resource "couchbase-capella_app_endpoint_log_streaming_config" "new_app_endpoint_log_streaming_config" {
  organization_id  = var.organization_id
  project_id = var.project_id
  cluster_id = var.cluster_id
  app_service_id = var.app_service_id
  app_endpoint_name = var.app_endpoint_name

  log_level = var.app_endpoint_log_streaming_config.log_level
  log_keys = var.app_endpoint_log_streaming_config.log_keys
}