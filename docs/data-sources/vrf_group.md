---
page_title: "device42_vrf_group Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  This data source allows you to retrieve details about a specific VRF group in your Device42 instance using Terraform.
---

# device42_vrf_group (Data Source)

The `device42_vrf_group` data source allows you to fetch details of a specific VRF group from your Device42 instance.

## Example Usage

```hcl
data "device42_vrf_group" "example" {
  name = "ABCD Test VRF"
}

output "vrf_group_id" {
  value = data.device42_vrf_group.example.id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the VRF group.
- `id` (Int - Optional) - Identifier of the VRF group.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for the VRF group in Device42.
- `description` - Description of the VRF group.
- `tags` - Comma separated list of tags attached to the VRF group.
- `custom_fields` - Map of custom fields associated with the VRF group.

Note: Only one of `name` or `id` should be provided to query a VRF group.