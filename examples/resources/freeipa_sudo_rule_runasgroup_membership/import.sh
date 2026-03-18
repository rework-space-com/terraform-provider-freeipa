# The import id uses the format: <sudo_rule_name>/<type>/<identifier>
# Use type "msrraug" for multi-runasgroup membership (non-deprecated).
# Note: slash characters in the rule name must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_runasgroup_membership.groups-0
  id = "sudo-rule-test/msrraug/groups-0"
}

resource "freeipa_sudo_rule_runasgroup_membership" "groups-0" {
  name        = "sudo-rule-test"
  runasgroups = ["group01", "group02"]
  identifier  = "groups-0"
}
