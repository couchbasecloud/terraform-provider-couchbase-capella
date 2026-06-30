terraform {
  required_providers {
    couchbase-capella = {
      source = "couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  authentication_token = var.auth_token
  host     = "http://localhost:8084"
}

