---
page_title: "Provider: FREEIPA"
description: |-
  
---

# FREEIPA Provider

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ca_certificate` (String) Path to the server's SSL CA certificate
- `host` (String) The FreeIPA host
- `insecure` (Boolean) Whether to verify the server's SSL certificate
- `password` (String, Sensitive) Password to use for connection
- `username` (String) Username to use for connection
