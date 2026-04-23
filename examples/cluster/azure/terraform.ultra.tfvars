auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"

cloud_provider = {
  name   = "azure",
  region = "eastus"
}

couchbase_server = {
  version = "8.0"
}

cluster = {
  name               = "TF Azure Ultra"
  cidr               = "10.10.0.0/23"
  node_count         = 3
  couchbase_services = ["data", "index", "query", "search"]
  availability_zone  = "single"
}

compute = {
  cpu = 4
  ram = 16
}

disk = {
  type          = "Ultra"
  size          = 64
  iops          = 5000
  autoexpansion = true
}

support = {
  plan     = "developer pro"
  timezone = "PT"
}
