output "new_import_filter" {
  value = couchbase-capella_app_endpoint_import_filter.new_import_filter
}

resource "couchbase-capella_app_endpoint_import_filter" "new_import_filter" {
  organization_id      = var.import_filter.organization_id
  project_id           = var.import_filter.project_id
  cluster_id           = var.import_filter.cluster_id
  app_service_id       = var.import_filter.app_service_id
  app_endpoint_name    = var.import_filter.app_endpoint_name
  scope                = var.import_filter.scope
  collection           = var.import_filter.collection
  import_filter        = var.import_filter.import_filter
}


