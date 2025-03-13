terraform {
  required_providers {
    couchbase-capella = {
      source = "hashicorp.com/couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  authentication_token = var.auth_token
}

