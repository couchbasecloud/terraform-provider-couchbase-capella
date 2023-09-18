variable "host" {
  default = "https://cloudapi.dev.nonprod-project-avengers.com"
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
}

variable "project_name" {
  default = "terraform-couchbasecapella-project"
  description = "Project Name for Project Created via Terraform"
}