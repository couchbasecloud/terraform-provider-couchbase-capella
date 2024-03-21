output "existing_auditlogexport" {
  value = data.couchbase-capella_audit_log_export.existing_auditlogexport
}

data "couchbase-capella_audit_log_export" "existing_auditlogexport" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
