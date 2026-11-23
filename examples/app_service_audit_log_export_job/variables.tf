variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "app_service_id" {
  description = "Capella App Service ID"
}

variable "start" {
  description = "Start date and time of the audit logs to export (RFC3339)"
}

variable "end" {
  description = "End date and time of the audit logs to export (RFC3339). Must be at least 15 minutes in the past."
}
