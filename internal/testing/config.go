package testing

var Cfg string = `
variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Project ID"
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
