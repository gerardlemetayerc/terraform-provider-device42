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
- Building (`device42_building`)
- Room (`device42_room`)
- Rack (`device42_rack`)
- Customer (`device42_customer`)
- Customer Contact (`device42_customer_contact`)
- Certificate (`device42_certificate`)
- Vendor (`device42_vendor`)

**Data sources**
- Device (`device42_device`)
- IP (`device42_ip`)
- Suggested IP (`device42_ip_suggested`)
- VLAN (`device42_vlan`)
- Subnet (`device42_subnet`)
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