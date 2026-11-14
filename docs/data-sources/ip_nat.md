---
page_title: "device42_ip_nat Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  This data source allows you to retrieve details about a specific IP NAT rule in your Device42 instance using Terraform.
---

# device42_ip_nat (Data Source)

The `device42_ip_nat` data source allows you to fetch details of a specific IP NAT rule from your Device42 instance.

## Example Usage

```hcl
data "device42_ip_nat" "example" {
  name = "foo"
}

output "nat_id" {
  value = data.device42_ip_nat.example.id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the IP NAT rule.
- `id` (Int - Optional) - Identifier of the IP NAT rule.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for the IP NAT rule in Device42.
- `ip_address_from` - Source IP address of the NAT rule.
- `ip_address_to` - Destination IP address of the NAT rule.
- `notes` - Any additional notes for the NAT rule.
- `protocol` - Protocol of the NAT rule.
- `two_way_relation` - Whether the NAT is two-way.
- `source_port_start` - Start of the source port range.
- `source_port_end` - End of the source port range.
- `target_port_start` - Start of the target port range.
- `target_port_end` - End of the target port range.
- `ip_address_from_end` - End of the source IP address range.
- `ip_address_to_end` - End of the destination IP address range.

Note: Only one of `name` or `id` should be provided to query an IP NAT rule.