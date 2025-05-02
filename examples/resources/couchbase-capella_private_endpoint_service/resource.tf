resource "couchbase-capella_private_endpoint_service" "new_service" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  enabled         = true
}