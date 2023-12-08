resource "couchbase-capella_allowlist" "new_allowlist" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  cidr            = var.cidr
  comment         = var.comment
  expires_at      = var.expires_at
}
