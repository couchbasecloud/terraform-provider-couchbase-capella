auth_token = "<v4-api-key-secret>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

backup = {}

restore = {
  target_cluster_id = "<cluster_id>"
  source_cluster_id = "<cluster_id>"
  services = [
    "data",
    "query"
  ]
  force_updates           = true
  auto_remove_collections = true
  restore_times           = 1
}

bucket = {
  name                       = "new_terraform_bucket"
  type                       = "couchbase"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 100
  bucket_conflict_resolution = "seqno"
  durability_level           = "none"
  replicas                   = 1
  flush                      = false
  time_to_live_in_seconds    = 0
  eviction_policy            = "fullEviction"
}
