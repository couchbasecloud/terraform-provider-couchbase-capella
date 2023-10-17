# Stores the API key details in an output variable.
# Can be viewed using `terraform output apikey` command
output "apikey" {
  value     = capella_apikey.new_apikey
  sensitive = true
}

resource "capella_apikey" "new_apikey" {
  organization_id    = var.organization_id
  name               = var.apikey.name
  organization_roles = var.apikey.organization_roles
  allowed_cidrs      = var.apikey.allowed_cidrs
  resources = [
    {
      id    = capella_project.new_project.id
      roles = ["projectManager", "projectDataReader"]
      type  = "project"
    }
  ]
}

