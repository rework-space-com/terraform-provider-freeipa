# The import id uses the format: <hbac_policy_name>/<type>/<identifier>
# Use type "mu" for multi-user/group membership (non-deprecated).

import {
  to = freeipa_hbac_policy_user_membership.hbac-users-1
  id = "test-hbac/mu/hbac-users-1"
}

resource "freeipa_hbac_policy_user_membership" "hbac-users-1" {
  name       = "test-hbac"
  users      = ["user-2", "user-3"]
  identifier = "hbac-users-1"
}
