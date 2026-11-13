---
page_title: "device42_software_license_key Resource - terraform-provider-device42"
subcategory: "Business Services & Software"
description: |-
  Create, update and delete a Device42 software license key.
---

# device42_software_license_key (Resource)

The `device42_software_license_key` resource manages a **software license key** in Device42, linked to a software component via `software_id`. It supports full create/read/update/delete and import.

## Example Usage

```hcl
resource "device42_software" "apache" {
  name   = "Apache HTTP Server"
  vendor = "Apache Software Foundation"
}

resource "device42_software_license_key" "apache" {
  key           = "LPKG-1234-ABCD-5678"
  software_id   = device42_software.apache.software_id
  license_count = 10
  notes         = "Enterprise subscription"
}
```

## Argument Reference

- `key` - (Required, Forces new resource) The license key value. Used as the identifier by the Device42 API.
- `software_id` - (Required) ID of the software component this license key belongs to.
- `software_name` - (Computed) Name of the software component (derived by the API).
- `license_count` - (Optional) Number of licenses covered by this key.
- `notes` - (Optional) Any additional notes.
- `license_key_id` - (Computed) License key ID assigned by Device42.

## Import

Software license keys can be imported using the Device42 license key `id`:

```shell
$ terraform import device42_software_license_key.apache <license_key_id>
```

The import key is the numeric license key `id` (e.g. `8`).