data "device42_vlan" "prod" {
  number = 100
}

data "device42_subnets" "in_vlan" {
  parent_vlan_id = data.device42_vlan.prod.vlan_id
}

output "subnet_networks" {
  value = [for s in data.device42_subnets.in_vlan.subnets : s.network]
}