terraform {
  required_providers {
    freeipa = {
      version = "5.1.0"
      source  = "rework-space-com/freeipa"
    }
  }
}

provider "freeipa" {
  host     = "ipa.example.test"
  username = "admin"
  password = "123456789"
  insecure = true
}
