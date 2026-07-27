output "existing_data_api" {
  value = data.couchbase-capella_data_api.existing_data_api
}

data "couchbase-capella_data_api" "existing_data_api" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
