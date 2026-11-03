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

variable "streaming_enabled" {
  description = "Set to true to stream audit logs. Set to false to pause streaming and keep the configuration."
  type        = bool
}

variable "disabled_app_endpoints" {
  description = "App Endpoints to exclude from audit log streaming"
  type        = set(string)
  default     = null
}

variable "datadog_url" {
  description = "Datadog log ingestion URL"
  sensitive   = true
}

variable "datadog_api_key" {
  description = "Datadog API key"
  sensitive   = true
}
