auth_token       = "<v4-api-key-secret>"
organization_id  = "<organization_id>"
project_id       = "<project_id>"
cluster_id       = "<source_cluster_id>"
source_bucket_id = "<source_bucket_id>"

target = {
  cluster = "<target_cluster_id>"
  bucket  = "<target_bucket_id>"
}

replication = {
  direction           = "oneWay"
  priority            = "high"
  network_usage_limit = 0
}

mobile_replication = {
  source_bucket_id = "<mobile_source_bucket_id>"
  target_bucket_id = "<mobile_target_bucket_id>"
  log_bucket_id    = "<log_bucket_id>"
}
