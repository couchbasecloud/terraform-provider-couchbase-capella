data "couchbase-capella_audit_log_event_ids" "event_list" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
}

# List of query event ids
locals {
  n1ql_event_ids = [for event in data.couchbase-capella_audit_log_event_ids.event_list.data : event.id if event.module == "n1ql"]
}

# Local variable n1ql_event_ids is used to provide
resource "couchbase-capella_audit_log_settings" "new_auditlogsettings" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  audit_enabled     = var.audit_log_settings.audit_enabled
  enabled_event_ids = local.n1ql_event_ids
  disabled_users    = var.audit_log_settings.disabled_users
}

output "new_auditlogsettings" {
  value = couchbase-capella_audit_log_settings.new_auditlogsettings
}

