terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchabasecloud/capella"
    }
  }
}

provider "capella" {
  host     = var.host
  bearer_token = var.auth_token
}

resource "capella_project" "example" {
  organization_id = var.organization_id
  name = var.project_name
  description = "A Capella Project that will host many Capella clusters."
}



