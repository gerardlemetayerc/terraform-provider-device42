---
page_title: "device42_pdu Data Source - terraform-provider-device42"
subcategory: "Power Management"
description: |-
  This data source allows you to retrieve details about a specific PDU in your Device42 instance using Terraform.
---

# device42_pdu (Data Source)

The `device42_pdu` data source allows you to fetch details of a specific physical PDU from your Device42 instance.

## Example Usage

```hcl
data "device42_pdu" "example" {
  name = "PDU-01"
}

output "pdu_id" {
  value = data.device42_pdu.example.pdu_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the PDU.
- `pdu_id` (Int - Optional) - Identifier of the PDU.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `pdu_id` - The unique identifier for the PDU in Device42.
- `pdu_serial_no` - Serial number of the PDU.
- `pdu_model` - Name of the PDU model this PDU is based on.
- `pdu_model_id` - Identifier of the PDU model.
- `notes` - Any additional notes for the PDU.
- `gateway_pdu_id` - Identifier of the parent/gateway PDU.
- `custom_fields` - Map of custom fields associated with the PDU.

Note: Only one of `name` or `pdu_id` should be provided to query a PDU.