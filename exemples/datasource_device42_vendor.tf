data "device42_vendor" "acme_ds" {
  name = "Acme Corp"
}

output "acme_id" {
  value = data.device42_vendor.acme_ds.vendor_id
}

data "device42_vendor" "all" {}

output "vendor_names" {
  value = [for v in data.device42_vendor.all.vendors : v.name]
}