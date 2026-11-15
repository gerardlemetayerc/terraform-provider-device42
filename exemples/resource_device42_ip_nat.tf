resource "device42_ip_nat" "myNewIpNat" {
  name            = "myNewIpNat"
  ip_address_from = "10.0.0.1"
  ip_address_to   = "172.16.0.1"
  notes           = "My new NAT rule"
  protocol        = "tcp"
}