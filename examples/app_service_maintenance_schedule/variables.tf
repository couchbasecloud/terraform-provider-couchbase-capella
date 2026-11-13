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

variable "maintenance_id" {
  description = "ID of a pending App Service maintenance, from the reschedulable_maintenance_ids output"
}

variable "queue_at" {
  description = "RFC3339 time to queue the maintenance at. It must be in the future and no later than the maintenance reschedule_by."
}
