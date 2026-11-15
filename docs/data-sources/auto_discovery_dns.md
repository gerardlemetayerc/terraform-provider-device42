---
page_title: "device42_auto_discovery_dns Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific DNS autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_dns (Data Source)

The `device42_auto_discovery_dns` data source allows you to fetch details of a specific DNS autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_dns" "example" {
  zonename = "example.local"
}

output "job_id" {
  value = data.device42_auto_discovery_dns.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `zonename` (String - Optional) - Zone name of the DNS autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `zonename` - Zone name of the DNS autodiscovery job.
- `job_id` - The unique identifier for the job in Device42.
- `nameserver` - IP or FQDN of the nameserver.
- `schedule_time` - Time of the schedule.
- `schedule_days` - Days of the week the job runs.

Note: Only one of `zonename` or `job_id` should be provided to query a job.