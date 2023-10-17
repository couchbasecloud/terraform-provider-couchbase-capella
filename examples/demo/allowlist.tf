resource "capella_allowlist" "new_allowlist" {
  organization_id = data.capella_organization.existing_organization.id
  project_id      = capella_project.new_project.id
  cluster_id      = capella_cluster.new_cluster.id
  cidr            = var.cidr
  comment         = var.comment
  expires_at      = "2023-11-30T23:59:59.465Z"
}