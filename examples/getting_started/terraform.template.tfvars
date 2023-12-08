auth_token = "<v4-api-key-secret>"
host       = "https://cloudapi.cloud.couchbase.com"

organization_id = "<organization_id>"
project_name    = "My First Terraform Project"

cloud_provider = {
  name   = "aws",
  region = "us-east-1"
}

cluster = {
  name               = "My First Terraform Cluster"
  cidr               = "10.250.250.0/23"
  node_count         = 3
  server_version     = "7.1"
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

database_credential_name = "terraform_db_credential"
password                 = "Secret12$#"

access = [
  {
    privileges = ["data_writer"]
    resources = {
      buckets = [{
        name = "new_terraform_bucket"
        scopes = [
          {
            name        = "_default"
            collections = ["_default"]
          }
        ]
      }]
    }
  },
  {
    privileges = ["data_reader"]
  }
]

comment    = "Allow access from a public IP"
cidr       = "8.8.8.8/32"
expires_at = "2043-11-30T23:59:59.465Z"

bucket = {
  name                       = "new_terraform_bucket"
  type                       = "couchbase"
  storage_backend            = "couchstore"
  memory_allocation_in_mb    = 100
  bucket_conflict_resolution = "seqno"
  durability_level           = "none"
  replicas                   = 1
  flush                      = false
  time_to_live_in_seconds    = 0
  eviction_policy            = "fullEviction"
}

user = {
  email = "johndoe@couchbase.com"
  name  = "John Doe"
}

apikey = {
  name               = "My First Terraform API Key"
  description        = "A Capella V4 API Key"
  allowed_cidrs      = ["10.1.42.0/23", "10.1.43.0/23"]
  organization_roles = ["organizationMember"]
  expiry             = 180
}

app_service = {
  name        = "new-terraform-app-service"
  description = "My first test app service."
  nodes       = 2
  compute = {
    cpu = 2
    ram = 4
  }
}
