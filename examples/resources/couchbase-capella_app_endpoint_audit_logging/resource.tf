resource "couchbase-capella_app_endpoint_audit_logging" "new_app_endpoint_audit_logging" {
  organization_id   = "<organization_id>"
  project_id        = "<project_id>"
  cluster_id        = "<cluster_id>"
  app_service_id    = "<app_service_id>"
  app_endpoint_name = "<app_endpoint_name>"

  audit_enabled     = true
  enabled_event_ids = [53280, 53281]

  disabled_users = [
    {
      domain = "sgw"
      name   = "<user_name>"
    }
  ]

  disabled_roles = [
    {
      domain = "sgw"
      name   = "<role_name>"
    }
  ]
}
