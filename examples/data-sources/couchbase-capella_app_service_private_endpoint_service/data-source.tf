data "couchbase-capella_app_service_private_endpoint_service" "app_service_private_endpoint_service" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  app_service_id  = "<app_service_id>"

  aws = {
    vpc_id     = "vpc-1234567890"
    subnet_ids = ["subnet-1234"]
  }
}
