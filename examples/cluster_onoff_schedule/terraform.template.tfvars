auth_token      = "<v4-api-key-secret>"
organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

cluster_onoff_schedule = {
  timezone = "US/Pacific"
}

days = [
  {
    day   = "monday"
    state = "custom"
    from = {
      hour   = 12
      minute = 30
    }
    to = {
      hour   = 14
      minute = 30
    }
  },
  {
    day   = "tuesday"
    state = "custom"
    from = {
      hour = 12
    }
    to = {
      hour   = 19
      minute = 30
    }
  },
  {
    day   = "wednesday"
    state = "on"
  },
  {
    day   = "thursday"
    state = "custom"
    from = {
      hour   = 12
      minute = 30
    }
  },
  {
    day   = "friday"
    state = "custom"
    from = {

    }
    to = {
      hour   = 12
      minute = 30
    }
  },
  {
    day   = "saturday"
    state = "custom"
    from = {
      hour   = 12
      minute = 30
    }
    to = {
      hour = 14
    }
  },
  {
    day   = "sunday"
    state = "off"
  }
]