---
page_title: "device42_service_level Data Source - terraform-provider-device42"
subcategory: "Asset Management"
description: |-
  Read Device42 service levels (tiers), either by name or as a full list.
---

# device42_service_level (Data Source)

The `device42_service_level` data source reads the **service levels** (availability tiers) defined in Device42, such as `Production`, `QA`, `Development`, `Decommissioned`.

It supports two modes:

- **Lookup by name** — set `name` to retrieve the `service_level_id` of a single tier.
- **List all** — omit `name` to get the complete list of tiers as the `levels` list attribute.

> **Note:** because the Device42 API (`/1.0/service_level/`) does not support deleting a service level, this is provided as a **data source only** — no `device42_service_level` resource exists.

## Example Usage

### Lookup a single service level by name

```hcl
data "device42_service_level" "prod" {
  name = "Production"
}

output "production_id" {
  value = data.device42_service_level.prod.service_level_id
}
```

### List all service levels

```hcl
data "device42_service_level" "all" {}

output "service_level_names" {
  value = [for l in data.device42_service_level.all.levels : l.name]
}
```

## Argument Reference

- `name` - (Optional) The exact name of a service level to look up (e.g. `Production`). When omitted, the data source returns all service levels.

## Attributes Reference

- `service_level_id` - (Optional/Computed) The Device42 ID of the service level matched by `name`. Only set in lookup mode.
- `levels` - (Computed) The complete list of service levels. Each element exposes:
  - `name` - Name of the service level.
  - `service_level_id` - Device42 ID of the service level.