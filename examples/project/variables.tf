variable "host" {
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
