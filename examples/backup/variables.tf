variable "host" {
  default = "https://cloudapi.dev.nonprod-project-avengers.com"
  description = "The Host URL of Couchbase Cloud."
}

variable "organization_id" {
  default = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  description = "Capella Organization ID"
}

variable "auth_token" {
  default = "SVJMcDhxUXdIaUY0TmkzSUlibEgwblBCYTRveDBwOEk6MmdyS2IxIVNjM09nczkjOWdsS1JRIyMycEBnTFV3Y3pLb0JZa0pjdm1seTRQNUR6Q0VPWVlxVTVmSWJjUVNEVw=="
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
    target_cluster_id = string
    source_cluster_id = string
    services = list(string)
    force_updates = optional(bool)
    auto_remove_collections = optional(bool)
    filter_keys = optional(string)
    filter_values = optional(string)
    include_data = optional(string)
    exclude_data = optional(string)
    map_data = optional(string)
    replace_ttl = optional(string)
    replace_ttl_with = optional(string)
    restore_times = number
  })
}