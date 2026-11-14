resource "device42_auto_discovery_cloudaccount" "example" {
  name          = "cloudaccount"
  type          = "OpenStack"
  server        = "hostname.domain.local"
  username      = "admin"
  schedule_days = "0,1,2,3,4"
}