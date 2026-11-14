# Terraform provider - Device42

This provider has the purpose to manage Device42 CMDB application. 

It currently provides support for the following resources and data sources :

**Resources**
- Business Application (`device42_businessapp`)
- Business Application Element (`device42_businessappelement`)
- Device (`device42_device`)
- IP (`device42_ip`)
- VLAN (`device42_vlan`)
- Subnet (`device42_subnet`)
- Suggested Subnet (`device42_subnet_suggested`)
- Building (`device42_building`)
- Room (`device42_room`)
- Rack (`device42_rack`)
- Customer (`device42_customer`)
- Customer Contact (`device42_customer_contact`)
- Certificate (`device42_certificate`)
- Vendor (`device42_vendor`)
- Software (`device42_software`)
- Software License Key (`device42_software_license_key`)
- Cable (`device42_cable`)
- Power Circuit (`device42_power_circuit`)
- PDU Model (`device42_pdu_model`)
- PDU (`device42_pdu`)
- VRF Group (`device42_vrf_group`)
- Subnet Category (`device42_subnet_category`)
- MAC (`device42_mac`)
- IP NAT (`device42_ip_nat`)

**Data sources**
- Device (`device42_device`)
- IP (`device42_ip`)
- IPs (`device42_ips`)
- Suggested IP (`device42_ip_suggested`) — **deprecated**: prefer `device42_ip`
- VLAN (`device42_vlan`)
- VLANs (`device42_vlans`)
- Subnet (`device42_subnet`)
- Subnets (`device42_subnets`)
- Business Application (`device42_businessapp`)
- Password (`device42_password`)
- Building (`device42_building`)
- Room (`device42_room`)
- Rack (`device42_rack`)
- Customer (`device42_customer`)
- Customer Contact (`device42_customer_contact`)
- Certificate (`device42_certificate`)
- Service Level (`device42_service_level`)
- Vendor (`device42_vendor`)
- Software (`device42_software`)
- Software Details (`device42_software_details`)
- Cable (`device42_cable`)
- Power Circuit (`device42_power_circuit`)
- PDU Model (`device42_pdu_model`)
- PDU (`device42_pdu`)
- VRF Group (`device42_vrf_group`)
- Subnet Category (`device42_subnet_category`)
- MAC (`device42_mac`)
- IP NAT (`device42_ip_nat`)

## Version v1.4.2

This release adds support for Device42 cables, connecting two endpoints (patch panel ports, netports, circuits), validated against the Device42 public API.

### What's New

- Added `device42_cable` resource to manage cables linking two Device42 endpoints (`origin`/`end` of type `patch_panel_port`, `cable`, `circuit`, `netport` or `tap_port`), with length, units, notes, vendor, tags and custom fields. The `cable_id` is the unique logical name used by the API to upsert.
- Added `device42_power_circuit` resource to manage power circuits attached to a breaker panel (`breakerpanel_id` + `number`), with label, amps, voltage and connection.
- Added `device42_pdu_model` resource to manage PDU models (manufacturer, depth, size, width, height) and `device42_pdu` resource to manage physical PDUs, linked to a PDU model by name.
- Added `device42_vrf_group` resource to manage VRF groups (name, description, tags). Tags are appended by the API on update, so drift is suppressed when the configured set is a subset of the state.
- Added `device42_subnet_category` resource to manage subnet categories, `device42_mac` resource to manage MAC addresses (upserted by address), and `device42_ip_nat` resource to manage IP NAT rules. The IP NAT `ip_address_to` field is not updatable by the API and therefore forces a recreate when changed.
- Added matching data sources for the new resources: `device42_cable`, `device42_power_circuit`, `device42_pdu_model`, `device42_pdu`, `device42_vrf_group`, `device42_subnet_category`, `device42_mac` and `device42_ip_nat`.

## Version v1.4.1

This release adds support for the physical infrastructure layer of Device42 (Building -> Room -> Rack), validated against the Device42 public API.

