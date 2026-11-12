resource "device42_certificate" "myNewCertificate" {
    issued_to           = "myNewCertificate"
    serial_number       = "1234567890"
    subject             = "CN=myNewCertificate,O=Acme"
    san                 = "DNS:myNewCertificate,DNS:www.acme.com"
    valid_from          = "2024-01-01"
    valid_to            = "2025-01-01"
    vendor              = "Let's Encrypt"
    signature_algorithm = "sha256WithRSAEncryption"
    signature_hash      = "sha256"
    extended_key_usage  = "Server Authentication"
    issued_by           = "R3"
    parent_cert         = "MyParentCert"
    version             = "3"
    custom_fields = {
        "certificate_type" = "server"
    }
}