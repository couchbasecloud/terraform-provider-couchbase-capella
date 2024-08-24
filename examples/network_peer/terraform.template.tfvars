auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

network_peer = {
  name          = "VPCPeerTFTestAWS"
  provider_type = "aws"
}

aws_config = {
  account_id = "123456789123"
  vpc_id     = "vpc-141f0fffff141aa00"
  region     = "us-east-1"
  cidr       = "10.0.0.0/16"
}

# Example GCP Config for creating network peer on GCP. Use this if you want to create a network peer for GCP.
# network_peer = {
#   name               = "VPCPeerTFTestGCP"
#   provider_type      = "gcp"
# }
# gcp_config = {
#     project_id = "rock-galaxy-123456"
#     network_name  = "cc-ffffffff-aaaa-1414-eeee-000000000000"
#     service_account = "service-account-name@project-id.iam.gserviceaccount.com"
#     cidr       = "10.0.0.0/16"
# }


# Example Azure Config for creating network peer on Azure. Use this if you want to create a network peer for Azure.
# network_peer = {
#   name               = "VNETPeerTFTestAZURE"
#   provider_type      = "azure"
# }
# azure_config = {
#   tenant_id       = "ffffffff-aaaa-1414-eeee-000000000000"
#   subscription_id = "ffffffff-aaaa-1414-eeee-000000000000"
#   cidr            = "10.0.0.0/16"
#   resource_group  = "test-rg"
#   vnet_id         = "test-vnet"
# }