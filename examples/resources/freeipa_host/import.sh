# The import id must be exactly the same as the name of the host, which must be the fqdn of the host.

import {
  to = freeipa_host.testhost
  id = "testhost.ipatest.lan"
}

resource "freeipa_host" "testhost" {
  name = "testhost.ipatest.lan"
}