resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
  organization_id         = "<organization_id>"
  project_id              = "<project_id>"
  cluster_id              = "<cluster_id>"
  name                    = "test_bucket"
  memory_allocation_in_mb = 250
}