data "device42_vendor" "acme" {
  name = "Acme Corp"
}

output "acme_id" {
  value = data.device42_vendor.acme.vendor_id
}

data "device42_vendor" "all" {}

output "vendor_names" {
  value = [for v in data.device42_vendor.all.vendors : v.name]
}