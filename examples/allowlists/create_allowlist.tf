output "new_allowlist" {
  value = capella_allowlist.new_allowlist
}

resource "capella_allowlist" "new_allowlist" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  cidr            = "10.0.0.0/16"
  comment         = "Allow access from another VPC"
  expires_at      = "2023-11-14T21:49:58.465Z"
}
