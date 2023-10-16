package testing

var Cfg string = `
variable "host" {
  default = "https://cloudapi.dev.nonprod-project-avengers.com"
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  default = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  description = "Capella Organization ID"
}

variable "auth_token" {
  default = "SVJMcDhxUXdIaUY0TmkzSUlibEgwblBCYTRveDBwOEk6MmdyS2IxIVNjM09nczkjOWdsS1JRIyMycEBnTFV3Y3pLb0JZa0pjdm1seTRQNUR6Q0VPWVlxVTVmSWJjUVNEVw=="
  description = "Authentication API Key"
  sensitive   = true
}
provider "capella" {
  host                 = var.host
  authentication_token = var.auth_token
}
`
