resource "couchbase-capella_app_service_audit_log_streaming" "new_app_service_audit_log_streaming" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  streaming_enabled = var.streaming_enabled
  # One of: datadog, sumologic, elastic, loki, splunk, dynatrace, generic_http.
  # Set the matching block in credentials.
  output_type            = "datadog"
  disabled_app_endpoints = var.disabled_app_endpoints

  credentials = {
    datadog = {
      url     = var.datadog_url
      api_key = var.datadog_api_key
    }

    # sumologic = {
    #   url = var.sumologic_url
    # }

    # elastic = {
    #   url      = var.elastic_url
    #   user     = var.elastic_user
    #   password = var.elastic_password
    # }

    # loki = {
    #   url      = var.loki_url
    #   user     = var.loki_user
    #   password = var.loki_password
    # }

    # splunk = {
    #   url          = var.splunk_url
    #   splunk_token = var.splunk_token
    # }

    # dynatrace = {
    #   url       = var.dynatrace_url
    #   api_token = var.dynatrace_api_token
    # }

    # generic_http = {
    #   url      = var.generic_http_url
    #   user     = var.generic_http_user
    #   password = var.generic_http_password
    # }
  }
}

output "new_app_service_audit_log_streaming" {
  value     = couchbase-capella_app_service_audit_log_streaming.new_app_service_audit_log_streaming
  sensitive = true
}
