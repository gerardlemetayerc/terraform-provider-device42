data "device42_software" "apache_ds" {
  name = "Apache HTTP Server"
}

output "apache_id" {
  value = data.device42_software.apache_ds.software_id
}

data "device42_software" "all" {}

output "software_names" {
  value = [for s in data.device42_software.all.software : s.name]
}