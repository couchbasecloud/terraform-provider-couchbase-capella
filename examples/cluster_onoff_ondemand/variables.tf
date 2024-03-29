variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Project Name for Project Created via Terraform"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "cluster_onoff_ondemand" {
    description = "Provides the means to turn the given cluster to on or off state"

    type = object({
        state = string
        turn_on_linked_app_service = optional(bool)
      })
}