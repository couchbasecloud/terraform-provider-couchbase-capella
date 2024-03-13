output "existing_auditlogsettings" {
  value = data.couchbase-capella_audit_log_settings.existing_auditlogsettings
}

data "couchbase-capella_audit_log_settings" "existing_auditlogsettings" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}