resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
  organization_id         = var.organization_id
  project_id              = couchbase-capella_project.new_project.id
  cluster_id              = couchbase-capella_free_tier_cluster.new_free_tier_cluster.id
  name                    = var.free_tier_bucket.name
  memory_allocation_in_mb = var.free_tier_bucket.memory_allocation_in_mb
}
