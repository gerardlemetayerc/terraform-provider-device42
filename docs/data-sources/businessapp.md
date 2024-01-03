---
page_title: "device42_businessapp Data Source - terraform-provider-device42"
subcategory: "asset management"
description: |-
  This data source allows you to retrieve details about a specific Business Application in your Device42 instance using Terraform.
---

# device42_businessapp (Data Source)

The `device42_businessapp` data source permits the retrieval of details related to a Business Application stored in a Device42 instance. This is especially useful when you need to fetch information about a Business Application based on its name to use in other parts of your Terraform configuration.

## Example Usage 

```hcl
data "device42_businessapp" "example" {
  name = "ExampleBusinessApp"
}

output "businessapp_custom_fields" {
  value = data.device42_businessapp.example.custom_fields
}
```

## Argument Reference

The following arguments are supported:

- `name` (String - Mandatory) - The name of the business application you wish to retrieve.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `custom_fields` - Any custom fields associated with the Business Application in Device42.
- `id` - ID of the Business Application