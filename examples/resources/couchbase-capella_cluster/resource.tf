resource "couchbase-capella_cluster" "new_cluster" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  name            = "Terraform Test Cluster"
  description     = "Test cluster created with Terraform"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "10.200.250.0/23"
  }
  couchbase_server = {
    version = "7.6"
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
          type    = "io2"
          iops    = 5000
        }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
