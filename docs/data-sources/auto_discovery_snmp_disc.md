---
page_title: "device42_auto_discovery_snmp_disc Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific SNMP autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_snmp_disc (Data Source)

The `device42_auto_discovery_snmp_disc` data source allows you to fetch details of a specific SNMP autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_snmp_disc" "example" {
  name = "mySnmpJob"
}

output "job_id" {
  value = data.device42_auto_discovery_snmp_disc.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the SNMP autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `job_id` - The unique identifier for the job in Device42.
- `server` - IP or FQDN of the target server.
- `start_ip` - Start of the IP address range to scan.
- `end_ip` - End of the IP address range to scan.
- `exclude_servers` - Comma separated list of addresses to exclude.
- `snmp_retries` - Number of SNMP retries.
- `snmp_timeout` - SNMP timeout in seconds.
- `snmp_min_parallel` - Minimum number of parallel SNMP requests.
- `do_port_check` - Whether a port check is performed on discovered devices.
- `clear_customer` - Whether the customer assignment is cleared on matched devices.
- `snmp_v3` - Whether SNMPv3 is used.
- `schedule_time` - Time of the schedule.
- `schedule_days` - Days of the week the job runs.
- `tags` - Comma separated list of tags.

Note: Only one of `name` or `job_id` should be provided to query a job.