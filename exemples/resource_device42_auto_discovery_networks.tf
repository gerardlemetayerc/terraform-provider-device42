resource "device42_auto_discovery_networks" "example" {
  name          = "example"
  server        = "hostname.domain.local"
  schedule_days = "0,2,4"
}