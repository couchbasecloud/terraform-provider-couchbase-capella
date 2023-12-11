auth_token = "<v4-api-key-secret>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

backup_schedule = {
  type = "weekly"
  weekly_schedule = {
    day_of_week              = "sunday"
    start_at                 = 10
    incremental_every        = 4
    retention_time           = "90days"
    cost_optimized_retention = false
  }
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
