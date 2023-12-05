terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchbasecloud/couchbase-capella"
    }
  }
}

provider "capella" {
  host                 = var.host
  authentication_token = var.auth_token
}