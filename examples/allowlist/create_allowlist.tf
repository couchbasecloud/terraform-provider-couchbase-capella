output "new_allowlist" {
  value = capella_allowlist.new_allowlist
}

output "allowlist_id" {
  value = capella_allowlist.new_allowlist.id
}

resource "capella_allowlist" "new_allowlist" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  cidr            = var.allowlist.cidr
  comment         = var.allowlist.comment
  expires_at      = var.allowlist.expires_at
}
