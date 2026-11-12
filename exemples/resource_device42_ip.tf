data "device42_subnet" "subnet" {
  name = "Infra"
}

resource "device42_device" "host" {
  name               = var.name
  archive_on_destroy = "true"
  service_level      = var.servicelevel
}

resource "device42_ip" "ip" {
  subnet_id    = data.device42_subnet.subnet.subnet_id
  suggest_ip   = true
  available    = "no"
  vrf_group_id = data.device42_subnet.subnet.vrf_group_id
  device_id    = tonumber(device42_device.host.id)
}

resource "device42_ip" "ip2" {
  subnet_id = data.device42_subnet.subnet.subnet_id
  ip        = "192.168.1.100"
  available = "no"
}