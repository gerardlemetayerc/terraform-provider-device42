data "device42_customer_contact" "searchContact" {
    customer = "acmeCorp"
    name     = "John Doe"
}

data "device42_customer_contact" "searchContact2" {
    contact_id = 1
}