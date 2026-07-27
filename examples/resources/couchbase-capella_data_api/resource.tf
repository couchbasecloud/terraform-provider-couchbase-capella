resource "couchbase-capella_data_api" "new_data_api" {
  organization_id        = "<organization_id>"
  project_id             = "<project_id>"
  cluster_id             = "<cluster_id>"
  enable_data_api        = true
  enable_network_peering = false
}
