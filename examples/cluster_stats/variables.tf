variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "host" {
  description = "Capella API Host"
  default     = "https://cloudapi.cloud.couchbase.com"
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
