data "device42_subnet_suggested" "new" {
  parent_subnet_id = data.device42_subnet.searchSubnet.subnet_id
  mask_bits        = 26
}

resource "device42_subnet" "child" {
  network   = data.device42_subnet_suggested.new.ip
  mask_bits = data.device42_subnet_suggested.new.mask_bits
  name      = "new child subnet"
}