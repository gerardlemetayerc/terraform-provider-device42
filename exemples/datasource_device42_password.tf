data "device42_password" "example" {
  username = "exampleUsername"
  category = "exampleCategory"
}

output "password_details" {
  value = {
    password = data.device42_password.example.password
    label    = data.device42_password.example.label
  }
}