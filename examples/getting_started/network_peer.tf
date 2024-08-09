# Stores the network peer details in an output variable.
# Can be viewed using `terraform output network_peer` command
output "network_peer" {
  value = couchbase-capella_network_peer.new_network_peer
}

resource "couchbase-capella_network_peer" "new_network_peer" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name            = var.network_peer.name
  provider_type   = var.network_peer.provider_type
  provider_config = {
    aws_config = {
      account_id = var.aws_config.account_id
      vpc_id     = var.aws_config.vpc_id
      cidr       = var.aws_config.cidr
      region     = var.aws_config.region
    }
  }
}

