---
page_title: "device42_vendor Data Source - terraform-provider-device42"
subcategory: "asset management"
description: |-
  Read Device42 vendors, either by name or as a full list.
---

# device42_vendor (Data Source)

The `device42_vendor` data source reads vendors defined in Device42.

It supports two modes:

- **Lookup by name** — set `name` to retrieve the `vendor_id` of a single vendor.
- **List all** — omit `name` to get the complete list of vendors as the `vendors` list attribute.

## Example Usage

### Lookup a single vendor by name

```hcl
data "device42_vendor" "acme" {
  name = "Acme Corp"
}

output "acme_id" {
  value = data.device42_vendor.acme.vendor_id
}
```

### List all vendors

```hcl
data "device42_vendor" "all" {}

output "vendor_names" {
  value = [for v in data.device42_vendor.all.vendors : v.name]
}
```

## Argument Reference

- `name` - (Optional) The exact name of a vendor to look up. When omitted, the data source returns all vendors.

## Attributes Reference

- `vendor_id` - (Optional/Computed) The Device42 ID of the vendor matched by `name`. Only set in lookup mode.
- `vendors` - (Computed) The complete list of vendors. Each element exposes:
  - `name` - Name of the vendor.
  - `vendor_id` - Device42 ID of the vendor.