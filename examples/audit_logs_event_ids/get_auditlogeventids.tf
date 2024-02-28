output "existing_auditlogeventids" {
  value = data.couchbase-capella_audit_log_event_ids.existing_auditlogeventids
}

data "couchbase-capella_audit_log_event_ids" "existing_auditlogeventids" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}