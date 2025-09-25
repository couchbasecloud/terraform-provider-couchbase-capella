# Stores the App Endpoint activation status resource in an output variable.
# Can be viewed using `terraform output app_endpoint_activation_status` command
output "app_endpoint_activation_status" {
  value = couchbase-capella_app_endpoint_activation_status.activation
}

resource "couchbase-capella_app_endpoint_activation_status" "activation" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.new_app_endpoint.name

  state = var.app_endpoint_activation.state
}


