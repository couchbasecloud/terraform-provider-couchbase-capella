terraform {
  required_providers {
    couchbase-capella = {
      source = "couchbasecloud/couchbase-capella"
    }
  }
}

variable "host" {}
variable "auth_token" { sensitive = true }
variable "organization_id" {}
variable "project_id" {}
variable "cluster_id" {}
variable "app_service_id" {}

# Set both to also get the AWS command. Leave unset to get only the status.
variable "vpc_id" { default = null }
variable "subnet_ids" {
  type    = list(string)
  default = null
}

provider "couchbase-capella" {
  host                 = var.host
  authentication_token = var.auth_token
}

data "couchbase-capella_app_service_private_endpoint_service" "this" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id

  aws = var.vpc_id == null ? null : {
    vpc_id     = var.vpc_id
    subnet_ids = var.subnet_ids
  }
}

output "state" {
  value = data.couchbase-capella_app_service_private_endpoint_service.this.state
}

output "target_state" {
  value = data.couchbase-capella_app_service_private_endpoint_service.this.target_state
}

output "command" {
  value = data.couchbase-capella_app_service_private_endpoint_service.this.command
}
