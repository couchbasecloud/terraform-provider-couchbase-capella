data "couchbase-capella_outbound_private_endpoint_service_verification" "service_check" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  service_name    = var.service_name
}

output "service_check" {
  value = data.couchbase-capella_outbound_private_endpoint_service_verification.service_check
}
