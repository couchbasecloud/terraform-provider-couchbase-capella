variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "event_id" {
  description = "Capella Event ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}