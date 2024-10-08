---
page_title: "device42_device Resource - terraform-provider-device42"
subcategory: "asset management"
description: |-
  
---

# device42_device (Resource)

## Exemple 

```hcl
resource "device42_device" "myNewVM" {
  name          = "myNewVMName"
  type          = "virtual"
  service_level = "Production"
  custom_fields = {
    field_name  = "field_value"
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The hostname of the device.

### Optional

- `custom_fields` (Map of String) Any custom fields that will be used in device42.
- `service_level` (String) Service Level of the device (default is d42null).
- `type` (String) The type of the device. Valid values are 'physical', 'virtual', 'unknown', 'cluster' (default is virtual)
- `archive_on_destroy (bool)` Specify if ressource need to be archived on destroy call **(default value: false)**

### Read-Only

- `id` (String) The ID of this resource.

### Import

Device object support ressource importation by device name.


Exemple :

```hcl
terraform import module.vm.device42_device.host myVMName
```