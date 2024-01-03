datasource "device42_subnet" "searchSubnet" {
  subnet_id    = 1
}

datasource "device42_subnet" "searchSubnet2" {
  network         = "192.168.1.0"
  mask_bits       = 24
  vrf_group_name  = "Customer A"
}

datasource "device42_subnet" "searchSubnet3" {
  name    = "My subnet name"
}