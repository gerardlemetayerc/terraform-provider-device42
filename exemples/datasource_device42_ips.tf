data "device42_ips" "subnet_ips" {
  subnet_id = data.device42_subnet.searchSubnet.subnet_id
}

output "ip_addresses" {
  value = [for i in data.device42_ips.subnet_ips.ips : i.ip]
}