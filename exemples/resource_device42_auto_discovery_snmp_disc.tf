resource "device42_auto_discovery_snmp_disc" "example" {
  name          = "mySnmpJob"
  server        = "hostname.domain.local"
  snmp_retries  = "0"
  snmp_timeout  = "60"
  schedule_days = "0,1,2,3,4,5,6"
}