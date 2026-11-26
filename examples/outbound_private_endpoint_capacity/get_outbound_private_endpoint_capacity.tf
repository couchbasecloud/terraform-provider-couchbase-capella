data "couchbase-capella_outbound_private_endpoint_capacity" "existing_capacity" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}

output "existing_capacity" {
  value = data.couchbase-capella_outbound_private_endpoint_capacity.existing_capacity
}
