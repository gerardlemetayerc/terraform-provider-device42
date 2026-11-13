# Release v1.4.1 — Test Report

**Branch:** `v1.4.1`
**Provider source:** `gerardlemetayerc/device42`
**Test date:** 2026-11-13

---

## 1. Unit tests — `go test ./...`

All tests pass.

| Test | Coverage |
|---|---|
| `TestProviderInternalValidate` | Validates the whole provider schema (config, all resources + data sources) via the SDK internal validator |
| `TestMsgID` | Extracts resource ID from Device42 `msg` array (numeric / `json.Number` / int / malformed) |
| `TestAPIDeleteResponseUnmarshal` | Delete response decoding |
| `TestAPIStringOrNumberUnmarshal` | Tolerates `code`/`id` encoded as number or string |
| `TestAPIStringOrNumberUnmarshalInvalid` | Rejects malformed values |
| `TestFilterCertificateByName` | Certificate lookup by name |
| `TestFindContact` / `TestFindContactByName` | Customer-contact lookup by id / by name (incl. ambiguous) |
| `TestFilterCustomerByName` | Customer lookup by name |
| `TestYesNo` | Boolean → "yes"/"no" conversion |
| `TestFlattenCustomFields` | Flattens API custom_fields map |
| `TestSuppressCustomFieldsDiffs` | Custom-field diff suppression |
| `TestAPICodeUnmarshal` | `code` field tolerates number / quoted string |
| `TestArchiveResponseDecodesStringCode` / `DecodesBareArray` | Archive response shape tolerance |
| `TestIsIPAddress` | IP address validation |
| `TestRoomLookupResponseUnmarshal` / `TestRackLookupResponseUnmarshal` | Room / rack lookup response decoding |

---

## 2. Functional tests — `terraform` (validation, plan, apply, destroy)

### 2.1 Data sources (read)

All **21** data-source reads succeed.

| Data source | Config used |
|---|---|
| `device42_building` | by `name` and by `building_id` |
| `device42_room` | by `name` and by `room_id` |
| `device42_customer` | by `name` |
| `device42_subnet` | by `subnet_id` and by `network`+`mask_bits`+`vrf_group_name` |
| `device42_ip` | by `ip` |
| `device42_ips` | by `subnet_id` |
| `device42_vlan` | by `vlan_id` |
| `device42_vlans` | empty (all) |
| `device42_subnets` | by `parent_vlan_id` |
| `device42_service_level` | by `name` and empty (all) |
| `device42_software` | by `name` and empty (all) |
| `device42_software_details` | by `device_id` |
| `device42_vendor` | by `name` and empty (all) |
| `device42_device` | by `name` |

### 2.2 Resources (create / destroy)

All **14** resources create and destroy successfully.

| Resource | Config used (key fields) |
|---|---|
| `device42_building` | `name`, `address`, `contact_name`, `contact_phone`, `notes`, `tags` |
| `device42_room` | `name`, `building_id`, `notes` |
| `device42_customer` | `name`, `notes`, `manager` |
| `device42_customer_contact` | `customer`, `name`, `type`, `email` |
| `device42_vendor` | `name`, `notes`, `contact_info` |
| `device42_software` | `name`, `notes`, `software_type` (`managed`) |
| `device42_software_license_key` | `key`, `software_id`, `license_count`, `notes` |
| `device42_subnet` | `name`, `network`, `mask_bits` |
| `device42_subnet_suggested` | `parent_subnet_id`, `mask_bits`, `name` |
| `device42_vlan` | `name`, `number`, `description` |
| `device42_device` | `name`, `building`, `customer`, `archive_on_destroy` |
| `device42_ip` | `subnet_id`, `ip`, `available` |
| `device42_businessapp` | `name`, `description` |
| `device42_certificate` | `issued_to`, `serial_number`, `subject`, `valid_from`, `valid_to` |

### 2.3 Idempotency (2nd apply) — drift findings

A second `terraform apply` on the same configuration reports **in-place updates** on the following resources. The read step overwrites fields with API default values that were not explicitly managed, producing a recurring diff:

| Resource | Drifting field(s) |
|---|---|
| `device42_building` | `tags` (order normalization) |
| `device42_device` | `service_level` |
| `device42_room` | grid start fields |
| `device42_businessapp` | owner / migration / service-level ID fields |
| `device42_software` | `licensing_model`, `track_licensed_count_by_keys` (and **`Update` fails to extract the software ID**) |

These are the same `Optional`-vs-`Computed` (unmanaged default-valued field) class fixed for `device42_device.building` in this branch. They are **not** create/read/destroy regressions but should be addressed for a fully stable `plan`.

---

## 3. Environment notes

- Tests ran against a local provider build via Terraform `dev_overrides` (`terraform.rc` pointing to the local provider binary).
- Provider `timeout` was raised to 90s for the slow software endpoints during the functional run.
- No Demo/CI instance credentials or host are recorded in this report.