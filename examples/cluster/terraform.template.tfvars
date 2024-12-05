auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"

cloud_provider = {
  name   = "aws",
  region = "us-east-1"
}

cluster = {
  name               = "New Terraform Cluster"
  cidr               = "192.168.0.0/20"
  node_count         = 3
  couchbase_services = ["data", "index", "query"]
  availability_zone  = "multi"
}

compute = {
  cpu = 4
  ram = 16
}

disk = {
  size = 50
  type = "io2"
  iops = 5000
}

support = {
  plan     = "developer pro"
  timezone = "PT"
}

# Example cluster Config for creating a single-node cluster.
# cloud_provider = {
#   name   = "aws",
#   region = "us-east-1"
# }
#
# cluster = {
#   name               = "New Terraform Single Node Cluster"
#   cidr               = "10.4.0.0/23"
#   node_count         = 1
#   couchbase_services = ["data", "index", "query"]
#   availability_zone  = "single"
#   zones  =  ["use1-az1"]
# }
#
# compute = {
#   cpu = 2
#   ram = 8
# }
#
# disk = {
#   size = 50
#   type = "gp3"
#   iops = 3000
# }
#
# support = {
#   plan     = "basic"
#   timezone = "PT"
# }