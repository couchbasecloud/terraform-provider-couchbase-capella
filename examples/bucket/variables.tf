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
        name = string
        type = string
        storage_backend = string
        memory_allocationinmb = number
        conflict_resolution = string
        durability_level = string
        replicas = number
        flush = bool
        ttl = number
        eviction_policy = optional(string)
    })
}