data "device42_software_details" "on_device" {
  device_id = "1014"
}

output "installed" {
  value = [for i in data.device42_software_details.on_device.instances : i.software]
}

data "device42_software_details" "runs_apache" {
  software_id = data.device42_software.apache.software_id
}

output "devices" {
  value = [for i in data.device42_software_details.runs_apache.instances : i.device]
}