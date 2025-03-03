
output "new_free_tier_cluster" {
  value = couchbase-capella_free_tier_cluster.new_free_tier_cluster
}

output "free_tier_cluster_id" {
  value = couchbase-capella_free_tier_cluster.new_free_tier_cluster.id
}

resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
  organization_id = var.organization_id
  project_id      = var.project_id
  name            = "New free tier cluster"
  description     = "New free tier test cluster for multiple services"
  cloud_provider = {
    type   = var.cloud_provider.name
    region = var.cloud_provider.region
    cidr   = var.cloud_provider.cidr
  }
}
