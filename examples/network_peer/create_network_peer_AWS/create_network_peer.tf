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
      AWS_config = {
        account_id =  var.AWS_config.account_id
        vpc_id     =  var.AWS_config.vpc_id
        cidr       =  var.AWS_config.cidr
        region     =  var.AWS_config.region
      }
    }
}
