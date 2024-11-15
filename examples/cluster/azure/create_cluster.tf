output "new_cluster" {
  value = couchbase-capella_cluster.new_cluster
}

output "cluster_id" {
  value = couchbase-capella_cluster.new_cluster.id
}

resource "couchbase-capella_cluster" "new_cluster" {
  organization_id = var.organization_id
  project_id      = var.project_id
  name            = var.cluster.name
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = var.cloud_provider.name
    region = var.cloud_provider.region
    cidr   = var.cluster.cidr
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = var.compute.cpu
          ram = var.compute.ram
        }
        disk = {
          type          = var.disk.type
          autoexpansion = var.disk.autoexpansion
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

