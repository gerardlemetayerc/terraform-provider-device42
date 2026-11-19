# Changelog

All notable changes to the Device42 provider are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
### Changed
### Fixed
### Removed

## [1.4.3] - 2026-11-19

### Added

- **`device42_ip`**: new `tags` attribute (comma-separated) to set tags on the IP address. The API is queried with `include_cols=tags` so the value is read back and no drift occurs. Note the Device42 IP API only **appends** tags (never removes them), so the attribute is best used to add tags rather than to fully replace the tag set.

### Fixed

- **`device42_ip`**: the read now passes a pointer to `apiDevice42Get` so resty actually unmarshals the response body. Previously the read silently failed to populate the state (the struct was passed by value), which meant the resource could never be read back after creation.

## [1.4.2] - 2026-11-15

### Added

#### New resources
- `device42_cable` — manage physical cables
- `device42_power_circuit` — manage power circuits
- `device42_pdu_model` — manage PDU models
- `device42_pdu` — manage PDUs
- `device42_vrf_group` — manage VRF groups
- `device42_subnet_category` — manage subnet categories
- `device42_mac` — manage MAC addresses
- `device42_ip_nat` — manage IP NAT rules

#### AutoDiscovery: 10 resources and data sources
- `device42_auto_discovery_certificate`
- `device42_auto_discovery_nmap`
- `device42_auto_discovery_blade`
- `device42_auto_discovery_networks`
- `device42_auto_discovery_snmp_disc`
- `device42_auto_discovery_dns`
- `device42_auto_discovery_pingsweep`
- `device42_auto_discovery_ipmi`
- `device42_auto_discovery_cloudaccount`
- `device42_auto_discovery_portscan`

#### New data sources (read-only)
- `device42_cable`
- `device42_power_circuit`
- `device42_pdu_model`
- `device42_pdu`
- `device42_vrf_group`
- `device42_subnet_category`
- `device42_mac`
- `device42_ip_nat`

#### Documentation & examples
- Reference documentation for all new resources and data sources
- Usage examples for every new resource
- `readme.md` updated

### Fixed

- **`device42_device`**: an HTTP 500 `Device matching query does not exist` (code 3) is now treated as a deleted resource instead of an error. Detection requires the specific code **and** message so an unrelated server error does not drop a live device.
- **`device42_mac`**: `macaddress` is now `ForceNew` (recreate on change), the id is preserved on network read error, and `vlan_id`/`port` became `Computed`.
- **`device42_pdu`**, **`device42_power_circuit`**, **`device42_vrf_group`**: the id is preserved on network read error to avoid an accidental destructive recreate.
- **`device42_pdu`**: `gateway_pdu_id` is now sent in the update payload (fixes drift); `pdu_model`/`pdu_model_id` are cleared when the API no longer returns a model.
- **Updates (pdu, pdu_model, power_circuit, subnet_category, vrf_group)**: HTTP responses ≥ 400 are rejected before inspecting the API code (resty left the code at zero and masked update failures).
- **`device42_auto_discovery_dns`**: documentation corrected — the API always creates a new job, there is no upsert.

### Tests
- Added `TestIsDeviceNotFound` (5 cases) and updated `TestSuppressTagsAppendDiff`.
- `go build`, `go vet`, `go test` all pass.

## [1.4.1] - 2025-XX-XX

### Added
- (content of the previous release — see git history)

[Unreleased]: https://github.com/gerardlemetayerc/terraform-provider-device42/compare/v1.4.3...HEAD
[1.4.3]: https://github.com/gerardlemetayerc/terraform-provider-device42/releases/tag/v1.4.3
[1.4.2]: https://github.com/gerardlemetayerc/terraform-provider-device42/releases/tag/v1.4.2
[1.4.1]: https://github.com/gerardlemetayerc/terraform-provider-device42/releases/tag/v1.4.1