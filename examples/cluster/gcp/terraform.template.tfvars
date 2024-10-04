auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cloud_provider = {
  name   = "gcp",
  region = "us-east1"
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
  storage = 64
  type    = "pd-ssd"
}

support = {
  plan     = "developer pro"
  timezone = "PT"
}

