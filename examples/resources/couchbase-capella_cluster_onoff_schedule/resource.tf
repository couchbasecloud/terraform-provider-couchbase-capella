resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
  organization_id = "organization_id"
  project_id      = "project_id"
  cluster_id      = "cluster_id"
  timezone        = "timezone(ex:US/Hawaii)"
  days = [
    for day in var.days : {
      state = "off"
      day   = "Monday"
      from =  {
                  hour   = 12
                  minute = 30
                }
       to = {
                  hour   = 14
                  minute = 30
                }
      from  = "from time"
      to    = "to time"
    }
  ]
}