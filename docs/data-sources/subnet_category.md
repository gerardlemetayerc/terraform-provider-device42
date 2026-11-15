---
page_title: "device42_subnet_category Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  This data source allows you to retrieve details about a specific subnet category in your Device42 instance using Terraform.
---

# device42_subnet_category (Data Source)

The `device42_subnet_category` data source allows you to fetch details of a specific subnet category from your Device42 instance.

## Example Usage

```hcl
data "device42_subnet_category" "example" {
  name = "Data"
}

output "subnet_category_id" {
  value = data.device42_subnet_category.example.id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the subnet category.
- `id` (Int - Optional) - Identifier of the subnet category.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for the subnet category in Device42.
- `description` - Description of the subnet category.
- `groups` - Comma separated list of groups associated with the subnet category.

Note: Only one of `name` or `id` should be provided to query a subnet category.