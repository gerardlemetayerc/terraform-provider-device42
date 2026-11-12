---
page_title: "device42_vendor Resource - terraform-provider-device42"
subcategory: "asset management"
description: |-
  Create, update and delete a Device42 vendor.
---

# device42_vendor (Resource)

The `device42_vendor` resource manages a **vendor** (supplier/manufacturer) in Device42. It supports full create/read/update/delete and import.

> **Note:** the Device42 API (`/api/1.0/vendors/`) upserts a vendor by its **`name`** — re-posting the same `name` updates the existing vendor in place. Because `name` is the upsert key, it is `ForceNew` (renaming a vendor recreates it).

## Example Usage

```hcl
resource "device42_vendor" "acme" {
  name         = "Acme Corp"
  notes        = "Hardware supplier"
  home_page    = "https://acme.example.com"
  contact_info = "support@acme.example.com"
  account_no   = "ACC-12345"
  aliases      = "Acme"

  custom_fields = {
    region = "eu"
  }
}
```

## Argument Reference

- `name` - (Required, Forces new resource) Name of the vendor. Used as the identifier by the Device42 API.
- `vendor_id` - (Computed) Vendor ID assigned by Device42.
- `home_page` - (Optional) Home page URL of the vendor.
- `contact_info` - (Optional) Contact information for the vendor.
- `account_no` - (Optional) Account number associated with the vendor.
- `aliases` - (Optional) Aliases of the vendor.
- `escalation_1` - (Optional) First escalation contact.
- `escalation_2` - (Optional) Second escalation contact.
- `notes` - (Optional) Any additional notes.
- `tags` - (Optional) Comma-separated tags to attach to the vendor.
- `custom_fields` - (Optional) Map of custom fields to set on the vendor.

## Import

Vendors can be imported using the Device42 `vendor_id`:

```shell
$ terraform import device42_vendor.acme <vendor_id>
```

The import key is the numeric `vendor_id` (e.g. `238`).