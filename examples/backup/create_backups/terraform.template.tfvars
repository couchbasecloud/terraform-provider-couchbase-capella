auth_token = "<v4-api-key-secret>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"
bucket_id       = "<bucket_id>"

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