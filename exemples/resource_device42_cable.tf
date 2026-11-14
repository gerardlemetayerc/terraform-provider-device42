resource "device42_cable" "myNewCable" {
  cable_id          = "myNewCable"
  origin_type       = "netport"
  origin_id         = "1016"
  end_type          = "netport"
  end_id            = "1017"
  cable_length      = 7
  cable_length_units = "m"
  notes             = "link between two netports"
  vendor            = "Cisco"
  tags              = "prod,ethernet"
}