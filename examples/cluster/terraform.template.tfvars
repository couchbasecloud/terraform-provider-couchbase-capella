auth_token      = "v4-api-key-secret"
organization_id = "bdb8662c-7157-46ea-956f-ed86f4c75211"
host            = "https://cloudapi.cloud.couchbase.com"

cloud_provider = {
  name   = "aws",
  region = "us-east-1"
}

cluster = {
  name               = "New Terraform Cluster"
  cidr               = "192.168.0.0/20"
  node_count         = 3
  server_version     = "7.1"
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
