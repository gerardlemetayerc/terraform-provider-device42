---
page_title: "device42_IP Resource - terraform-provider-device42"
subcategory: "network management"
description: |-
---

# device42_ip (Resource)


## Exemple 

```hcl
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
```

## Schema

- `ip` (String - Optional) - Network of the subnet. Required for creation, cannot be modified after - subnet creation.
- `subnet` (String - Optional) - Subnet name of the IP.
- `subnet_id` (Int - Optional) - Subnet ID.
- `available` (String - Optional) - Whether the IP is available or not. Typically "yes" or "no".
- `vrf_group_id` (Int -Optional) - Subnet VRF Group ID.
- `device_id` (Int - Optional) - ID of the device to attach to the network.
- `suggest_ip` (Bool - Optional) - A boolean indicating if the Device42 IP suggestion functionality should be used to retrieve an available IP. If true and ip is not specified, the provider will automatically fetch a suggested IP.

### Read-Only

- `id` (String) The ID of this resource.


## Import

IP addresses can be imported using the Device42 IP id, e.g.

```bash
$ terraform import device42_ip.myIp <ip_id>
```