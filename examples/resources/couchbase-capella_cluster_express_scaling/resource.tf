resource "couchbase-capella_cluster_express_scaling" "example" {
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  express_scaling = true

  timeouts {
    create = "4h"
    update = "4h"
  }
}
