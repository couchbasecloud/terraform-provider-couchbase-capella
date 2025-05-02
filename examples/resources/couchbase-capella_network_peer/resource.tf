# example for AWS config
resource "couchbase-capella_network_peer" "new_network_peer" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  name            = "VPCPeerTFTestAWS"
  provider_type   = "aws"
  provider_config = {
    aws_config = {
      account_id = "123456789123"
      vpc_id     = "vpc-141f0fffff141aa00ff"
      cidr       = "10.1.0.0/23"
      region     = "us-east-1"
    }
  }
}


# example for GCP config
resource "couchbase-capella_network_peer" "new_network_peer" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  name            = "VPCPeerTFTestGCP"
  provider_type   = "gcp"
  provider_config = {
    aws_config = {
       cidr            = "10.0.4.0/23"
       network_name    = "cc-ffffffff-aaaa-1414-eeee-000000000000"
       project_id      = "test-project-id"
       service_account = "service-account-name@project-id.iam.gserviceaccount.com"
    }
  }
}


# example for Azure config
resource "couchbase-capella_network_peer" "new_network_peer" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  name            = "VPCPeerTFTestAzure"
  provider_type   = "azure"
  provider_config = {
    aws_config = {
      cidr            = "10.6.0.0/16"
      resource_group  = "test-rg"
      subscription_id = "ffffffff-aaaa-1414-eeee-00000000000"
      tenant_id       = "ffffffff-aaaa-1414-eeee-00000000000"
      vnet_id         = "test-vnet"
    }
  }
}