# Reads a replication by its ID. cluster_id is the cluster where the replication starts.
data "couchbase-capella_replication" "existing_replication" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  id              = "<replication_id>"
}

output "existing_replication" {
  value = data.couchbase-capella_replication.existing_replication
}

# Pauses the replication that the data source reads.
resource "couchbase-capella_replication_activation_status" "pause_existing_replication" {
  organization_id = data.couchbase-capella_replication.existing_replication.organization_id
  project_id      = data.couchbase-capella_replication.existing_replication.project_id
  cluster_id      = data.couchbase-capella_replication.existing_replication.cluster_id
  replication_id  = data.couchbase-capella_replication.existing_replication.id
  state           = "paused"
}
