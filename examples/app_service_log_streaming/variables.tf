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

variable "output_type" {
  description = "Type of Log Collector to stream logs to"
  type = string
}

variable "credentials" {
  description = "Credentials for the Log Collector to stream logs to. The structure of this variable will depend on the output_type specified."
  type = map(any)
  sensitive   = true
}