data "device42_password" "example" {
  username = "exampleUsername"
  id       = 123
}

output "password_details" {
  value = {
    password = data.device42_password.example.password
    label    = data.device42_password.example.label
  }
}