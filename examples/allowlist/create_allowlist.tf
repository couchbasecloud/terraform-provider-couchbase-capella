output "new_allowlist" {
  value = capella_allowlist.new_allowlist
}

resource "capella_allowlist" "new_allowlist" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  cidr            = var.cidr
  comment         = var.comment
  expires_at      = var.expires_at
}
