data "device42_vlan" "example" {
  vlan_id = 1234
}

output "vlan_number" {
  value = data.device42_vlan.example.number
}
