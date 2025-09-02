output "existing_import_filter" {
  value = data.couchbase-capella_import_filter.existing_import_filter
}

# Note: This data source is illustrative. Ensure the provider exposes the corresponding data source.
data "couchbase-capella_import_filter" "existing_import_filter" {
  organization_id = var.import_filter.organization_id
  project_id      = var.import_filter.project_id
  cluster_id      = var.import_filter.cluster_id
  app_service_id  = var.import_filter.app_service_id
  keyspace        = var.import_filter.keyspace
}


