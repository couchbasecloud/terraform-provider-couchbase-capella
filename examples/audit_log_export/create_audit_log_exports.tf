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
