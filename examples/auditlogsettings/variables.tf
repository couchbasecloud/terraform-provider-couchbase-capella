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

variable "host" {
  description = "The Host URL of Couchbase Cloud."
}

variable "auditlogsettings" {
  description = "configure cluster audit log settings"

  type = object({
    auditenabled = bool
    enabledeventids = list(number)
    disabledusers = list(object({
      name  = string
      domain    = string
    }))
  })
}