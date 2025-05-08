resource "couchbase-capella_bucket" "new_bucket" {
  name                       = "terraform bucket"
  organization_id            = "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  project_id                 = "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  cluster_id                 = "aaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  type                       = "couchbase"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 100
  bucket_conflict_resolution = "seqno"
  durability_level           = "none"
  replicas                   = 1
  flush                      = false
  time_to_live_in_seconds    = 1
  eviction_policy            = "fullEviction"
}