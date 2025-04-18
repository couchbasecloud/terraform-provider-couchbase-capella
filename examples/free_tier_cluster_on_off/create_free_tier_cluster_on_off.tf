output "new_free_tier_cluster_on_off" {
  value = couchbase-capella_free_tier_cluster_on_off.new_free_tier_cluster_on_off
}

resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  state           = var.state

}

