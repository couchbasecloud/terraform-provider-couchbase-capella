output "new_bucket" {
  value = capella_bucket.new_bucket
}

resource "capella_bucket" "new_bucket" {
  name                  = var.bucket.name
  organization_id       = var.organization_id
  project_id            = var.project_id
  cluster_id            = var.cluster_id
  type                  = var.bucket.type
  storage_backend       = var.bucket.storage_backend
  memory_allocationinmb = var.bucket.memory_allocationinmb
  conflict_resolution   = var.bucket.conflict_resolution
  durability_level      = var.bucket.durability_level
  replicas              = var.bucket.replicas
  flush                 = var.bucket.flush
  ttl                   = var.bucket.ttl
  eviction_policy       = var.bucket.eviction_policy
}