project_id      = "3b74be24-9787-49f6-bb4e-2ec8ea01e7ba"

cloud_provider = {
  name   = "aws",
  region = "us-east-2"
}

cluster = {
  name               = "New Terraform Cluster name"
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
