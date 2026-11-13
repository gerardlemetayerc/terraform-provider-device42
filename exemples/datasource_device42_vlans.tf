data "device42_vlans" "all" {}

output "vlan_names" {
  value = [for v in data.device42_vlans.all.vlans : v.name]
}