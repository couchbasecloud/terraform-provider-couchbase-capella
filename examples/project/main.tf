terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchabasecloud/capella"
    }
  }
}

provider "capella" {
  host                 = var.host
  authentication_token = var.auth_token
}

resource "capella_project" "my_new_project" {
  organization_id = var.organization_id
  name            = var.project_name
  description     = "A Capella Project that will host many Capella clusters."
}

data "capella_projects" "existing_projects" {
  organization_id = var.organization_id
}
