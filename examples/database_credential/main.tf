terraform {
  required_providers {
    couchbase-capella = {
      source = "registry.terraform.io/couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  host                 = var.host
  authentication_token = var.auth_token
}

