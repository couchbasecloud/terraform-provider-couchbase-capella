resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_free_tier_cluster.new_free_tier_cluster.id
  state           = var.free_tier_cluster_state.state
}
