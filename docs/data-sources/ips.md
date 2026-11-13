---
page_title: "device42_ips Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  Read Device42 IPs of a subnet (server-side filtered).
---

# device42_ips (Data Source)

The `device42_ips` data source reads **multiple** Device42 IP addresses, typically all the IPs of a given subnet. It is the **plural** counterpart of [`device42_ip`](ip.md).

> **Note:** the Device42 API (`/2.0/ips/`) is **paginated** and returns at most 1000 records per response. To reliably list the IPs of a subnet, always set `subnet_id` — that filter is applied **server-side**. Without `subnet_id`, only the first 1000 IPs are returned.

## Example Usage

### All IPs of a subnet

```hcl
data "device42_subnets" "in_vlan" {
  parent_vlan_id = data.device42_vlan.prod.id
}

data "device42_ips" "subnet_ips" {
  subnet_id = data.device42_subnet.searchSubnet.subnet_id
}

output "ip_addresses" {
  value = [for i in data.device42_ips.subnet_ips.ips : i.ip]
}
```

## Argument Reference

- `subnet_id` - (Optional) Filter IPs belonging to the given subnet id. This filter is applied server-side and is the reliable way to list a subnet's IPs.

## Attributes Reference

- `ips` - (Computed) The IPs matching the filter. Each element exposes:
  - `ip_id` - Device42 ID of the IP.
  - `ip` - The IP address.
  - `available` - Whether the IP is marked as available.
  - `subnet_id` - Subnet ID the IP belongs to.
  - `subnet` - Subnet network the IP belongs to.
  - `vrf_group_name` - VRF group name of the subnet.
  - `mac_address` - MAC address associated with the IP.
  - `label` - Label of the IP.
  - `notes` - Notes for the IP.
  - `device_id` - ID of the device attached to the IP, if any.