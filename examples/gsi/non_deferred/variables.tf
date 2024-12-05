variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
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

variable "bucket_name" {
  description = "Bucket Name"
}

variable "scope_name" {
  description = "Scope Name"
}

variable "collection_name" {
  description = "collection Name"
}

variable "index_name" {
  description = "index Name"
}

variable "index_keys" {
  description = "index keys"
}
