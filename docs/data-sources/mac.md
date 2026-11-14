---
page_title: "device42_mac Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  This data source allows you to retrieve details about a specific MAC address in your Device42 instance using Terraform.
---

# device42_mac (Data Source)

The `device42_mac` data source allows you to fetch details of a specific MAC address from your Device42 instance.

## Example Usage

```hcl
data "device42_mac" "example" {
  macaddress = "00:11:22:33:44:55"
}

output "mac_id" {
  value = data.device42_mac.example.macaddress_id
}

```

## Argument Reference
The following arguments are supported:

- `macaddress` (String - Optional) - MAC address value.
- `macaddress_id` (Int - Optional) - Identifier of the MAC address.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `macaddress_id` - The unique identifier for the MAC address in Device42.
- `vlan_id` - Identifier of the VLAN this MAC address belongs to.
- `port` - Port associated with the MAC address.

Note: Only one of `macaddress` or `macaddress_id` should be provided to query a MAC address.