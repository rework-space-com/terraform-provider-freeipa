# The import id uses the format: <sudo_rule_name>/<type>/<identifier>
# Use type "msrh" for multi-host/hostgroup membership (non-deprecated).
# Note: slash characters in the rule name must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_host_membership.hosts-0
  id = "sudo-rule-test/msrh/hosts-0"
}

resource "freeipa_sudo_rule_host_membership" "hosts-0" {
  name       = "sudo-rule-test"
  hosts      = ["test.example.test"]
  identifier = "hosts-0"
}
