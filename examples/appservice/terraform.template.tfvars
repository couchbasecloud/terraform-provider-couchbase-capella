auth_token      = "v4-api-key-secret"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

app_service = {
  name        = "new-terraform-app-service"
  description = "My first test app service."
  nodes       = 2
  # load_balancer_cidr = "10.1.0.0/24" # Azure App Services only
  compute = {
    cpu = 2
    ram = 4
  }
}