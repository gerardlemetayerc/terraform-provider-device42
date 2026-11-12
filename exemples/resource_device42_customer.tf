resource "device42_customer" "acmeCorp" {
    name          = "acmeCorp"
    manager       = "John Doe"
    contact_info  = "contact@acme.com"
    notes         = "Main customer"
    tags          = "Prod,East"
    custom_fields = {
        "customer_type" = "enterprise"
    }
}