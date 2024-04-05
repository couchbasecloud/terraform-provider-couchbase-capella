data "couchbase-capella_audit_log_event_ids" "event_list" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

# List of query event ids
locals {
  n1ql_event_ids = [for event in data.couchbase-capella_audit_log_event_ids.event_list.data : event.id if event.module == "n1ql"]
}

# Local variable n1ql_event_ids is used to provide
resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
  organization_id   = var.organization_id
  project_id        = var.project_id
  cluster_id        = var.cluster_id
  audit_enabled     = var.audit_log_settings.audit_enabled
  enabled_event_ids = local.n1ql_event_ids
  disabled_users    = var.audit_log_settings.disabled_users
}

output "new_auditlogsettings" {
  value = couchbase-capella_audit_log_settings.new_auditlogsettings
}

# Obtain audit logs for a 1 hour time frame.
# This example assumes there is query traffic in this time frame.
output "new_auditlogexport" {
  value = couchbase-capella_audit_log_export.new_auditlogexport
}

# Create this resource block as needed.  It is included for completeness.
resource "couchbase-capella_audit_log_export" "new_auditlogexport" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  start           = var.audit_log_export.start
  end             = var.audit_log_export.end
}
