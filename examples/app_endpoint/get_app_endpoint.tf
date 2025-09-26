output "some_app_endpoints" {
  value = data.couchbase-capella_app_endpoints.filtered_app_endpoints
}

data "couchbase-capella_app_endpoints" "filtered_app_endpoints" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  app_service_id = var.app_service.id

  # filter app endpoints by attribute name (i.e. app endpoint name)
  # values are the actual names
  filter {
    name   = "name"
    values = ["app-endpoint-1"]
  }
}