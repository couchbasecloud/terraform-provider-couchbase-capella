output "new_bucket" {
  value = couchbase-capella_bucket.new_bucket
}

output "bucket_id" {
  value = couchbase-capella_bucket.new_bucket.id
}

resource "couchbase-capella_bucket" "new_bucket" {
  name                       = var.bucket.name
  organization_id            = var.organization_id
  project_id                 = var.project_id
  cluster_id                 = var.cluster_id
  type                       = var.bucket.type
  storage_backend            = var.bucket.storage_backend
  memory_allocation_in_mb    = var.bucket.memory_allocation_in_mb
  bucket_conflict_resolution = var.bucket.bucket_conflict_resolution
  durability_level           = var.bucket.durability_level
  replicas                   = var.bucket.replicas
  flush                      = var.bucket.flush
  time_to_live_in_seconds    = var.bucket.time_to_live_in_seconds
  eviction_policy            = var.bucket.eviction_policy
}
