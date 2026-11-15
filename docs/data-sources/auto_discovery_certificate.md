---
page_title: "device42_auto_discovery_certificate Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific certificate autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_certificate (Data Source)

The `device42_auto_discovery_certificate` data source allows you to fetch details of a specific certificate autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_certificate" "example" {
  name = "myCertJob"
}

output "job_id" {
  value = data.device42_auto_discovery_certificate.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the certificate autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `job_id` - The unique identifier for the job in Device42.
- `server` - IP or FQDN of the target server.
- `start_ip_address` - Start of the IP address range.
- `end_ip_address` - End of the IP address range.
- `exclude_servers` - Comma separated list of addresses to exclude.
- `ports` - Comma separated list of ports.
- `get_cipher_suites` - Find cipher suites supported by the server.
- `follow_chain` - Whether the certificate chain is followed.
- `debug_level` - Debug level of the job.
- `schedule_time` - Time of the schedule.
- `schedule_days` - Days of the week the job runs.
- `tags` - Comma separated list of tags.

Note: Only one of `name` or `job_id` should be provided to query a job.