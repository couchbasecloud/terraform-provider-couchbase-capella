# Stores the App Endpoint Access Control Function in an output variable.
# Can be viewed using `terraform output app_endpoint_access_function` command
output "app_endpoint_access_function" {
  value = couchbase-capella_app_endpoint_access_control_function.acf
}

resource "couchbase-capella_app_endpoint_access_control_function" "acf" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.new_app_endpoint.name

  scope      = var.app_endpoint_function.scope
  collection = var.app_endpoint_function.collection
  access_control_function = var.app_endpoint_function.access_control_function
}


