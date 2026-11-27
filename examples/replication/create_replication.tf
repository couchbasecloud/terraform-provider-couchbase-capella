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

output "new_mobile_replication" {
  value = couchbase-capella_replication.new_mobile_replication
}

output "mobile_replication_id" {
  value = couchbase-capella_replication.new_mobile_replication.id
}

# A replication with mobile sync mode and conflict logging. The create turns on cross cluster versioning
# for the source and target buckets, and cross cluster versioning cannot be turned off.
# Conflict logging needs Couchbase Server 8.0 or later.
resource "couchbase-capella_replication" "new_mobile_replication" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  source_bucket   = var.mobile_replication.source_bucket_id

  target = {
    cluster = var.target.cluster
    bucket  = var.mobile_replication.target_bucket_id
  }

  mobile_sync_mode = "active"

  conflict_logging = {
    disabled   = var.conflict_logging_disabled
    bucket     = var.mobile_replication.log_bucket_id
    collection = "logs.default"

    logging_rules = {
      "inventory"         = { bucket = var.mobile_replication.log_bucket_id, collection = "logs.inventory" } # log to this target
      "inventory.airline" = {}                                                                               # log to the default
      "inventory.airport" = null                                                                             # do not log
    }
  }
}
