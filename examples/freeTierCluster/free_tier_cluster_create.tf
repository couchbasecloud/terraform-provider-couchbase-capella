
output "new_free_tier_cluster" {
  value = couchbase-capella_cluster_free_tierr.new_cluster_free_tier
}

output "cluster_free_tierid" {
  value = couchbase-capella_cluster_free_tier.new_cluster_free_tier.id
}

resource "couchbase-capella__cluster_free_tier" "new_cluster_free_tier" {
  organization_id = var.organization_id
  project_id      = var.project_id
  name            = var.cluster.name
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = var.cloud_provider.name
    region = var.cloud_provider.region
    cidr   = var.cluster.cidr
  }
}
