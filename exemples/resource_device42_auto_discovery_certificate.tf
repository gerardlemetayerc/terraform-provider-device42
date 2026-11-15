resource "device42_auto_discovery_certificate" "myCertJob" {
  name         = "myCertJob"
  server       = "hostname.domain.local"
  ports        = "443"
  follow_chain = "yes"
}