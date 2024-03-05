
output "new_auditlogsettings" {
  value = couchbase-capella_auditlogsettings.new_auditlogsettings
}

resource "couchbase-capella_auditlogsettings" "new_auditlogsettings" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  auditenabled    = var.auditlogsettings.auditenabled
  enabledeventids = var.auditlogsettings.enabledeventids
  disabledusers   = var.auditlogsettings.disabledusers
}