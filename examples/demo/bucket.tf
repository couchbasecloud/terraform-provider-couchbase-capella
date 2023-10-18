# Stores the bucket name in an output variable.
# Can be viewed using `terraform output bucket` command
output "bucket" {
  value = capella_bucket.new_bucket.name
}

resource "capella_bucket" "new_bucket" {
  name                       = var.bucket.name
  organization_id            = data.capella_organization.existing_organization.id
  project_id                 = capella_project.new_project.id
  cluster_id                 = capella_cluster.new_cluster.id
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
