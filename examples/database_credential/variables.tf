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

variable "database_credential" {
  type = object({
    database_credential_name = string
    password                 = optional(string)
  })
}

variable "access" {
  type = list(object({
    privileges = list(string)
    resources = optional(object({
      buckets = list(object({
        name = string
        scopes = optional(list(object({
          name        = string
          collections = optional(list(string))
        })))
      }))
    }))
  }))
}
