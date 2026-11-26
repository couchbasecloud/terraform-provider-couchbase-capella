auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

outbound_private_endpoint = {
  name         = "cng-prod"
  service_name = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
  # Leave hostname out to have Capella assign one.
  hostname = "cng.acme.example.com"
}
