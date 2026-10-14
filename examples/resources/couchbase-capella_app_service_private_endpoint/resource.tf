resource "couchbase-capella_app_service_private_endpoint" "accept_endpoint" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
  endpoint_id     = "<endpoint_id>"
}
