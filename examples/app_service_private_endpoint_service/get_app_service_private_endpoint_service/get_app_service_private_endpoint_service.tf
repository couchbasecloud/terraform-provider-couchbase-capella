data "couchbase-capella_app_service_private_endpoint_service" "existing_app_service_private_endpoint_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  # Omit the aws block to read only the state and target_state.
  aws = {
    vpc_id     = var.vpc_id
    subnet_ids = var.subnet_ids
  }
}

output "existing_app_service_private_endpoint_service" {
  value = data.couchbase-capella_app_service_private_endpoint_service.existing_app_service_private_endpoint_service
}
