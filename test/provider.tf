terraform {
  required_providers {
    freeipa = {
      source  = "rework-space-com/freeipa"
    }
  }
}

provider "freeipa" {
  host = "ipa.ipatest.lan"
  username = "admin"
  password = "P@ssword"
  insecure = true
}