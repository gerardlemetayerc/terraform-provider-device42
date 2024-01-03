---
page_title: "device42_password Data Source - terraform-provider-device42"
subcategory: "Identity and Access Management"
description: |-
  This data source allows you to retrieve details about a specific password stored in Device42.

---

# device42_password (Data Source)

The `device42_password` data source allows you to retrieve details about a specific password stored in Device42. This can be used to fetch password details based on certain criteria like username, device, or application component.

## Example Usage 

```hcl
data "device42_password" "example" {
  username = "exampleUsername"
  category = "exampleCategory"
}

output "password_details" {
  value = {
    password = data.device42_password.example.password
    label    = data.device42_password.example.label
  }
}
```

## Argument Reference

The following arguments are supported:

- `category` - (Optional) The category of the password.

- `label` - (Optional) The label of the password.

- `username` - (Optional) The username associated with the password.

- `device` - (Optional) The device associated with the password.

- `appcomp` - (Optional) The application component associated with the password.

- `id` - (Optional) The ID of the password.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `password` - The password.

- `device` - The associated device.

- `category` - The category of the password.

- `label` - The label of the password.

- `username` - The username associated with the password.
