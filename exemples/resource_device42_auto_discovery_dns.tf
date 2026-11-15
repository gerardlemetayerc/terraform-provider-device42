resource "device42_auto_discovery_dns" "example" {
  zonename      = "example.local"
  nameserver    = "ns1.example.local"
  schedule_time = "02:00"
  schedule_days = "0,1,2,3,4,5,6"
}