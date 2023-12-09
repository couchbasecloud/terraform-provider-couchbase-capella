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
  configuration_type = "multiNode"
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

couchbase_server = {
  version = 7.2
}
