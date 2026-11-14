resource "device42_auto_discovery_nmap" "example" {
  name          = "example"
  server        = "hostname.domain.local"
  target        = "192.168.1.0/24"
  schedule_time = "02:00"
  schedule_days = "0,2,4"
}