resource "device42_auto_discovery_ipmi" "example" {
  name           = "ipmi"
  server         = "hostname.domain.local"
  ip_start       = "10.0.0.1"
  ip_end         = "10.0.0.100"
  discovery_type = "redfish"
  schedule_days  = "0,1,2,3,4"
}