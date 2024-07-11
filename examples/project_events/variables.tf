variable "organization_id" {
description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}