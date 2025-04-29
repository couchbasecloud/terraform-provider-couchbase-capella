data "couchbase-capella_free_tier_buckets" "existing_buckets" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
}