---
page_title: "device42_vlan Data Source - terraform-provider-device42"
subcategory: "network management"
description: |-
  This data source allows you to retrieve details about a specific VLAN in your Device42 instance using Terraform.
---

# device42_vlan (Data Source)

The `device42_vlan` data source allows you to fetch details of a specific VLAN from your Device42 instance. This can be useful when you need to access VLAN information, such as its ID, number, and description, for use in other parts of your Terraform configurations.

## Example Usage 

```hcl
data "device42_vlan" "example" {
  vlan_id = 1234
}

output "vlan_number" {
  value = data.device42_vlan.example.number
}

```

## Argument Reference
The following arguments are supported:

- `vlan_id` (Int - Optional) - The unique identifier for the VLAN.
- `number` (Int - Optional) - The number associated with the VLAN.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `name` - The name of the VLAN.
- `description` - A description of the VLAN.

Note: Only one of vlan_id or number should be provided to query a VLAN. If both are provided, vlan_id will take precedence.
