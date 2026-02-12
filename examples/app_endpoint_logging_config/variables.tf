variable "organization_id" {
  description = "Capella Organization ID"
  type = string
}

variable "auth_token" {
  description = "Authentication API Key"
  type = string
}

variable "project_id" {
  description = "Capella Project ID"  
  type = string
}

variable "cluster_id" {
  description = "Capella Cluster ID"
  type = string
}

variable "app_service_id" {
    description = "Capella App Service ID"
    type = string
}

variable "app_endpoint_name" {
    description = "Capella App Endpoint Name"
    type = string
}

variable "app_endpoint_log_streaming_config" {
  description = "App endpoint logging config"
  type = object({
    log_level = string
    log_keys = list(string)
  })

  validation {
    condition = contains(["info", "warn", "error"], var.app_endpoint_log_streaming_config.log_level)
    error_message = "log_level must be 'info', 'warn' or 'error'."
  }

  validation {
    condition = length(var.app_endpoint_log_streaming_config.log_keys) >= 1
    error_message = "There must be one or more log_keys."
  }

  validation {
    condition = alltrue([
      for log_key in var.app_endpoint_log_streaming_config.log_keys : contains([
        "Admin",
		"Access",
		"Auth",
		"Cache",
		"Changes",
		"CRUD",
		"HTTP",
		"HTTP+",
		"Import",
		"Javascript",
		"Query",
		"Sync",
		"SyncMsg"], log_key)
    ])
    error_message = "log_key must be 'Admin', 'Access, 'Auth', 'Cache', 'Changes, 'CRUD', 'HTTP', 'HTTP+', 'Import', 'Javascript', 'Query', 'Sync' or 'SyncMsg'"
  }

}