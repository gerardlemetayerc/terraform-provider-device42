---
page_title: "device42_auto_discovery_pingsweep Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific pingsweep autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_pingsweep (Data Source)

The `device42_auto_discovery_pingsweep` data source allows you to fetch details of a specific pingsweep autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_pingsweep" "example" {
  name = "myPingSweepJob"
}

output "job_id" {
  value = data.device42_auto_discovery_pingsweep.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the pingsweep autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `job_id` - The unique identifier for the job in Device42.
- `name` - Name of the pingsweep autodiscovery job.
- `networks` - Comma separated list of networks to sweep.
- `schedule_time` - Time of the schedule.
- `subnet_category` - Subnet category assigned to discovered subnets.
- `create_new_subnet` - Whether a new subnet is created for the discovered addresses.
- `debug_level` - Debug level of the job.
- `overwrite_subnet_categories` - Whether the subnet categories are overwritten.
- `schedule_days` - Days of the week the job runs.
- `reverse_dns` - Whether a reverse DNS lookup is performed on discovered addresses.
- `strip_domain_suffix` - Whether the domain suffix is stripped from discovered hostnames.
- `type` - Type of the pingsweep job.

Note: Only one of `name` or `job_id` should be provided to query a job.