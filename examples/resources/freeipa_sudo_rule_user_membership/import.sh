# The import id uses the format: <sudo_rule_name>/<type>/<identifier>
# Use type "msru" for multi-user/group membership (non-deprecated).
# Note: slash characters in the rule name must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_user_membership.users-1
  id = "sudo-rule-test/msru/users-1"
}

resource "freeipa_sudo_rule_user_membership" "users-1" {
  name       = "sudo-rule-test"
  users      = ["user01"]
  identifier = "users-1"
}
