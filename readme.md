# Terraform provider - Device42

This provider has the purpose to manage Device42 CMDB application. 

It currently provide support for following ressources :
- Business Application
- Business Application Element
- IP
- VLAN
- Subnet
- Device

## Usage exemples

Usage examples can be found [here](/exemples/)

## How to compile sources from your own side.

- For Windows users
```
go build -o build\terraform-provider-device42.exe 
```

- For Linux Users
```
go build -o build\terraform-provider-device42
```