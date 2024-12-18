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
  provider_config = {
    aws_config = {
      account_id = var.aws_config.account_id
      vpc_id     = var.aws_config.vpc_id
      cidr       = var.aws_config.cidr
      region     = var.aws_config.region
    }
  }
}

# Example GCP Config for creating network peer on GCP. Use this instead of aws_config above if you want to create a network peer for GCP.
#      gcp_config = {
#             project_id = var.gcp_config.project_id
#             network_name = var.gcp_config.network_name
#             service_account = var.gcp_config.service_account
#             cidr  = var.gcp_config.cidr
#              }


# Example Azure Config for creating network peer on Azure. Use this instead of aws_config above if you want to create a network peer for Azure.
#      azure_config = {
#             tenant_id       = var.azure_config.tenant_id
#             resource_group  = var.azure_config.resource_group
#             subscription_id = var.azure_config.subscription_id
#             cidr            = var.azure_config.cidr
#             vnet_id         = var.azure_config.vnet_id
#              }