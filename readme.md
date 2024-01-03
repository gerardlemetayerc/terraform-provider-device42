# Terraform provider - Device42

This provider has the purpose to manage Device42 CMDB application. 

It currently provide support for following ressources :
- Business Application
- Business Application Element
- IP
- VLAN
- Subnet
- Device

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