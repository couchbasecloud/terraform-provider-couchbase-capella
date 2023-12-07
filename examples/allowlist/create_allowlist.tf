output "new_allowlist" {
  value = couchbase-capella_allowlist.new_allowlist
}

output "allowlist_id" {
  value = couchbase-capella_allowlist.new_allowlist.id
}

resource "couchbase-capella_allowlist" "new_allowlist" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  cidr            = var.allowlist.cidr
  comment         = var.allowlist.comment
  expires_at      = var.allowlist.expires_at
}
