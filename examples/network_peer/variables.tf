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

variable "aws_config" {
  description = "AWS configuration details useful for network peer creation"

  type = object({
    account_id = optional(string)
    vpc_id     = optional(string)
    cidr       = string
    region     = optional(string)
  })
}

# variable "gcp_config" {
#   description = "GCP configuration details useful for network peer creation"
#
#   type = object({
#     network_name    = string
#     project_id      = string
#     cidr            = string
#     service_account = string
#   })
# }

# locals {
#   config_check = (var.aws_config != null && var.gcp_config != null) ?
#     error("Only one of aws_config or gcp_config should be provided") :
#     (var.aws_config != null ? var.aws_config : var.gcp_config)
# }

variable "host" {
  description = "The Host URL of Couchbase Cloud."
}