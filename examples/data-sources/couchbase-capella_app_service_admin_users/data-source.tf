data "couchbase-capella_app_service_admin_users" "existing_admin_users" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  # Optional: return only the admin users with these names.
  filter {
    name   = "name"
    values = ["<admin_user_name>"]
  }
}
