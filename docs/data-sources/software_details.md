---
page_title: "device42_software_details Data Source - terraform-provider-device42"
subcategory: "Business Services & Software"
description: |-
  Read Device42 software-in-use instances, optionally filtered by software or device.
---

# device42_software_details (Data Source)

The `device42_software_details` data source reads the **software instances detected on devices** (software "in use") in Device42. It is a read-only listing that can be filtered by `software_id` or `device_id`.

> **Note:** these records are auto-discovered by Device42 and are therefore exposed as a **data source only** — no resource exists to create them.

## Example Usage

### List software installed on a device

```hcl
data "device42_software_details" "on_device" {
  device_id = "1014"
}

output "installed" {
  value = [for i in data.device42_software_details.on_device.instances : i.software]
}
```

### List devices running a given software

```hcl
data "device42_software_details" "runs_apache" {
  software_id = data.device42_software.apache.software_id
}

output "devices" {
  value = [for i in data.device42_software_details.runs_apache.instances : i.device]
}
```

## Argument Reference

- `software_id` - (Optional) Filter instances to the given software component id.
- `device_id` - (Optional) Filter instances to the given device id.

## Attributes Reference

- `instances` - (Computed) The software-in-use instances matching the filters. Each element exposes:
  - `software_id` - ID of the software component.
  - `software` - Name of the software component.
  - `software_alias` - Alias of the software component.
  - `device_id` - ID of the device the software is installed on.
  - `device` - Name of the device the software is installed on.
  - `version` - Version of the software.
  - `arch` - Architecture of the software.
  - `user` - User associated with the instance.
  - `vendor` - Vendor of the software.
  - `appcomp` - Application component associated with the instance.
  - `install_date` - Installation date of the software.
  - `license_use_count` - License use count for the instance.
  - `count_in_licensing` - Whether the instance counts towards licensing.