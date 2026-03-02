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

variable "app_service_id" {
  description = "Capella App Service ID"
}

variable "output_type" {
  description = "Type of Log Collector to stream logs to"
  type        = string
}

variable "credentials" {
  description = "Credentials for the Log Collector to stream logs to. The structure of this variable will depend on the output_type specified."
  sensitive   = true
  type = object({
    datadog = optional(object({
      url     = string
      api_key = string
    }))
    dynatrace = optional(object({
      url       = string
      api_token = string
    }))
    elastic = optional(object({
      url      = string
      user     = string
      password = string
    }))
    generic_http = optional(object({
      url      = string
      user     = optional(string)
      password = optional(string)
    }))
    loki = optional(object({
      url      = string
      user     = string
      password = string
    }))
    splunk = optional(object({
      url          = string
      splunk_token = string
    }))
    sumologic = optional(object({
      url = string
    }))
  })
}