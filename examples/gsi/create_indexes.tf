locals {
  index_template = templatefile("${path.module}/indexes.json", {
    organization_id = var.organization_id
    project_id      = var.project_id
    cluster_id      = var.cluster_id
  })
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
}
