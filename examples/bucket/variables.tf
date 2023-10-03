# variable "host" {
#   description = "The Host URL of Couchbase Cloud."
# }
#
# variable "organization_id" {
#   description = "Capella Organization ID"
# }
#
# variable "auth_token" {
#   description = "Authentication API Key"
#   sensitive   = true
# }
#
# variable "project_id" {
#   description = "Capella Project ID"
# }
#
# variable "cluster_id" {
#   description = "Capella Cluster ID"
# }


variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
  default = "c3JVVEJkTWNVNWtxekFxaWJBZFNvVmYwNlN6dnhidXU6TTRMZ0VrJWVYcmwwUkBlZ2c1Smhtc0RuaEVvMWlLUWlmdnBDVCVaeDQyS0FFTUJTQW9zYzhVckJIRWZpeWJuSA=="
}
variable "host" {
  description = "The Host URL of Couchbase Cloud."
  default = "https://cloudapi.dev.nonprod-project-avengers.com"
}

variable "organization_id" {
  description = "Capella Organization ID"
  default = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
}

variable "project_id" {
  description = "Capella Project ID"
  default = "f14134f2-7943-4e7b-b2c5-fc2071728b6e"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
  default = "146574f7-349c-4ca0-9714-5b2d3278a901"
}
