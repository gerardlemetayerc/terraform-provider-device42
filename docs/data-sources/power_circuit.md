---
page_title: "device42_power_circuit Data Source - terraform-provider-device42"
subcategory: "Power Management"
description: |-
  This data source allows you to retrieve details about a specific power circuit in your Device42 instance using Terraform.
---

# device42_power_circuit (Data Source)

The `device42_power_circuit` data source allows you to fetch details of a specific power circuit attached to a breaker panel from your Device42 instance.

## Example Usage

```hcl
data "device42_power_circuit" "example" {
  breakerpanel_id = "12"
  number          = "4"
}

output "power_circuit_label" {
  value = data.device42_power_circuit.example.label
}

```

## Argument Reference
The following arguments are supported:

- `power_circuit_id` (Int - Optional) - Identifier of the power circuit.
- `breakerpanel_id` (String - Optional) - Identifier of the breaker panel.
- `number` (String - Optional) - Number of the circuit within the breaker panel.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `power_circuit_id` - The unique identifier for the power circuit in Device42.
- `breakerpanel` - Name of the breaker panel.
- `label` - Label of the power circuit.
- `amps` - Amperage rating of the circuit.
- `voltage` - Voltage rating of the circuit.
- `connection_name` - Name of the power connection type.

Note: Provide either `power_circuit_id`, or both `breakerpanel_id` and `number`, to query a power circuit.