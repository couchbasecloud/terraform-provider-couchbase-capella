variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "resource_group_name" {
  description = "resource group name"
}

variable "virtual_network" {
  description = "virtual network"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
