resource "couchbase-capella_app_endpoint_access_control_function" "acf1" {
  organization_id         = var.organization_id
  project_id              = couchbase-capella_project.new_project.id
  cluster_id              = couchbase-capella_cluster.new_cluster.id
  app_service_id          = couchbase-capella_app_service.new_app_service.id
  app_endpoint_name       = var.app_endpoint
  scope                   = var.scope.scope_name
  collection              = var.collection.collection_name
  access_control_function = "function (doc, oldDoc, meta) {channel('c1');}"
}


