project_id      = "0f3952be-5ddf-47ab-862c-0e29afa317b9"

cloud_provider = {
  name   = "aws",
  region = "us-east-1"
}

cluster = {
  name               = "New Terraform Cluster"
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
