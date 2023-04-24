terraform {
  required_providers {
    device42 = {
      source  = "github.com/gerardlemetayerc/terraform-provider-device42"
    }
  }
}

provider "device42" {
  d42_username = "${var.d42_username}"
  d42_password = "${var.d42_password}"
  d42_host     = "${var.d42_host}"

  # If TLS certificate not trusted by your system, set following value to true
  d42_tls_unsecure = true

}

