---
page_title: "device42_vlans Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  Read multiple Device42 VLANs, optionally filtered by name or id.
---

# device42_vlans (Data Source)

The `device42_vlans` data source reads **multiple** Device42 VLANs at once. It is the **plural** counterpart of [`device42_vlan`](vlan.md), used when you need to enumerate or iterate VLANs.

## Example Usage

### List all VLANs

```hcl
data "device42_vlans" "all" {}

output "vlan_names" {
  value = [for v in data.device42_vlans.all.vlans : v.name]
}
```

### Look up a VLAN by name

```hcl
data "device42_vlans" "prod" {
  name = "CRMConsultant"
}

output "vlan_id" {
  value = data.device42_vlans.prod.vlans[0].vlan_id
}
```

## Argument Reference

- `name` - (Optional) Filter VLANs by exact name.
- `vlan_id` - (Optional) Filter VLANs by id.

When both are omitted, all VLANs are returned.

## Attributes Reference

- `vlans` - (Computed) The VLANs matching the filters. Each element exposes:
  - `vlan_id` - Device42 VLAN id.
  - `number` - VLAN number.
  - `name` - Name of the VLAN.
  - `description` - Description of the VLAN.
  - `tags` - Tags of the VLAN.