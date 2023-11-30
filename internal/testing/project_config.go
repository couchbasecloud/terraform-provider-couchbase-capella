package testing

var ProjectCfg string = `
variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
provider "capella" {
  host                 = var.host
  authentication_token = var.auth_token
}
`
