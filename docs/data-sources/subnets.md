---
page_title: "device42_subnets Data Source - terraform-provider-device42"
subcategory: "network management"
description: |-
  Read multiple Device42 subnets, optionally filtered by VLAN, name or VRF group.
---

# device42_subnets (Data Source)

The `device42_subnets` data source reads **multiple** Device42 subnets at once. It is the **plural** counterpart of [`device42_subnet`](subnet.md) and is used when you need to enumerate several subnets — for example all subnets attached to a given VLAN.

> **Use `device42_subnet` (singular)** when you are looking for a **single** subnet by its `name`/`subnet_id`. **Use `device42_subnets` (plural)** when you want to **list/iterate** subnets matching a filter (e.g. all subnets of a VLAN).

## Relationship to VLANs

A subnet is attached to a VLAN through its **`parent_vlan_id`** field, which equals the `id` of the parent VLAN. The Device42 API does **not** filter subnets by VLAN server-side, so this data source fetches all subnets and filters them client-side — meaning the `parent_vlan_id` filter always returns subnets whose `parent_vlan_id` matches exactly.

## Example Usage

### All subnets attached to a VLAN

```hcl
data "device42_vlan" "prod" {
  name = "CRMConsultant"
}

data "device42_subnets" "in_vlan" {
  parent_vlan_id = data.device42_vlan.prod.id
}

output "subnet_networks" {
  value = [for s in data.device42_subnets.in_vlan.subnets : s.network]
}
```

### All subnets of a VRF group

```hcl
data "device42_subnets" "all" {}

output "subnets" {
  value = data.device42_subnets.all.subnets
}
```

## Argument Reference

- `parent_vlan_id` - (Optional) Filter subnets attached to the given parent VLAN id.
- `name` - (Optional) Filter subnets by exact name.
- `vrf_group_id` - (Optional) Filter subnets by VRF group id.
- `vrf_group_name` - (Optional) Filter subnets by VRF group name.

At least one filter may be omitted to return all subnets.

## Attributes Reference

- `subnets` - (Computed) The subnets matching the filters. Each element exposes:
  - `subnet_id` - Device42 subnet id.
  - `name` - Name of the subnet.
  - `network` - Network address of the subnet.
  - `mask_bits` - Mask bits (CIDR) of the subnet.
  - `gateway` - Gateway of the subnet.
  - `range_begin` - First usable IP of the subnet.
  - `range_end` - Last usable IP of the subnet.
  - `vrf_group_id` - VRF group id of the subnet.
  - `vrf_group_name` - VRF group name of the subnet.
  - `parent_vlan_id` - ID of the parent VLAN the subnet is attached to.
  - `parent_subnet_id` - ID of the parent subnet, if any.
  - `description` - Description of the subnet.
  - `allocated` - Allocation status of the subnet.
  - `customer` - Customer attached to the subnet.
  - `custom_fields` - Custom fields of the subnet.