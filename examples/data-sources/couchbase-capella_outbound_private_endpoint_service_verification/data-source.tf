data "couchbase-capella_outbound_private_endpoint_service_verification" "service_check" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  service_name    = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
}
