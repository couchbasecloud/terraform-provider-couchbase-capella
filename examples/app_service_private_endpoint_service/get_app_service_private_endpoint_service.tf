data "couchbase-capella_app_service_private_endpoint_service" "existing_app_service_private_endpoint_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = couchbase-capella_app_service_private_endpoint_service.new_app_service_private_endpoint_service.app_service_id

  # TODO(human)
}

output "existing_app_service_private_endpoint_service" {
  value = data.couchbase-capella_app_service_private_endpoint_service.existing_app_service_private_endpoint_service
}
