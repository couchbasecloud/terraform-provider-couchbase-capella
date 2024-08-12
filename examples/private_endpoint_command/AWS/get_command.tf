output "aws_command" {
  value = data.couchbase-capella_aws_private_endpoint_command.aws_command
}

data "couchbase-capella_aws_private_endpoint_command" "aws_command" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  vpc_id          = var.vpc_id
  subnet_ids      = var.subnet_ids
}