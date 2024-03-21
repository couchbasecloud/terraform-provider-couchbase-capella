
output "new_auditlogsettings" {
  value = couchbase-capella_audit_log_settings.new_auditlogsettings
}

resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  audit_enabled     = var.auditlogsettings.audit_enabled
  enabled_event_ids = var.auditlogsettings.enabled_event_ids
  disabled_users    = var.auditlogsettings.disabled_users
}