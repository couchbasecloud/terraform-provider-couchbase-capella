variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
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

variable "version_id" {
  description = "ID of the version to upgrade to, from the available_version_ids output"
}

variable "queue_at" {
  description = "RFC3339 time to queue the upgrade at, in the future and no more than 7 days ahead. Leave null to upgrade now."
  default     = null
}
