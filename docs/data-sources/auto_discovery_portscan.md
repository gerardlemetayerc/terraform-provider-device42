---
page_title: "device42_auto_discovery_portscan Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific portscan autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_portscan (Data Source)

The `device42_auto_discovery_portscan` data source allows you to fetch details of a specific portscan autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_portscan" "example" {
  name = "portscan"
}

output "job_id" {
  value = data.device42_auto_discovery_portscan.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the portscan autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `job_id` - The unique identifier for the job in Device42.
- `name` - Name of the job.
- `server` - IP or FQDN of the target server.
- `start_ip` - Start of the IP address range.
- `end_ip` - End of the IP address range.
- `capture_hosts_file` - Whether the hosts file is captured.
- `discover_cluster` - Whether clusters are discovered.
- `overwrite_device_categories` - Whether the device categories of existing devices are overwritten.
- `schedule_time` - Time of the schedule.
- `schedule_days` - Days of the week the job runs.
- `tags` - Comma separated list of tags.

Note: Only one of `name` or `job_id` should be provided to query a job.