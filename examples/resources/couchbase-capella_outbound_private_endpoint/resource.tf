resource "couchbase-capella_outbound_private_endpoint" "new_outbound_private_endpoint" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  name            = "cng-prod"
  service_name    = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
  hostname        = "cng.acme.example.com"
}
