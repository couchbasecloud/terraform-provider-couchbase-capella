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
  description = "Capella Cluster ID of the source cluster"
}

variable "source_bucket_id" {
  description = "Bucket ID (not name) of the source bucket"
}

variable "target" {
  description = "Target of the replication. For a Capella target, cluster is the cluster ID and bucket is the bucket ID."
  type = object({
    cluster = string
    bucket  = string
    type    = optional(string)
  })
}

variable "replication" {
  description = "Replication settings"
  type = object({
    direction           = optional(string)
    priority            = optional(string)
    network_usage_limit = optional(number)
  })
  default = {}
}

variable "mobile_replication" {
  description = "Bucket IDs (not names) for the replication with mobile sync mode and conflict logging. The log bucket is on the source cluster."
  type = object({
    source_bucket_id = string
    target_bucket_id = string
    log_bucket_id    = string
  })
}

variable "conflict_logging_disabled" {
  description = "Set to true to keep the conflict logging configuration without logging conflicts."
  type        = bool
  default     = false
}
