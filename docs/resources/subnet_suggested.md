---
page_title: "device42_subnet_suggested Resource - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  Suggest and reserve the next available sub-subnet of a parent subnet.
---

# device42_subnet_suggested (Resource)

The `device42_subnet_suggested` resource asks Device42 to **suggest the next available sub-subnet** under a given parent subnet for a requested mask size, then **reserves it** by actually creating the subnet. Because the suggestion is consumed on creation, repeated applies stay stable (the same value is not re-offered in a loop). It is the subnet counterpart of [`device42_ip_suggested`](ip_suggested.md).

## Example Usage

```hcl
resource "device42_subnet_suggested" "new" {
  parent_subnet_id = data.device42_subnet.searchSubnet.subnet_id
  mask_bits        = 26
  name             = "new child subnet"
}

output "suggested_subnet_id" {
  value = device42_subnet_suggested.new.subnet_id
}
```

## Argument Reference

- `parent_subnet_id` - (Required, ForceNew) ID of the parent subnet to suggest a child from.
- `mask_bits` - (Required, ForceNew) Mask bits of the suggested sub-subnet.
- `name` - (Optional, ForceNew) Name to give the created sub-subnet.

## Attributes Reference

- `network` - (Computed) Network address of the suggested (and created) sub-subnet.
- `subnet_id` - (Computed) ID of the created sub-subnet.

## Import

This resource can be imported using the created sub-subnet id:

```hcl
terraform import device42_subnet_suggested.new 1234
```

## Note

The API returns an error when no free sub-subnet is available under the parent subnet (for example if it is fully allocated).