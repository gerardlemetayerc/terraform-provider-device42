resource "device42_vendor" "acme" {
  name         = "Acme Corp"
  notes        = "Hardware supplier"
  home_page    = "https://acme.example.com"
  contact_info = "support@acme.example.com"
  account_no   = "ACC-12345"
  aliases      = "Acme"

  custom_fields = {
    region = "eu"
  }
}