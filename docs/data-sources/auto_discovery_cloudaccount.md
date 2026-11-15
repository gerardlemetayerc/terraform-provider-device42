---
page_title: "device42_auto_discovery_cloudaccount Data Source - terraform-provider-device42"
subcategory: "Auto Discovery"
description: |-
  This data source allows you to retrieve details about a specific cloudaccount autodiscovery job in your Device42 instance using Terraform.
---

# device42_auto_discovery_cloudaccount (Data Source)

The `device42_auto_discovery_cloudaccount` data source allows you to fetch details of a specific cloudaccount autodiscovery job from your Device42 instance.

## Example Usage

```hcl
data "device42_auto_discovery_cloudaccount" "example" {
  name = "cloudaccount"
}

output "job_id" {
  value = data.device42_auto_discovery_cloudaccount.example.job_id
}

```

## Argument Reference
The following arguments are supported:

- `name` (String - Optional) - Name of the cloudaccount autodiscovery job.
- `job_id` (Int - Optional) - Identifier of the job.

## Attribute Reference
In addition to all arguments above, the following attributes are exported:

- `job_id` - The unique identifier for the job in Device42.
- `name` - Name of the job.
- `type` - The cloud type.
- `server` - IP or FQDN of the cloud account server.
- `ip` - IP of the cloud account server.
- `username` - Username to authenticate with.
- `tenant` - Tenant of the cloud account.
- `vendor` - Vendor of the cloud account.
- `object_category` - Object category assigned to discovered devices.
- `vrf_group` - VRF group assigned to discovered devices.
- `notes` - Notes attached to the job.
- `service_level` - Service level assigned to discovered devices.
- `regions` - Comma separated list of regions.
- `aws_roles` - Comma separated list of AWS roles.
- `oracle_regions` - Comma separated list of Oracle cloud regions.
- `google_regions` - Comma separated list of Google Cloud regions.
- `alibaba_regions` - Comma separated list of Alibaba cloud regions.
- `outscale_regions` - Comma separated list of Outscale regions.
- `collect_db_info` - Whether database information is collected.
- `add_suffix` - Whether a suffix is added to discovered names.
- `strip_domain_name` - Whether the domain name is stripped from discovered device names.
- `match_name_only_for_virtuals` - Whether names are matched only for virtual devices.
- `match_name_only_for_hypervisor` - Whether names are matched only for hypervisor devices.
- `overwrite_object_category` - Whether the object category of existing devices is overwritten.
- `remove_unfound_instances` - Whether instances that are no longer found are removed.
- `add_tags_as_custom_fields` - Whether tags are added as custom fields.
- `schedule_time` - Time of the schedule.
- `schedule_days` - Days of the week the job runs.
- `tags` - Comma separated list of tags.

Note: Only one of `name` or `job_id` should be provided to query a job.