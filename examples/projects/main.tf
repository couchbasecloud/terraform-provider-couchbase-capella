terraform {
  required_providers {
    capella = {
      source = "hashicorp.com/couchabasecloud/capella"
    }
  }
}

provider "capella" {
  host     = "hostname of the capella"
  bearer_token = "capella api key bearer token"
}

data "capella_projects" "example" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
}


output "example_projects" {
  value = data.capella_projects.example
}