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

variable "vpc_id" {
  description = "AWS VPC ID used to generate the private endpoint command"
}

variable "subnet_ids" {
  description = "AWS subnet IDs used to generate the private endpoint command"
  type        = list(string)
}
