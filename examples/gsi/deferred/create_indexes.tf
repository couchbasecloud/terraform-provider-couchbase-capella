locals {
  index_template = templatefile("${path.module}/indexes.json", {
    organization_id = var.organization_id
    project_id      = var.project_id
    cluster_id      = var.cluster_id
  })

  decoded_template = jsondecode(local.index_template)
  index_names      = [for idx, details in local.decoded_template.resource["couchbase-capella_indexes"] : details.index_name]
}

resource "couchbase-capella_query_indexes" "new_indexes" {
  for_each        = jsondecode(local.index_template).resource["couchbase-capella_indexes"]
  organization_id = each.value.organization_id
  project_id      = each.value.project_id
  cluster_id      = each.value.cluster_id
  bucket_name     = each.value.bucket_name
  scope_name      = each.value.scope_name
  collection_name = each.value.collection_name
  index_name      = each.value.index_name
  index_keys      = each.value.index_keys
  with = {
    defer_build = each.value.with.defer_build
  }
}

resource "couchbase-capella_query_indexes" "build_idx" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id

  bucket_name     = var.bucket_name
  scope_name      = var.scope_name
  collection_name = var.collection_name

  build_indexes = local.index_names

  depends_on = [couchbase-capella_query_indexes.new_indexes]
}

data "couchbase-capella_query_index_monitor" "mon_indexes" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_name     = var.bucket_name
  scope_name      = var.scope_name
  collection_name = var.collection_name
  indexes         = local.index_names

  depends_on = [couchbase-capella_query_indexes.build_idx]
}

output "new_indexes" {
  value = couchbase-capella_query_indexes.new_indexes
}