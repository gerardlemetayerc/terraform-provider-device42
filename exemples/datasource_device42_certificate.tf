data "device42_certificate" "searchCertificate" {
    name = "myNewCertificate"
}

data "device42_certificate" "searchCertificate2" {
    certificate_id = 1
}