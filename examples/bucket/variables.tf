variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

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

variable "bucket" {
  description = "Bucket configuration details useful for creation"

  type = object({
    name                  = string
    type                  = optional(string)
    storage_backend       = optional(string)
    memory_allocationinmb = optional(number)
    conflict_resolution   = optional(string)
    durability_level      = optional(string)
    replicas              = optional(number)
    flush                 = optional(bool)
    ttl                   = optional(number)
    eviction_policy       = optional(string)
  })
}