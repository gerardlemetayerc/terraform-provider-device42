terraform {
  required_providers {
    device42 = {
      source  = "github.com/gerardlemetayerc/device42"
    }
  }
}

provider "device42" {
  D42_USERNAME = "yourusername"
  D42_PASSWORD = "yourpassword"
  D42_HOST     = "youhostname"

  # If TLS certificate not trusted by your system, set following value to true
  D42_TLS_INSECURE = true

}

