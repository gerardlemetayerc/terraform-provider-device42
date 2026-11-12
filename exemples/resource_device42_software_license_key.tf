resource "device42_software" "apache" {
  name   = "Apache HTTP Server"
  vendor = "Apache Software Foundation"
}

resource "device42_software_license_key" "apache" {
  key           = "LPKG-1234-ABCD-5678"
  software_id   = device42_software.apache.software_id
  license_count = 10
  notes         = "Enterprise subscription"
}