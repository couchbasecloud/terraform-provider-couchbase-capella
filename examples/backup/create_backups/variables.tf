variable "organization_id" {
  description = "Capella Organization ID"
}

variable "auth_token" {
  description = "Authentication API Key"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "backup" {
  description = "Backup configuration details useful for creation"
}

variable "restore" {
  type = object({
    target_cluster_id       = string
    source_cluster_id       = string
    services                = list(string)
    force_updates           = optional(bool)
    auto_remove_collections = optional(bool)
    filter_keys             = optional(string)
    filter_values           = optional(string)
    include_data            = optional(string)
    exclude_data            = optional(string)
    map_data                = optional(string)
    replace_ttl             = optional(string)
    replace_ttl_with        = optional(string)
    restore_times           = number
  })
  default = null
}

variable "bucket" {
  description = "Bucket configuration details useful for creation"

  type = object({
    name                       = string
    type                       = optional(string)
    storage_backend            = optional(string)
    memory_allocation_in_mb    = optional(number)
    bucket_conflict_resolution = optional(string)
    durability_level           = optional(string)
    replicas                   = optional(number)
    flush                      = optional(bool)
    time_to_live_in_seconds    = optional(number)
    eviction_policy            = optional(string)
  })
}
