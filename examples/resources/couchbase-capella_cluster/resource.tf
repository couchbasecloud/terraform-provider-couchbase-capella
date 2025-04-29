resource "couchbase-capella_cluster" "new_cluster" {
  organization_id               = "<organization_id>"
  project_id                    = "<project_id>"
  name                          = "<cluster_id>"
  description                   = "My first test cluster for multiple services."
  enable_private_dns_resolution = true
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "10.1.30.0/23"
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
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes =  3
      services     = [ "search","index","data", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
