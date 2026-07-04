import {
  to = freeipa_sysaccount.sudo
  id = "sudo"
}

resource "freeipa_sysaccount" "sudo" {
  name       = "sudo"
}
