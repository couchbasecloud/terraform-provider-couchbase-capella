output "azure_network_peer_command" {
  value = data.couchbase-capella_azure_network_peer_command.azure_network_peer_command
}

data "couchbase-capella_azure_network_peer_command" "azure_network_peer_command" {
  organization_id                = var.organization_id
  project_id                     = var.project_id
  cluster_id                     = var.cluster_id
  tenant_id                      = var.tenant_id
  vnet_id                        = var.vnet_id
  subscription_id                = var.subscription_id
  resource_group                 = var.resource_group
  vnet_peering_service_principal = var.vnet_peering_service_principal
}