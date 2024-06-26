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

# Reference variable for GCP Config to create a network peer on GCP and use these in the create_network_peer.tf file under provider_config.
# variable "gcp_config" {
#   description = "GCP configuration details useful for network peer creation"
#
#   type = object({
#     network_name    = optional(string)
#     project_id      = optional(string)
#     cidr            = string
#     service_account = optional(string)
#   })
# }
