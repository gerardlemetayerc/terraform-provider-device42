---
page_title: "device42_device Data Source - terraform-provider-device42"
subcategory: "asset management"
description: |-
  
---

# device42_device (Data Source)

Allow to query data about a device in Device42

## Exemple 

```hcl
datasource "device42_device" "searchVM" {
  name    = "myNewVMName"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The hostname of the device.
- `device_id`(Int) The ID of the device