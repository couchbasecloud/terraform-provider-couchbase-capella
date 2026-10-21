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

variable "endpoint_id" {
  description = "Cloud service provider private endpoint ID, for example an AWS VPC endpoint ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
