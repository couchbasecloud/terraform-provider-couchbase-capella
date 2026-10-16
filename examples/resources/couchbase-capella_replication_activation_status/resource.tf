# Pauses the replication. Set state = "running" to resume it.
resource "couchbase-capella_replication_activation_status" "new_replication_activation_status" {
  organization_id = couchbase-capella_replication.new_replication.organization_id
  project_id      = couchbase-capella_replication.new_replication.project_id
  cluster_id      = couchbase-capella_replication.new_replication.cluster_id
  replication_id  = couchbase-capella_replication.new_replication.id
  state           = "paused"
}

# Pause and resume change one direction only. For a twoWay replication, the reverse direction
# needs its own resource, with the target cluster ID and the reverse replication ID.
resource "couchbase-capella_replication_activation_status" "reverse_replication_activation_status" {
  organization_id = couchbase-capella_replication.new_replication.organization_id
  project_id      = couchbase-capella_replication.new_replication.project_id
  cluster_id      = couchbase-capella_replication.new_replication.target.cluster
  replication_id  = couchbase-capella_replication.new_replication.reverse_replication_id
  state           = "paused"
}
