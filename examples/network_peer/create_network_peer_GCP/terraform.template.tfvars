auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

network_peer = {
  name               = "VPCPeerTFTestAWS"
  provider_type      = "gcp"
}

GCP_config = {
    network_name = "VPCPeerTestGCP"
    project_id   = "rock-galaxy-123456"
    cidr       = "10.0.0.0/16"
    service_account = "service-account-name@project-id.iam.gserviceaccount.com"
}