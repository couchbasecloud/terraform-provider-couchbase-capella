resource "couchbase-capella_app_service_log_streaming" "app_service_log_streaming" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  output_type = "generic_http"
  credentials = {
    generic_http = {
      user = "<user>"
      password = "<password>"
      url = "<log_collector_url>"
    }
  }
}
