resource "device42_customer_contact" "johnDoe" {
    customer = "acmeCorp"
    name     = "John Doe"
    type     = "primary"
    email    = "john.doe@acme.com"
    phone    = "+1 555 000 000"
    address  = "100 Main Street, New York"
}