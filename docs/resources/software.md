---
page_title: "device42_software Resource - terraform-provider-device42"
subcategory: "asset management"
description: |-
  Create, update and delete a Device42 software component.
---

# device42_software (Resource)

The `device42_software` resource manages a **software component** (catalog entry) in Device42. It supports full create/read/update/delete, custom fields and import.

> **Note:** the Device42 API (`/api/1.0/software/`) upserts a software component by its **`name`** — re-posting the same `name` updates the existing component in place. Because `name` is the upsert key, it is `ForceNew` (renaming recreates the component).

## Example Usage

```hcl
resource "device42_software" "apache" {
  name          = "Apache HTTP Server"
  vendor        = "Apache Software Foundation"
  notes         = "Web server"
  software_type = "Server"
  category      = "Web"
  aliases       = "httpd"
  tags          = "prod,web"

  custom_fields = {
    owner = "platform-team"
  }
}
```

## Argument Reference

- `name` - (Required, Forces new resource) Name of the software component. Used as the identifier by the Device42 API.
- `software_id` - (Computed) Software component ID assigned by Device42.
- `vendor` - (Optional) Vendor of the software component.
- `description` - (Optional) Description of the software component.
- `notes` - (Optional) Any additional notes.
- `software_type` - (Optional) Type of the software component.
- `deployment_type` - (Optional) Deployment type of the software component.
- `category` - (Optional) Category of the software component.
- `aliases` - (Optional) Aliases of the software component.
- `licensing_model` - (Optional) Licensing model of the software component.
- `track_licensed_count_by_keys` - (Optional) Whether licensed counts are tracked by license keys (`yes`/`no`).
- `tags` - (Optional) Comma-separated tags to attach to the software component.
- `custom_fields` - (Optional) Map of custom fields to set on the software component.

## Import

Software components can be imported using the Device42 software `id`:

```shell
$ terraform import device42_software.apache <software_id>
```

The import key is the numeric software `id` (e.g. `5076`).