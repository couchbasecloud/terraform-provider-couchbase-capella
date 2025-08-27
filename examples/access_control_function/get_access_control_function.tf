output "existing_acess_control_function" {
    value = data.couchbase-capella_access_control_function.existing_acess_control_function
}

data "couchbase-capella_access_control_function." "existing_acess_control_function" {
  organization_id         = var.access_control_function.organization_id
  project_id              = var.access_control_function.project_id
  cluster_id              = var.access_control_function.cluster_id
  app_service_id          = var.access_control_function.app_service_id
  app_endpoint_name       = var.access_control_function.app_endpoint_name
  scope                   = var.access_control_function.scope
  collection              = var.access_control_function.collection
}