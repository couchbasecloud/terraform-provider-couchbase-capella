data "couchbase-capella_app_services_cidr" "existing_allowed_cidr" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
}