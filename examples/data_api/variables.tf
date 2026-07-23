variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "data_api" {
  description = "Data API configuration for the cluster"

  type = object({
    enable_data_api        = bool
    enable_network_peering = bool
  })
}
