data "couchbase-capella_azure_private_endpoint_command" "azure_command" {
  organization_id     = "<organization_id>"
  project_id          = "<project_id>"
  cluster_id          = "<cluster_id>"
  resource_group_name = "test-rg"
  virtual_network     = "vnet-1/subnet-1"
}
