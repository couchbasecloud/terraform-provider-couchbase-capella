resource "couchbase-capella_outbound_private_endpoint" "new_outbound_private_endpoint" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  name            = var.outbound_private_endpoint.name
  service_name    = var.outbound_private_endpoint.service_name
  hostname        = var.outbound_private_endpoint.hostname
}

output "new_outbound_private_endpoint" {
  value = couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint
}

output "outbound_private_endpoint_id" {
  value = couchbase-capella_outbound_private_endpoint.new_outbound_private_endpoint.id
}
