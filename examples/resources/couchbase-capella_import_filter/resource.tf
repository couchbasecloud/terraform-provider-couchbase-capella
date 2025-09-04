resource "couchbase-capella_import_filter" "example" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  app_endpoint_name = var.app_endpoint_name
  scope            = var.scope
  collection       = var.collection
  # JavaScript function body
  import_filter = <<-EOT
  function (doc, meta) {
    return doc.type == "order";
  }
  EOT
}

