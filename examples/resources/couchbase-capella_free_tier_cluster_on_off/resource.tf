resource "couchbase-capella_free_tier_cluster_on_off" "new_free_tier_cluster_on_off" {
  organization_id = "organization_id"
  project_id      = "project_id"
  cluster_id      = "cluster_id"
  state           = "on"
}