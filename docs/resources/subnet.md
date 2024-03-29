---
page_title: "device42_subnet Resource - terraform-provider-device42"
subcategory: "network management"
description: |-
---

# device42_subnet (Resource)


## Exemple 

```hcl
resource "device42_subnet" "myNewNetwork" {
    name        = "myNewNetwork"
    mask_bits   = 24
    network     = "192.168.100.0"
    vrf_group   = "Chicago Data Center"
}
```



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network` (String) Network of the subnet. Required for creation, cannot be modified after subnet creation.
- `mask_bits` (String) Mask bits of the subnet. Required for creation, can be modified after subnet creation.

### Optional

- `vrf_group` (String) Subnet VRF Group
- `custom_fields` (Map of String) Any custom fields that will be used in device42.
- `parent_vlan_id` (Int) Parent vlan ID of the subnet.


### Read-Only

- `id` (String) The ID of this resource.


