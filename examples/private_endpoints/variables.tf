variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "endpoint_id" {
  description = "endpoint ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}