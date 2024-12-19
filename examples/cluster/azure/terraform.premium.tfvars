auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"

cloud_provider = {
  name   = "azure",
  region = "eastus"
}

cluster = {
  name               = "New Terraform Azure Cluster 6"
  cidr               = "10.0.6.0/23"
  node_count         = 3
  couchbase_services = ["data"]
  availability_zone  = "single"
}

compute = {
  cpu = 4
  ram = 16
}

disk = {
  type          = "P6"
  autoexpansion = true
}

support = {
  plan     = "basic"
  timezone = "PT"
}