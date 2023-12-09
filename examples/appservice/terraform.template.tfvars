auth_token      = "v4-api-key-secret"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

app_service = {
  name        = "new-terraform-app-service"
  description = "My first test app service."
  nodes       = 2
  compute = {
    cpu = 2
    ram = 4
  }
}