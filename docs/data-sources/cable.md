---
page_title: "device42_cable Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  This data source allows you to retrieve details about a specific cable in your Device42 instance using Terraform.
---

# device42_cable (Data Source)

The `device42_cable` data source allows you to fetch details of a specific cable from your Device42 instance.

## Example Usage

```hcl
data "device42_cable" "example" {
  cable_id = "CableEnd"
}

output "cable_details" {
  value = data.device42_cable.example
}

```

## Argument Reference
The following arguments are supported:

- `cable_id` (String - Optional) - Identifier of the cable.
- `id` (Int - Optional) - Identifier of the cable.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for the cable in Device42.
- `cable_length` - Length of the cable.
- `cable_length_units` - Unit of the cable length.
- `notes` - Any additional notes for the cable.
- `vendor` - Vendor of the cable.
- `tags` - Comma separated list of tags attached to the cable.
- `custom_fields` - Map of custom fields associated with the cable.

Note: Only one of `cable_id` or `id` should be provided to query a cable.