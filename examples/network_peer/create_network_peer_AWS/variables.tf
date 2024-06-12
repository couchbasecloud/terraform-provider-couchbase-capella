variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "network_peer" {
  description = "Network Peer configuration details useful for creation"

  type = object({
    name               = string
    provider_type      = string
  })
}

variable "AWS_config" {
  description = "AWS configuration details useful for network peer creation"

  type = object({
    account_id = optional(string)
    vpc_id     = optional(string)
    cidr       = string
    region     = optional(string)
  })
}