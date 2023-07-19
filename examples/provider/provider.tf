terraform {
  required_providers {
    freeipa = {
      version = "0.1.1"
      source  = "[Terraform registry provider path]"
    }
  }
}

provider "freeipa" {
  host = "ipa.example.test"
  username = "admin"
  password = "123456789"
  insecure = true
}
