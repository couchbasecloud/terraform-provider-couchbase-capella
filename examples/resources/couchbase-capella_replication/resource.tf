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
