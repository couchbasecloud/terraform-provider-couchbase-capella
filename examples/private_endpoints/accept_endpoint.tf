resource "couchbase-capella_private_endpoints" "accept_endpoint" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  endpoint_id     = var.endpoint_id
}