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

provider "couchbase-capella" {
  host                 = var.host
  authentication_token = var.auth_token
}

resource "couchbase-capella_app_service_private_endpoint_service" "this" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  app_service_id  = var.app_service_id
}

output "state" {
  value = couchbase-capella_app_service_private_endpoint_service.this.state
}
