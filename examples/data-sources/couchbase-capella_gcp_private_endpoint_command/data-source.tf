data "couchbase-capella_gcp_private_endpoint_command" "gcp_command" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  vpc_network_id  = "vpcnet-1234"
  subnet_ids      = ["subnet-1234"]
}