output "new_network_peer" {
  value = couchbase-capella_network_peer.new_network_peer
}

output "peer_id" {
  value = couchbase-capella_network_peer.new_network_peer.id
}

resource "couchbase-capella_network_peer" "new_network_peer" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = var.network_peer.name
  provider_type   = var.network_peer.provider_type
  provider_config =
    {
      GCP_config = {
        network_name    = var.GCP_config.network_name
        project_id      = var.GCP_config.project_id
        cidr            = var.GCP_config.cidr
        service_account = var.GCP_config.service_account
      }
    }
}
