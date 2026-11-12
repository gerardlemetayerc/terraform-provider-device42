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
  id       = 123
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

- `label` - (Optional) The label of the password.

- `username` - (Optional) The username associated with the password.

- `id` - (Optional) The ID of the password.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `password` - The password of the account. This value is marked as sensitive and will not be displayed in plan or state output.

- `label` - The label of the password.

- `username` - The username associated with the password.
