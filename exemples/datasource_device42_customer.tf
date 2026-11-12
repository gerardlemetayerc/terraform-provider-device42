data "device42_customer" "searchCustomer" {
  name = "acmeCorp"
}

data "device42_customer" "searchCustomer2" {
  customer_id = 1
}