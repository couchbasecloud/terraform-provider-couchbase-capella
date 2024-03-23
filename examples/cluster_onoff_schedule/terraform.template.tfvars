auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

days = [{
  day   = "monday"
  state = "custom"
  from  = {
    hour   = 12
    minute = 30
    }
  to  = {
      hour   = 14
      minute = 30
      }
},
{
  day   = "tuesday"
  state = "custom"
  from  = {
    hour   = 12
    minute = 30
    }
  to  = {
      hour   = 14
      minute = 30
      }
},
{
  day   = "wednesday"
  state = "on"
},
{
  day   = "thursday"
  state = "on"
},
{
  day   = "friday"
  state = "custom"
  from  = {
    hour   = 12
    minute = 30
    }
  to  = {
      hour   = 14
      minute = 30
      }
},
{
  day   = "saturday"
  state = "off"
},
{
  day   = "sunday"
  state = "off"
}
]