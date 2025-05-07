resource "couchbase-capella_free_tier_cluster" "new_free_tier_cluster" {
  organization_id = "organization_id"
  project_id      = "project_id"
  name            = "New free tier cluster"
  description     = "New free tier test cluster for multiple services"
  cloud_provider = {
    type   = "aws"
    region = "capella_free_tier_cluster"
    cidr   = "10.1.0.0/16"
  }
}
