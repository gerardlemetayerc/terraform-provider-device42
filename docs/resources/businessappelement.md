---
page_title: "device42_businessappelement Resource - terraform-provider-device42"
subcategory: "asset management"
description: |-
  This resource allows you to manage and associate devices with Business Applications in your Device42 instance using Terraform.
---

# device42_businessappelement (Resource)

The `device42_businessappelement` resource permits the management of devices associated with Business Applications in a Device42 instance. It can be used to assign, update, and delete device associations with Business Applications.

## Example Usage 

```hcl
resource "device42_businessappelement" "example" {
  businessapp_id = 123
  device_id      = 456
}
```


## Schema

The following arguments are supported:

- `businessapp_id` (Int - Required) - The ID of an existing Business Application to add elements (devices) to.
- `device_id` (Int - Required) - ID of an element (device) to add to the business app.

## Computed

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the device association with the Business Application in Device42.