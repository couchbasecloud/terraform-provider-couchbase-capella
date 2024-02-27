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

variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "scope" {
  description = "Scope configuration details useful for creation"

  type = object({
    name                       = string
  })
}
