resource "device42_subnet" "myNewNetwork" {
    name        = "myNewNetwork"
    mask_bits   = 24
    network     = "192.168.100.0"
    vrf_group   = "Chicago Data Center"
}