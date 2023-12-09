variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "project_name" {
  default     = "terraform-couchbasecapella-project"
  description = "Project Name for Project Created via Terraform"
}
