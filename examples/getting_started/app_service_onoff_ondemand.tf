# Stores the app service onoff details in an output variable.
# Can be viewed using `terraform output app_service_onoff_ondemand` command
output "app_service_onoff_ondemand" {
  value = couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand
}

resource "couchbase-capella_app_service_onoff_ondemand" "new_app_service_onoff_ondemand" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  state           = var.app_service_state
}

