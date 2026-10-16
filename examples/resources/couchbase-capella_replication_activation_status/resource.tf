# Pauses the replication. Set state = "running" to resume it.
resource "couchbase-capella_replication_activation_status" "new_replication_activation_status" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  replication_id  = "<replication_id>"
  state           = "paused"
}

# Pause and resume change one direction only. For a twoWay replication, the reverse direction
# needs its own resource, with the target cluster ID and the reverse_replication_id of the replication.
resource "couchbase-capella_replication_activation_status" "reverse_replication_activation_status" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "aaaaaaaa-bbbb-1414-eeee-000000000000"
  replication_id  = "<reverse_replication_id>"
  state           = "paused"
}
