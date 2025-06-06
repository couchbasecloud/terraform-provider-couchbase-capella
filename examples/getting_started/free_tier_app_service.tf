resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_free_tier_cluster.new_free_tier_cluster.id
  name            = var.free_tier_app_service.name
  description     = var.free_tier_app_service.description
}
