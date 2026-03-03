auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"
app_service_id  = "<app_service_id>"

output_type = "generic_http"
credentials = {
  generic_http = {
    user = "example_user"
    password = "example_password"
    url = "https://couchbase.com"
  }
}