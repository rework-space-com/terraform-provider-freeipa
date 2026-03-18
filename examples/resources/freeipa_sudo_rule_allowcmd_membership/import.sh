# The import id uses the format: <sudo_rule_name>/<type>/<identifier>
# Use type "msrac" for multi-sudocmd/sudocmd_group membership (non-deprecated).
# Note: slash characters in the rule name or sudocmd paths must be percent-encoded (%2F).

import {
  to = freeipa_sudo_rule_allowcmd_membership.allowed_bash
  id = "sudo-rule-editors/msrac/allowed_bash"
}

resource "freeipa_sudo_rule_allowcmd_membership" "allowed_bash" {
  name       = "sudo-rule-editors"
  sudocmds   = ["/bin/bash"]
  identifier = "allowed_bash"
}