### What's New

- Added `device42_building` resource and datasource to manage Device42 buildings (address, contact, latitude/longitude, tags, custom fields, import by name).
- Added `device42_room` resource and datasource to manage rooms, linked to a building via `building_id` or `building`.
- Added `device42_rack` resource and datasource to manage racks, with size, manufacturer, building/room placement and orientation.
- Added `device42_customer` resource and datasource to manage Device42 customers (manager, notes, contact_info, tags, custom fields, import by name).
- Added `device42_customer_contact` resource and datasource to manage customer contacts (name, type, email, phone, address).
- Added `device42_certificate` resource and datasource to manage certificates (subject, SAN, validity dates, serial number, vendor, custom fields, import by common name).
- Added import support to the `device42_ip` resource.
- Added `device42_ip` datasource to look up an existing Device42 IP by its address (subnet, availability, attached device).
- Added `device42_service_level` datasource to read Device42 service levels, either by name or as a full list (`TypeList`).
- Added `device42_vendor` resource to manage vendors (upsert by name, import by `vendor_id`, custom fields) and datasource to look up a vendor by name or list all vendors.
- Added `device42_software` resource and datasource to manage software components (vendor, description, type, category, aliases, custom fields, import by `software_id`).
- Added `device42_software_license_key` resource to manage software license keys.
- Added `device42_software_details` datasource to list software-in-use instances, optionally filtered by software or device.
- Added `device42_subnets` datasource to read multiple subnets at once, filtered by parent VLAN, name or VRF group (the plural counterpart of `device42_subnet`).
- Added `device42_ips` datasource to read the IPs of a subnet (server-side `subnet_id` filter, since the `/2.0/ips/` API is paginated at 1000 records).
- Added `device42_vlans` datasource to read multiple VLANs at once, optionally filtered by name or id (the plural counterpart of `device42_vlan`).
- Added `device42_subnet_suggested` resource to suggest and reserve the next available sub-subnet under a parent subnet (`GET /1.0/suggest_subnet/{id}/?mask_bits=` then `POST /1.0/subnets/`, so the suggestion is consumed and repeated applies stay stable).
- All resources support import (`terraform import`) and custom fields, reusing the existing `device42` error-handling and custom-field helpers.
- Fixed type handling for building `latitude`/`longitude` (the API returns JSON numbers or null, not always strings) and ensured room custom fields survive a read refresh.
- Fixed device archive response decoding (the archive endpoint returns a bare JSON array, not the usual `{code,msg}` object).
- Fixed certificate read to detect a deleted certificate via the HTTP status code (resty does not surface a 404 as an error).

## Version v1.3.13

This release includes new features :
- Implemented import functionnality on device (works now)
- Added build for Mac OS (darwin architecture - issue #70 reported by @gowgopal83)
- Fixed issue on datasource (resty client was displaying too much information because debug mode)


## Version v1.3.10

This release includes new features, performance improvements, enhanced debugging, and critical bug fixes.

### What's New

- Added `device42_password` datasource for fetching password details from Device42.
- Dependencies updated: `github.com/go-resty/resty/v2` to v2.11, `github.com/hashicorp/terraform-plugin-sdk/v2` to v2.31.
- Go language updated to version 1.20 for improved performance and stability.
- Enhanced error handling and response validation for better reliability and clearer error messages.

### Installation

To install this provider, follow the standard procedure for installing a Terraform plugin.

### Documentation

For detailed information on the available resources and datasources, visit [Terraform Registry](https://registry.terraform.io/providers/gerardlemetayerc/device42/latest/docs.

## Usage exemples

Usage examples can be found [here](/exemples/)

### Contributions

Contributions to this provider are welcome. Please refer to the contributing guidelines for more information.

## How to compile sources from your own side.

- For Windows users
```
go build -o build\terraform-provider-device42.exe 
```

- For Linux Users
```
go build -o build\terraform-provider-device42
```