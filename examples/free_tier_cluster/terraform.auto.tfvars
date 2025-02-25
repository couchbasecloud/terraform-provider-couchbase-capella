project_id      = "62207c06-572e-4b9e-b116-65605f446feb"

cloud_provider = {
  name   = "aws",
  region = "us-east-2"
}

cluster = {
  name               = "New Terraform Cluster modified"
  cidr               = "192.168.0.0/20"
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
