output "new_app_service_onoff_ondemand" {
  value = couchbase-capella_app_service_onoff_ondemand.new_app_service_onoff_ondemand
}

resource "couchbase-capella_app_service_onoff_ondemand" "new_app_service_onoff_ondemand" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
  state           = var.state
}

