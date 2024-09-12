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

variable "subnet_ids" {
  description = "subnet IDs"
  type        = list(string)
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

