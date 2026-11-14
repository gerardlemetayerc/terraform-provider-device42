---
page_title: "device42_auto_discovery_ipmi Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific ipmi autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_ipmi (Data Source)

The `device42_auto_discovery_ipmi` data source allows you to fetch details of a specific ipmi autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_ipmi" "example" {
  name = "ipmi"
}

output "job_id" {
  value = data.device42_auto_discovery_ipmi.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the ipmi autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `job_id` - The unique identifier for the job in Device42.
- `name` - Name of the job.
- `server` - IP or FQDN of the target server.
- `schedule_time` - Time of the schedule.
- `object_category` - Object category assigned to discovered devices.
- `update_model_if_found` - Whether the model of devices that are found is updated.
- `overwrite_object_categories` - Whether the object categories of existing devices are overwritten.
- `ip_start` - Start of the IP address range.
- `ip_end` - End of the IP address range.
- `overwrite_device_name` - Whether the name of existing devices is overwritten.
- `hostname_to_use` - Hostname to use for the discovered devices.
- `debug_level` - Debug level of the job.
- `discovery_type` - Type of discovery to run.
- `run_as_operator` - Whether discovery is run as operator.
- `schedule_days` - Days of the week the job runs.
- `exclude_servers` - Comma separated list of addresses to exclude.
- `tags` - Comma separated list of tags.

Note: Only one of `name` or `job_id` should be provided to query a job.