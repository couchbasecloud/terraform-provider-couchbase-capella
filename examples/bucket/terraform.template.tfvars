auth_token      = "v4-api-key-secret"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"
host            = "https://cloudapi.cloud.couchbase.com"

bucket = {
  name                  = "test_bucket"
  type                  = "couchbase"
  storage_backend       = "couchstore"
  memory_allocationinmb = 105
  conflict_resolution   = "seqno"
  durability_level      = "majorityAndPersistActive"
  replicas              = 2
  flush                 = true
  ttl                   = 100
}



