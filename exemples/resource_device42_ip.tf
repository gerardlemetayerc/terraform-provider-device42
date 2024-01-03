data "device42_subnet" "subnet" {
  name    = "Infra"
}

resource "device42_device" "host" {
   name  				        = var.name
   archive_on_destroy 	= true
   service_level        = var.servicelevel
}

resource "device42_ip" "ip" {
  subnet 	      = data.device42_subnet.subnet.name
  suggest_ip    = true
  available 	  = "no"
  vrf_group_id 	= data.device42_subnet.subnet.vrf_group_id
  device_id     = resource.device42_device.host.id
}

resource "device42_ip" "ip2" {
  subnet 	      = data.device42_subnet.subnet.name
  ip            = "192.168.1.100"
  available 	  = "no"
}