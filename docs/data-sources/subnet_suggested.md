---
page_title: "device42_subnet_suggested Data Source - terraform-provider-device42"
subcategory: "network management"
description: |-
  Suggest the next available sub-subnet of a parent subnet.
---

# device42_subnet_suggested (Data Source)

The `device42_subnet_suggested` data source asks Device42 to **suggest the next available sub-subnet** under a given parent subnet, for a requested mask size. It is the subnet counterpart of [`device42_ip_suggested`](ip_suggested.md).

> **Note:** this is a **data source only** — a suggestion is read-only, so no `device42_subnet_suggested` resource exists.

## Example Usage

```hcl
data "device42_subnet_suggested" "new" {
  parent_subnet_id = data.device42_subnet.searchSubnet.subnet_id
  mask_bits        = 26
}

resource "device42_subnet" "child" {
  network   = data.device42_subnet_suggested.new.ip
  mask_bits = data.device42_subnet_suggested.new.mask_bits
  name      = "new child subnet"
}
```

## Argument Reference

- `parent_subnet_id` - (Required) ID of the parent subnet to suggest a child from.
- `mask_bits` - (Required) Mask bits of the suggested sub-subnet.

## Attributes Reference

- `ip` - (Computed) The suggested sub-subnet network address. Also used as the data source id.

## Note

The API returns an error when no free sub-subnet is available under the parent subnet (for example if it is fully allocated). In that case the data source returns an empty state without failing the plan — verify `ip` is non-empty before using it.