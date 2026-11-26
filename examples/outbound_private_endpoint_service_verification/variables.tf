variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Capella Project ID"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "service_name" {
  description = "The private link service to check: an AWS endpoint service name, a GCP service attachment URI or an Azure Private Link service alias"
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}
