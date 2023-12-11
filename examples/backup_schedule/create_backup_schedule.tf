output "new_backup_schedule" {
  value = couchbase-capella_backup_schedule.new_backup_schedule
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

resource "couchbase-capella_backup_schedule" "new_backup_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  bucket_id       = couchbase-capella_bucket.new_bucket.id
  type            = var.backup_schedule.type
  weekly_schedule = {
    day_of_week              = var.backup_schedule.weekly_schedule.day_of_week
    start_at                 = var.backup_schedule.weekly_schedule.start_at
    incremental_every        = var.backup_schedule.weekly_schedule.incremental_every
    retention_time           = var.backup_schedule.weekly_schedule.retention_time
    cost_optimized_retention = var.backup_schedule.weekly_schedule.cost_optimized_retention
  }
}