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

variable "app_endpoint_name" {
  description = "Capella App Endpoint Name"
}

variable "scope" {
  description = "Scope name containing the collection"
}

variable "collection" {
  description = "Collection name for which the access function is defined"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "access_control_function" {
  description = "JavaScript access control function"
  type        = string
  default     = "function (doc, oldDoc, meta) {channel(doc.channels); }"
} 