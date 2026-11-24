resource "couchbase-capella_app_service_public_access" "app_service_public_access" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  public_access = "restricted"
}
