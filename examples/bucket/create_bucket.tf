output "new_bucket" {
    value = capella_bucket.new_bucket
}

resource "capella_bucket" "new_bucket" {
    name = "test_bucket"
    organization_id = var.organization_id
    project_id      = var.project_id
    cluster_id = var.cluster_id
    type = "couchbase"
    storage_backend = "couchstore"
    memory_allocationinmb = 105
    conflict_resolution = "seqno"
    durability_level = "majorityAndPersistActive"
    replicas = 2
    flush = true
    ttl = 100
}