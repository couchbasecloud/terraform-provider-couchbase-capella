terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchabasecloud/capella"
    }
  }
}

provider "capella" {
  host     = var.host
  authentication_token = var.auth_token
}

resource "capella_allowlist" "example" {
  organization_id = var.organization_id
  cidr = "0.00.00.00"
  comment = "Allow access from any ip address"
  expiresAt = "2023-05-14T21:49:58.465Z"
}
