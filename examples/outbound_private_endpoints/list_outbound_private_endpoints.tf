data "couchbase-capella_outbound_private_endpoints" "existing_outbound_private_endpoints" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

output "existing_outbound_private_endpoints" {
  value = data.couchbase-capella_outbound_private_endpoints.existing_outbound_private_endpoints
}
