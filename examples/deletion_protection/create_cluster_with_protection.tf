# This example demonstrates how to create a cluster with deletion protection enabled.
# When deletion_protection is true, terraform destroy will be blocked until the
# attribute is set to false.

output "protected_cluster" {
  value = couchbase-capella_cluster.protected_cluster
}

resource "couchbase-capella_cluster" "protected_cluster" {
  organization_id     = var.organization_id
  project_id          = var.project_id
  name                = var.cluster.name
  description         = "Cluster with deletion protection enabled."
  deletion_protection = true

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

