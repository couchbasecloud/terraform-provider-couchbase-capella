resource "couchbase-capella_import_filter" "example" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  # keyspace format: <app_endpoint_name>[.<scope_name>[.<collection_name>]]
  keyspace = var.keyspace

  # JavaScript function body
  import_filter = <<-EOT
  function (doc, meta) {
    return doc.type == "order";
  }
  EOT
}

