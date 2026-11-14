resource "device42_auto_discovery_blade" "example" {
  name          = "example"
  server        = "hostname.domain.local"
  snmp_retries  = 3
  snmp_timeout  = 5
  schedule_days = "0,2,4"
}