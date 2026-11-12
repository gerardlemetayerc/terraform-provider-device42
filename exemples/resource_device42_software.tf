resource "device42_software" "apache" {
  name          = "Apache HTTP Server"
  vendor        = "Apache Software Foundation"
  notes         = "Web server"
  software_type = "Server"
  category      = "Web"
  aliases       = "httpd"
  tags          = "prod,web"

  custom_fields = {
    owner = "platform-team"
  }
}