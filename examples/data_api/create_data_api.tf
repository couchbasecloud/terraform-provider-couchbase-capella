output "new_data_api" {
  value = couchbase-capella_data_api.new_data_api
}

resource "couchbase-capella_data_api" "new_data_api" {
  organization_id        = var.organization_id
  project_id             = var.project_id
  cluster_id             = var.cluster_id
  enable_data_api        = var.data_api.enable_data_api
  enable_network_peering = var.data_api.enable_network_peering
}
