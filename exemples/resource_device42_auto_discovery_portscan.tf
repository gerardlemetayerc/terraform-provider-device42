resource "device42_auto_discovery_portscan" "example" {
  name          = "portscan"
  server        = "hostname.domain.local"
  start_ip      = "10.0.0.1"
  schedule_days = "0,1,2,3,4"
}