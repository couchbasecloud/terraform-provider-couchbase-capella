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

resource "capella_project" "sampleProject" {
  organization_id = var.organization_id
  name = var.project_name
  description = "A Capella Project that will host many Capella clusters."
}

resource "capella_cluster" "sampleCluster" {
  organization_id = var.organization_id
  project_id = capella_project.sampleProject.id
  name =  "Sample CLuster"
  description = "My first test aws cluster for multiple services."
  cloud_provider = {
    type =  "aws"
    region = "us-east-1"
    cidr = "10.1.88.0/23"
  }
  couchbase_server = {
    version = "7.1"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type = "io2"
          iops = 3000
        }
      }
      num_of_nodes = 3
      services = ["data", "query", "index"]
    }
  ]
  availability = {
    "type": "multi"
  }
  support = {
    plan =  "developer pro"
    timezone = "PT"
  }
}

