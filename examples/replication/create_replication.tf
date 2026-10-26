output "new_replication" {
  value = couchbase-capella_replication.new_replication
}

output "replication_id" {
  value = couchbase-capella_replication.new_replication.id
}

resource "couchbase-capella_replication" "new_replication" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  source_bucket   = var.source_bucket_id

  target = var.target

  direction           = var.replication.direction
  priority            = var.replication.priority
  network_usage_limit = var.replication.network_usage_limit

  mappings = [
    {
      source_scope = "inventory"
      target_scope = "inventory"
      collections = [
        { source_collection = "airline", target_collection = "airline" },
        { source_collection = "airport", target_collection = "airport" },
      ]
    },
  ]

  filter = {
    document_exclude_options = {
      ttl = true
    }
    expressions = {
      reg_ex = "REGEXP_CONTAINS(META().id, '^airline_')"
    }
  }
}
