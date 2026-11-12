data "device42_service_level" "prod" {
  name = "Production"
}

output "production_id" {
  value = data.device42_service_level.prod.service_level_id
}

data "device42_service_level" "all" {}

output "service_level_names" {
  value = [for l in data.device42_service_level.all.levels : l.name]
}