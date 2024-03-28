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

variable "bucket_id" {
  description = "Capella Bucket ID"
}

variable "scope_name" {
  description = "Capella Scope Name"
}

variable "collection" {
  description = "Collection configuration details useful for creation"

  type = object({
    collection_name = string
    max_ttl         = optional(number)
  })
}
