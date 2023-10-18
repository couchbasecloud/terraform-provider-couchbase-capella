variable "host" {
  description = "The Host URL of Couchbase Cloud."
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

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "comment" {
  description = "comment describing the allowlist details"
}

variable "cidr" {
  description = "CIDR that will have access to the cluster"
}

variable "expires_at" {
  description = "timestamp when the allowlist expires"
}
