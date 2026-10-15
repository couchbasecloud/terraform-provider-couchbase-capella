resource "couchbase-capella_replication" "new_replication" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  source_bucket   = "dHJhdmVsLXNhbXBsZQ=="

  target = {
    cluster = "aaaaaaaa-bbbb-1414-eeee-000000000000"
    bucket  = "dHJhdmVsLXNhbXBsZQ=="
  }

  direction           = "oneWay"
  priority            = "high"
  network_usage_limit = 0

  mappings = [
    {
      source_scope = "inventory"
      target_scope = "inventory"
      collections = [
        { source_collection = "airline", target_collection = "airline" },
      ]
    },
  ]

  filter = {
    document_exclude_options = {
      deletion = false
      ttl      = true
    }
    expressions = {
      reg_ex = "REGEXP_CONTAINS(META().id, '^airline_')"
    }
  }
}

# Mobile sync mode and conflict logging. Both need cross cluster versioning on the source and
# target buckets. A create with mobile_sync_mode = "active" turns cross cluster versioning on.
# Conflict logging needs Couchbase Server 8.0 or later.
resource "couchbase-capella_replication" "mobile_replication" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  source_bucket   = "dHJhdmVsLXNhbXBsZQ=="

  target = {
    cluster = "aaaaaaaa-bbbb-1414-eeee-000000000000"
    bucket  = "dHJhdmVsLXNhbXBsZQ=="
  }

  mobile_sync_mode = "active"

  conflict_logging = {
    disabled   = false
    bucket     = "Y29uZmxpY3Rsb2dz"
    collection = "logs.default"

    logging_rules = {
      "inventory"       = { bucket = "Y29uZmxpY3Rsb2dz", collection = "logs.inventory" } # log to this target
      "sales"           = {}                                                             # log to the default
      "agents.partners" = null                                                           # do not log
    }
  }
}
