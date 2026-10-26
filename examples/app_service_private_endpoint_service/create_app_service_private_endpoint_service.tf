resource "couchbase-capella_app_service_private_endpoint_service" "new_app_service_private_endpoint_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}

output "new_app_service_private_endpoint_service" {
  value = couchbase-capella_app_service_private_endpoint_service.new_app_service_private_endpoint_service
}
