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

variable "outbound_private_endpoint" {
  description = "Outbound private endpoint configuration details"

  type = object({
    name         = string
    service_name = string
    hostname     = optional(string)
  })
}
