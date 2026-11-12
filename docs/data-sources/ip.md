---
page_title: "device42_ip Data Source - terraform-provider-device42"
subcategory: "network management"
description: |-
  Look up an existing Device42 IP by its address.
---

# device42_ip (Data Source)

The `device42_ip` data source retrieves details of an **existing** Device42 IP address. Unlike `device42_ip_suggested` (which finds a *free* IP), this data source looks up an IP that is already recorded in Device42 so you can reference its metadata (subnet, device, availability) elsewhere in your configuration.

## Example Usage

```hcl
data "device42_ip" "lookup" {
  ip = "10.1.10.5"
}

output "device_id" {
  value = data.device42_ip.lookup.device_id
}

output "subnet" {
  value = data.device42_ip.lookup.subnet
}
```

## Argument Reference

The following arguments are supported:

- `ip` - (Required) The IP address to look up (e.g. `10.1.10.5`).

## Attributes Reference

In addition to the argument above, the following attributes are exposed:

- `id` - The Device42 ID of the IP (used as the data source id).
- `ip_id` - The Device42 ID of the IP.
- `available` - Whether the IP is marked as available (`yes`/`no`).
- `device_id` - ID of the device attached to the IP, if any.
- `subnet_id` - Subnet ID the IP belongs to.
- `subnet` - Subnet network the IP belongs to (e.g. `10.1.0.0/17`).
- `vrf_group_name` - VRF group name of the subnet.
- `mac_address` - MAC address associated with the IP, if any.
- `label` - Label of the IP.
- `notes` - Notes for the IP.