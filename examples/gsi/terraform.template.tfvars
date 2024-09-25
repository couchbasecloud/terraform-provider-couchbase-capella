auth_token = "<token>"

organization_id = "<organization_id>"
project_id = "<project_id>"
cluster_id = "<cluster_name>"

bucket_name = "test"
scope_name = "test"
collection_name = "test"

index_name = "idx_pe9"
index_keys = ["sourceairport", "destinationairport", "stops", "airline", "id", "ARRAY_COUNT(schedule)"]
partition_by = ["sourceairport", "destinationairport"]

with = {
  defer_build = false
  num_replica = 1
  num_partition = 8
}