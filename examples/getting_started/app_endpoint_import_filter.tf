# Stores the App Endpoint Import Filter in an output variable.
# Can be viewed using `terraform output app_endpoint_import_filter` command
output "app_endpoint_import_filter" {
  value = couchbase-capella_app_endpoint_import_filter.import_filter
}

resource "couchbase-capella_app_endpoint_import_filter" "import_filter" {
  organization_id   = var.organization_id
  project_id        = couchbase-capella_project.new_project.id
  cluster_id        = couchbase-capella_cluster.new_cluster.id
  app_service_id    = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name = couchbase-capella_app_endpoint.new_app_endpoint.name

  scope         = var.app_endpoint_import_filter.scope
  collection    = var.app_endpoint_import_filter.collection
  import_filter = var.app_endpoint_import_filter.import_filter
}


