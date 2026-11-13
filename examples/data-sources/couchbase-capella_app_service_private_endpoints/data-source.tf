data "couchbase-capella_app_service_private_endpoints" "app_service_private_endpoints" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  # Optional. Omit to return private endpoints in every status.
  status = ["linked", "pendingAcceptance"]
}
