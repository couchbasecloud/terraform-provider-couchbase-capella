resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  name            = var.free_tier_cluster.name
  description     = var.free_tier_cluster.description
  cloud_provider = {
    type   = var.free_tier_cloud_provider.name
    region = var.free_tier_cloud_provider.region
    cidr   = var.free_tier_cluster.cidr
  }
}
