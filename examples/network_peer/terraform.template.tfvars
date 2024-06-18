auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

network_peer = {
  name               = "VPCPeerTFTestAWS"
  provider_type      = "aws"
}

aws_config = {
    account_id = "123456789123"
    vpc_id     = "vpc-141f0fffff141aa00"
    region     = "us-east-1"
    cidr       = "10.0.0.0/16"
}