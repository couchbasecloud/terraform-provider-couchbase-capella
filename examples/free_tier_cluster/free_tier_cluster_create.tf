
output "new_free_tier_cluster" {
  value = couchbase-capella_cluster_free_tier.new_cluster_free_tier
}

output "cluster_free_tierid" {
  value = couchbase-capella_cluster_free_tier.new_cluster_free_tier.id
}

resource "couchbase-capella_cluster_free_tier" "new_cluster_free_tier" {
  organization_id = var.organization_id
  project_id      = var.project_id
  name            = "New free tier cluster modifed"
  description     = "new test cluster for multiple services modified"
  cloud_provider = {
    type   = var.cloud_provider.name
    region = var.cloud_provider.region
    cidr   = var.cluster.cidr
  }
}
