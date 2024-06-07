variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "vpc_id" {
  description = "VPC ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

