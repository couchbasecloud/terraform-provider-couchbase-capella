output "new_free_tier_bucket" {
  value = couchbase-capella_free_tier_bucket.new_free_tier_bucket
}

output "free_tier_bucket_id" {
  value = couchbase-capella_free_tier_bucket.new_free_tier_bucket.id
}

resource "couchbase-capella_free_tier_bucket" "new_free_tier_bucket" {
  organization_id         = var.organization_id
  project_id              = var.project_id
  cluster_id              = var.cluster_id
  name                    = "test_bucket"
  memory_allocation_in_mb = 250
}

