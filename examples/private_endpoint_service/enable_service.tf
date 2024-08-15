resource "couchbase-capella_private_endpoint_service" "new_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  enabled         = var.enabled
}