---
page_title: "device42_pdu_model Data Source - terraform-provider-device42"
subcategory: "Power Management"
description: |-
  This data source allows you to retrieve details about a specific PDU model in your Device42 instance using Terraform.
---

# device42_pdu_model (Data Source)

The `device42_pdu_model` data source allows you to fetch details of a specific PDU model from your Device42 instance.

## Example Usage

```hcl
data "device42_pdu_model" "example" {
  name = "AP7800"
}

output "pdu_model_id" {
  value = data.device42_pdu_model.example.pdu_model_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the PDU model.
- `pdu_model_id` (Int - Optional) - Identifier of the PDU model.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `pdu_model_id` - The unique identifier for the PDU model in Device42.
- `manufacturer` - Manufacturer of the PDU model.
- `description` - Description of the PDU model.
- `notes` - Any additional notes for the PDU model.
- `depth` - Depth of the PDU model.
- `size` - Size of the PDU model.
- `width` - Width of the PDU model.
- `height` - Height of the PDU model.

Note: Only one of `name` or `pdu_model_id` should be provided to query a PDU model.