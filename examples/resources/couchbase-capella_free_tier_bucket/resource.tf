resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
  organization_id         = "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  project_id              = "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa"
  cluster_id              = "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  name                    = "test_bucket"
  memory_allocation_in_mb = 250
}