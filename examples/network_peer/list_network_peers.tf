output "network_peers_list" {
  value = data.couchbase-capella_network_peers.existing_network_peers
}

data "couchbase-capella_network_peers" "existing_network_peers" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
