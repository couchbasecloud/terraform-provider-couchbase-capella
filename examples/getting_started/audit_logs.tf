# Enable audit logs for data operations (event ids 20488, 20490, 20491).
output "new_auditlogsettings" {
  value = couchbase-capella_audit_log_settings.new_auditlogsettings
}

resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  audit_enabled     = var.audit_log_settings.audit_enabled
  enabled_event_ids = var.audit_log_settings.enabled_event_ids
  disabled_users    = var.audit_log_settings.disabled_users
}

# Obtain audit logs for a 1 hour time frame.
# This example assumes there is data traffic in this time frame.
output "new_auditlogexport" {
  value = couchbase-capella_audit_log_export.new_auditlogexport
}

resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  start           = var.audit_log_export.start
  end             = var.audit_log_export.end
}
