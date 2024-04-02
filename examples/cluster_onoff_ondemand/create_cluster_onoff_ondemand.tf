output "new_cluster_onoff_ondemand" {
  value = couchbase-capella_cluster_onoff_ondemand.new_cluster_onoff_ondemand
}

resource "couchbase-capella_cluster_onoff_ondemand" "new_cluster_onoff_ondemand" {
  organization_id            = var.organization_id
  project_id                 = var.project_id
  cluster_id                 = var.cluster_id
  state                      = var.state
  turn_on_linked_app_service = var.cluster_onoff_ondemand.turn_on_linked_app_service
}

