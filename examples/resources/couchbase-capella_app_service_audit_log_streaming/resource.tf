resource "couchbase-capella_app_service_audit_log_streaming" "new_app_service_audit_log_streaming" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  streaming_enabled      = true
  output_type            = "datadog"
  disabled_app_endpoints = ["<app_endpoint_name>"]

  credentials = {
    datadog = {
      url     = "<datadog_url>"
      api_key = "<datadog_api_key>"
    }
  }
}
