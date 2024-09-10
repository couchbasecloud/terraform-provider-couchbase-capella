resource "couchbase-capella_private_endpoint_service" "new_service" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  enabled         = var.enabled
}
