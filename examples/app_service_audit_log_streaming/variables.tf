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

variable "streaming_enabled" {
  description = "Set to true to stream audit logs. Set to false to pause streaming and keep the configuration."
  type        = bool
}

variable "disabled_app_endpoints" {
  description = "App Endpoints to exclude from audit log streaming"
  type        = set(string)
  default     = null
}

# Credentials for the log collector. This example uses Datadog.
# To use another log collector, uncomment its variables, set them in terraform.tfvars,
# and change output_type and credentials in create_app_service_audit_log_streaming.tf.

variable "datadog_url" {
  description = "Datadog log ingestion URL"
  sensitive   = true
}

variable "datadog_api_key" {
  description = "Datadog API key"
  sensitive   = true
}

# variable "sumologic_url" {
#   description = "Sumo Logic signed URL for log ingestion"
#   sensitive   = true
# }

# variable "elastic_url" {
#   description = "Elasticsearch log ingestion URL. Only Elasticsearch 8 and newer is supported."
#   sensitive   = true
# }

# variable "elastic_user" {
#   description = "Elasticsearch username for HTTP Basic Authentication"
#   sensitive   = true
# }

# variable "elastic_password" {
#   description = "Elasticsearch password for HTTP Basic Authentication"
#   sensitive   = true
# }

# variable "loki_url" {
#   description = "Grafana Loki log ingestion URL"
#   sensitive   = true
# }

# variable "loki_user" {
#   description = "Grafana Loki username for HTTP Basic Authentication"
#   sensitive   = true
# }

# variable "loki_password" {
#   description = "Grafana Loki password for HTTP Basic Authentication"
#   sensitive   = true
# }

# variable "splunk_url" {
#   description = "Splunk HTTP Event Collector (HEC) endpoint URL"
#   sensitive   = true
# }

# variable "splunk_token" {
#   description = "Splunk HTTP Event Collector (HEC) token"
#   sensitive   = true
# }

# variable "dynatrace_url" {
#   description = "Dynatrace log ingestion URL"
#   sensitive   = true
# }

# variable "dynatrace_api_token" {
#   description = "Dynatrace API token"
#   sensitive   = true
# }

# variable "generic_http_url" {
#   description = "Custom HTTP log collector URL"
#   sensitive   = true
# }

# variable "generic_http_user" {
#   description = "Optional username for HTTP Basic Authentication"
#   sensitive   = true
#   default     = null
# }

# variable "generic_http_password" {
#   description = "Optional password for HTTP Basic Authentication"
#   sensitive   = true
#   default     = null
# }
