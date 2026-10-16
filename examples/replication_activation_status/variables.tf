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
  description = "Capella Cluster ID of the cluster where this direction of the replication starts"
}

variable "replication_id" {
  description = "Replication ID. For the reverse direction of a twoWay replication, use the reverse replication ID and the target cluster ID."
}

variable "state" {
  description = "The requested status of the replication: running or paused."
  type        = string
  default     = "paused"
  validation {
    condition     = contains(["running", "paused"], var.state)
    error_message = "State can only be set to either 'running' or 'paused'."
  }
}
