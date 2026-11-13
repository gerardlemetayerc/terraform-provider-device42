resource "device42_subnet_suggested" "new" {
  parent_subnet_id = data.device42_subnet.searchSubnet.subnet_id
  mask_bits        = 26
  name             = "new child subnet"
}

output "suggested_subnet_id" {
  value = device42_subnet_suggested.new.subnet_id
}