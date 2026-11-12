data "device42_ip" "lookup" {
  ip = "10.1.10.5"
}

output "device_id" {
  value = data.device42_ip.lookup.device_id
}

output "subnet" {
  value = data.device42_ip.lookup.subnet
}