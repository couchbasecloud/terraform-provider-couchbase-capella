resource "couchbase-capella_app_service_audit_log_streaming" "new_app_service_audit_log_streaming" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  streaming_enabled      = var.streaming_enabled
  output_type            = "datadog"
  disabled_app_endpoints = var.disabled_app_endpoints

  credentials = {
    datadog = {
      url     = var.datadog_url
      api_key = var.datadog_api_key
    }
  }
}

output "new_app_service_audit_log_streaming" {
  value     = couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming
  sensitive = true
}
