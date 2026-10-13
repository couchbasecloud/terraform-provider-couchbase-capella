resource "couchbase-capella_app_service_admin_user" "new_admin_user" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"
  name            = "admin-user"
  password        = "<password>"
  endpoints       = ["<app_endpoint_name>"]
}
