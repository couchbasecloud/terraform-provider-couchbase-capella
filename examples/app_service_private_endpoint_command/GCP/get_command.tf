output "gcp_command" {
  value = data.couchbase-capella_app_service_gcp_private_endpoint_command.gcp_command
}

data "couchbase-capella_app_service_gcp_private_endpoint_command" "gcp_command" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  vpc_network_id  = var.vpc_network_id
  subnet_ids      = var.subnet_ids
}
