output "new_replication_activation_status" {
  value = couchbase-capella_replication_activation_status.new_replication_activation_status
}

resource "couchbase-capella_replication_activation_status" "new_replication_activation_status" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  replication_id  = var.replication_id

  state = var.state
}
