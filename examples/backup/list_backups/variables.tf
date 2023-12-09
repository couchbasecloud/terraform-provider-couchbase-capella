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

variable "bucket_id" {
  description = "Capella Bucket ID"
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