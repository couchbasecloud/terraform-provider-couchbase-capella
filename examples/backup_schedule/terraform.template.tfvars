auth_token = "<v4-api-key-secret>"
host       = "https://cloudapi.cloud.couchbase.com"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"
bucket_id       = "<bucket_id>"

backup_schedule = {
  type = "weekly"
  weekly_schedule = {
    day_of_week = "sunday"
    start_at = 10
    incremental_every = 4
    retention_time = "90days"
    cost_optimized_retention = false
  }
}