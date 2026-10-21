data "couchbase-capella_app_service_admin_user" "existing_admin_user" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
  id              = "<admin_user_id>"
}
