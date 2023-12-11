# Stores the cluster details in an output variable.
# Can be viewed using `terraform output cluster` command
output "cluster" {
  value = couchbase-capella_cluster.new_cluster
}

resource "couchbase-capella_cluster" "new_cluster" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  name            = var.cluster.name
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = var.cloud_provider.name
    region = var.cloud_provider.region
    cidr   = var.cluster.cidr
  }
  couchbase_server = var.couchbase_server
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

