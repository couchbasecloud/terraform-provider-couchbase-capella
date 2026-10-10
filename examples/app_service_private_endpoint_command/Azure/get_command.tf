output "azure_command" {
  value = data.couchbase-capella_app_service_azure_private_endpoint_command.azure_command
}

data "couchbase-capella_app_service_azure_private_endpoint_command" "azure_command" {
  organization_id     = var.organization_id
  project_id          = var.project_id
  cluster_id          = var.cluster_id
  app_service_id      = var.app_service_id
  resource_group_name = var.resource_group_name
  virtual_network     = var.virtual_network
}
