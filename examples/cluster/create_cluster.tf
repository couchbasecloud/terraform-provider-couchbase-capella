output "existing_project" {
  value = capella_project.existing_project
}

output "new_cluster" {
  value = capella_cluster.new_cluster
}
resource "capella_project" "existing_project" {
  organization_id = var.organization_id
  name            = var.project_name
  description     = "A Capella Project that will host many Capella clusters."
}

resource "capella_cluster" "new_cluster" {
  organization_id = var.organization_id
  project_id      = capella_project.existing_project.id
  name            = var.cluster.name
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = var.cloud_provider.name
    region = var.cloud_provider.region
    cidr   = var.cluster.cidr
  }
  couchbase_server = {
    version = var.cluster.server_version
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = var.compute.cpu
          ram = var.compute.ram
        }
        disk = {
          storage = var.disk.size
          type    = var.disk.type
          iops    = var.disk.iops
        }
      }
      num_of_nodes = var.cluster.node_count
      services     = var.cluster.couchbase_services
    }
  ]
  availability = {
    "type" : var.cluster.availability_zone
  }
  support = {
    plan     = var.support.plan
    timezone = var.support.timezone
  }
}

