# The import id uses the format: <hbac_policy_name>/<type>/<identifier>
# Use type "mh" for multi-host/hostgroup membership (non-deprecated).

import {
  to = freeipa_hbac_policy_host_membership.hbac-hosts-1
  id = "test-hbac/mh/hbac-hosts-1"
}

resource "freeipa_hbac_policy_host_membership" "hbac-hosts-1" {
  name       = "test-hbac"
  hosts      = ["ipaclient1.ipatest.lan", "ipaclient2.ipatest.lan"]
  identifier = "hbac-hosts-1"
}
