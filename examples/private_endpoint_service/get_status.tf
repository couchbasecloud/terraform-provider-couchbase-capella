output "service_status" {
  value = data.couchbase-capella_private_endpoint_service.service_status
}

data "couchbase-capella_private_endpoint_service" "service_status" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}