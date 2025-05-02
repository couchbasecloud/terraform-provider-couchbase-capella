resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  audit_enabled     = true
  enabled_event_ids = [20488, 20490,20491,]
  disabled_users    = []
}