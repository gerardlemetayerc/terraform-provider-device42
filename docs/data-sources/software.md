---
page_title: "device42_software Data Source - terraform-provider-device42"
subcategory: "asset management"
description: |-
  Read Device42 software components, either by name or as a full list.
---

# device42_software (Data Source)

The `device42_software` data source reads software components (catalog entries) defined in Device42.

It supports two modes:

- **Lookup by name** — set `name` to retrieve the `software_id` of a single component.
- **List all** — omit `name` to get the complete list of components as the `software` list attribute.

## Example Usage

### Lookup a single software component by name

```hcl
data "device42_software" "apache" {
  name = "Apache HTTP Server"
}

output "apache_id" {
  value = data.device42_software.apache.software_id
}
```

### List all software components

```hcl
data "device42_software" "all" {}

output "software_names" {
  value = [for s in data.device42_software.all.software : s.name]
}
```

## Argument Reference

- `name` - (Optional) The exact name of a software component to look up. When omitted, the data source returns all components.

## Attributes Reference

- `software_id` - (Optional/Computed) The Device42 ID of the component matched by `name`. Only set in lookup mode.
- `software` - (Computed) The complete list of software components. Each element exposes:
  - `name` - Name of the software component.
  - `software_id` - Device42 ID of the software component.