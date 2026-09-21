variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "bucket_id" {
  description = "The bucket ID as returned by Capella, which is the URL-compatible base64 encoding of the bucket name, not the bucket name itself."
}

variable "backup_id" {
  description = "The ID of the backup to export"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
