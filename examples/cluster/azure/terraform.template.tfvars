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


# Example for Ultra disk type Azure Cluster:

# cloud_provider = {
#   name   = "azure",
#   region = "eastus"
# }
# cluster = {
#   name               = "TF Azure Ultra"
#   cidr               = "10.10.0.0/23"
#   node_count         = 3
#   couchbase_services = ["data", "index", "query", "search"]
#   availability_zone  = "single"
# }
#
# compute = {
#   cpu = 4
#   ram = 16
# }
#
# disk = {
#   type          = "Ultra"
#   storage       = 128
#   iops          = 5000
#   autoexpansion = true
# }
#
# support = {
#   plan     = "developer pro"
#   timezone = "PT"
# }
