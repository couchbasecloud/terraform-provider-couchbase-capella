resource "couchbase-capella_app_endpoint_audit_logging" "new_app_endpoint_audit_logging" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  app_service_id    = var.app_service_id
  app_endpoint_name = var.app_endpoint_name

  audit_enabled     = true
  enabled_event_ids = [53280, 53281]

  disabled_users = [
    {
      domain = "sgw"
      name   = "<user_name>"
    }
  ]
}

output "new_app_endpoint_audit_logging" {
  value = couchbase-capella_app_endpoint_audit_logging.new_app_endpoint_audit_logging
}
