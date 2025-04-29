resource "couchbase-capella_free_tier_app_service" "new_free_tier_app_service" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  name            = "free-tier-app-service"
  description     = "Free Tier App Service created by terraform"
}