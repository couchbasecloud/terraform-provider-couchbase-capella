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

variable "tenant_id" {
  description = "Azure Tenant ID"
}

variable "vnet_id" {
  description = "Azure virtual network name"
}

variable "subscription_id" {
  description = "Azure Subscription ID"
}

variable "resource_group" {
  description = "Azure resource group name holding the resource you’re connecting with Capella"
}

variable "vnet_peering_service_principal" {
  description = "Azure enterprise application object ID for the Capella service principal"
}

