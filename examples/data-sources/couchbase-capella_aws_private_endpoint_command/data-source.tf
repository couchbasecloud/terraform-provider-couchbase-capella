data "couchbase-capella_aws_private_endpoint_command" "aws_command" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  vpc_id          = "vpc-1234"
  subnet_ids      = "["subnet-1234",]"
}