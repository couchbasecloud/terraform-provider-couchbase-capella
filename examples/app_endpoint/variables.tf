
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

variable "app_service_id" {
    description = "Capella App Service ID"
}

variable "bucket" {
    description = "bucket"
}

variable "name" {
    description = "name"
}

variable "delta_sync_enabled" {
    description = "delta_sync_enabled"
}

variable "scope" {
    description = "scope"
}

variable "collections" {
    description = "collections"
}

variable "cors" {
    description = "cors"
}

variable "oidc" {
    description = "oidc"
}
