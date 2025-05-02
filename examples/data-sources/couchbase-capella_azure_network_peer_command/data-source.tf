data "couchbase-capella_azure_network_peer_command" "azure_network_peer_command" {
  organization_id                = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id                     = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id                     = "ffffffff-aaaa-1414-eeee-000000000000"
  tenant_id                      = "ffffffff-aaaa-1414-eeee-000000000000"
  vnet_id                        = "test_vnet"
  subscription_id                = "ffffffff-aaaa-1414-eeee-000000000000"
  resource_group                 = "test_rg"
  vnet_peering_service_principal = "ffffffff-aaaa-1414-eeee-000000000000"
}
