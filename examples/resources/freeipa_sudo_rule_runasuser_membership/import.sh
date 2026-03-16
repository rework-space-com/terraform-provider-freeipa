# The import id uses the format: <sudo_rule_name>/<type>/<identifier>
# Use type "msrrau" for multi-runasuser membership (non-deprecated).
# Note: slash characters in the rule name must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_runasuser_membership.users-0
  id = "sudo-rule-test/msrrau/users-0"
}

resource "freeipa_sudo_rule_runasuser_membership" "users-0" {
  name       = "sudo-rule-test"
  runasusers = ["user01", "user02"]
  identifier = "users-0"
}
