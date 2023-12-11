terraform {
  required_providers {
    couchbase-capella = {
      source = "registry.terraform.io/couchbasecloud/couchbase-capella"
    }
  }
}

provider "couchbase-capella" {
  authentication_token = var.auth_token
  host                 = var.host_url
}
